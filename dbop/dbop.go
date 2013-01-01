// Package dbop includes the most commonly used db operations - select, insert, delete
// in an easy to use manner without the need to use or know anything about the underlying 
// sql queries.
package dbop

import (
		"fmt"
		"database/sql"
		_ "github.com/ziutek/mymysql/godrv"
		"strconv"
)

type Table interface {
	InitTable(tableName string, fieldNames []string, fieldTypes []string)
	ResetTable()
	GetTableName() string
	GetFieldNameList() []string
	GetFieldTypeList() []string
	GetFieldValue(fieldName string) string
	SetFieldValue(fieldName string, fieldValue string) bool
}

func toStmtStr(value string, valueType string) string {
	switch valueType {
		case "BIT", "TINYINT", "BOOL", "BOOLEAN", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "SERIAL", "DECIMAL", "DEC", "FLOAT", "DOUBLE", "YEAR":
			return value
			
		case "DATE","DATETIME", "TIMESTAPM", "TIME", "CHAR", "VARCHAR", "BINARY", "VARBINARY", "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT", "ENUM", "SET":
			return fmt.Sprintf("'%s'",value)
	}
	
	return ""
}

func anytypeToStr(value interface{}) string {
	switch value.(type) {
		case int, int8, int16, int32, int64:
			return strconv.FormatInt(value.(int64),10)
		case uint, uint8, uint16, uint32, uint64:
			return strconv.FormatUint(value.(uint64), 10)
		case float32, float64:
			return strconv.FormatFloat(value.(float64), 'f', -1, 64)
		case bool:
			return strconv.FormatBool(value.(bool))
		default:
			return fmt.Sprintf("%s", value)
	}
	
	return ""
}

type RecId struct {
	Value 	int64
	Exists 	bool
	AutoInc bool
	IsSet	bool
}

// Defines the database table base type
type DbTable struct {
	tableName 		string
	fieldNames		[]string
	fieldTypes		[]string
	fieldValue 		[]string
	fieldValueSet   []bool
	recid			RecId
}

// Returns the recid value of the table and the IsSet value. IsSet will be true if the 
// value has been set and not read, it will be false if the value is empty or has been read from db
// The methd will panic if recid does not exist for this table.
func (t DbTable) RecId() (int64, bool) {
	if !t.recid.Exists {
		panic ("Rec id doesn't exist for this table")
	}
	
	return t.recid.Value, t.recid.IsSet
}

// Sets the value for recid to be used when executing db operations. This value will be included in ther
// where clause. If 0 is passed in, the method will clear the recid value. This method also must be used
// if recid field is not set as auto_increment in the database and must be maintained in the application.
// The methd will panic if recid does not exist for this table.
func (t *DbTable) SetRecId(recId int64) {
	if !t.recid.Exists {
		panic ("Rec id dosn't exist for this table")
	}
	
	if recId != 0 {
		t.recid.Value = recId
		t.recid.IsSet = true
	} else {
		t.recid.Value = recId
		t.recid.IsSet = false
	}
}

func (t DbTable) newTableInstance() DbTable {
	var tbl 		DbTable
	var fieldNames	[]string
	var fieldTypes 	[]string
	var recid		[2]bool
	
	fieldNames = make([]string, len(t.fieldNames))
	copy(fieldNames, t.fieldNames)
	
	fieldTypes = make([]string, len(t.fieldTypes))
	copy(fieldTypes, t.fieldTypes)
	
	recid[0] = t.recid.Exists
	recid[1] = t.recid.AutoInc
	
	tbl.InitTable(t.tableName, fieldNames, fieldTypes, recid)
	
	return tbl;
}

// Initiates the base type with info from a specific table in the database.
// recid is a field and a unique index for that field that can be created on the table
// for easily performing DoUpdate and DoDelete operations after selecting a single record
func (t *DbTable) InitTable(tableName string, fieldNames []string, fieldTypes []string, recid [2]bool) {
	t.tableName 	= tableName
	t.fieldNames 	= fieldNames
	t.fieldTypes 	= fieldTypes
	t.fieldValue 	= make([]string, len(fieldTypes))
	t.fieldValueSet = make([]bool, len(fieldTypes))
	t.recid.Exists	= recid[0]
	if t.recid.Exists {
		t.recid.AutoInc = recid[1]
	}
}

// Resets the table variable for initiating as a different database table
func (t *DbTable) ResetTable() {
	t.tableName 	= ""
	t.fieldNames 	= nil
	t.fieldTypes 	= nil
	t.recid.AutoInc = false
	t.recid.Exists	= false
	t.recid.Value 	= 0
}

// Returns a slice of all the field names for the initiated table
func (t DbTable) GetFieldNameList() []string {
	return t.fieldNames
}

// Returns a slice of all the field types for the initiated table
func (t DbTable) GetFieldTypeList() []string {
	return t.fieldTypes
}

// Returns a the table name of the initiated table
func (t DbTable) GetTableName() string {
	return t.tableName
}

// Returns the value of a field specified by the field name
func (t DbTable) GetFieldValue(fieldName string) string {
	for fId, fn := range t.fieldNames {
		if fieldName == fn {
			return t.fieldValue[fId]
		}
	}
	
	return ""
}

// Sets the value of a field
func (t *DbTable) SetFieldValue(fieldName string, fieldValue string) bool {
	for fId, fn := range t.fieldNames {
		if fieldName == fn {
			t.fieldValue[fId] = fieldValue
			t.fieldValueSet[fId] = true
			return true;
		}
	}
	
	return false
}

// Clears all field values
func (t *DbTable) ClearFields() {
	for fId := 0; fId < len(t.fieldValue); fId++ {
		t.fieldValue[fId] = ""
		t.fieldValueSet[fId] = false
	}
}

// Clears the value of a field specified by the field name
func (t *DbTable) ClearField(fieldName string) bool {
	for fId, fn := range t.fieldNames {
		if fieldName == fn {
			t.fieldValue[fId] = ""
			t.fieldValueSet[fId] = false
			return true
		}
	}
	
	return false
}

// Database connection type used when executing a db operation
type DbConnection struct {
	connection *sql.DB
}

// Opens a database connection
func (dbc *DbConnection) Open(connectionStr string) {
	con, err := sql.Open("mymysql", connectionStr)
	
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	
	dbc.connection = con
}

// Executes a custom sql statement
func (dbc *DbConnection) Exec(queryStr string) (int64, error) {
	result, err := dbc.connection.Exec(queryStr)
	
	if err != nil {
		return 0, err
	}
	
	rowsAffected, _ := result.RowsAffected()
	
	return rowsAffected, nil
}

func (t *DbTable) buildSelectStr(firstonly bool) string {
	var selectStr 	string
	var whereStr 	string
	var hasWhere 	bool
	
	selectStr = "SELECT * FROM " + t.tableName
	
	for fId, isSet := range t.fieldValueSet {
		if isSet {
			hasWhere = true
			if len(whereStr) != 0 {
				whereStr = whereStr + " AND "
			}
			whereStr = whereStr + t.tableName + "." + t.fieldNames[fId] + " = " + toStmtStr(t.fieldValue[fId], t.fieldTypes[fId])
			
		}
	}
	
	if t.recid.Exists && t.recid.IsSet {
		if len(whereStr) != 0 {
			whereStr = whereStr + " AND "
		}
		whereStr = whereStr + t.tableName + ".recid = " + toStmtStr(anytypeToStr(t.recid.Value), "BIGINT")
	}
	
	if hasWhere {
		selectStr = selectStr + " WHERE " + whereStr
	}
	
	if firstonly {
		selectStr = selectStr + " LIMIT 1"
	}
	
	return selectStr
}

// Builds and executes a select statement based on the field values that have been set using
// using the SetFieldValue() function. Will return true if successfull and will populate the 
// field values for the variable called from. All fields returned by the db will be populated,
// but fields will be considered not set. Selects only the first line from the table.
func (t *DbTable) DoSelectFirstonly(dbc *DbConnection) bool {
	row := dbc.connection.QueryRow(t.buildSelectStr(true))
	
	fields := make([]interface{}, len(t.fieldNames))
	fieldValues := make([]*interface{}, len(t.fieldNames))
	
	for fId := range fields {
		fields[fId] = &fieldValues[fId]
	}
	
	err := row.Scan(fields...)	
	
	if err != nil {
		fmt.Printf("%s\n", err)
		return false
	}
	
	for fId := range fieldValues {
		value := *fieldValues[fId]
		if t.recid.Exists && fId == 1 {
			t.recid.Value = value.(int64)
			t.recid.IsSet = false
		} else {
			t.fieldValue[fId] = anytypeToStr(value)
			t.fieldValueSet[fId] = false
		}
	}
	
	return true
}

// Builds and executes a select statement based on the field values that have been set using
// using the SetFieldValue() function. Will return a slice of DbTable objects that represent
// the selected table rows. If no lines are found, will return a 0 sized slice. Will return 
// nil value and an error if a problem was encountered.
func (t DbTable) DoSelect(dbc *DbConnection) ([]DbTable, error) {
	var retRows 	[]DbTable
	var counter		int
	
	rows, err := dbc.connection.Query(t.buildSelectStr(false))
	
	if err != nil {
		return nil, err
	}	
		
	for rows.Next() {
		fields := make([]interface{}, len(t.fieldNames))
		fieldValues := make([]*interface{}, len(t.fieldNames))
		
		for fId := range fields {
			fields[fId] = &fieldValues[fId]
		}
	
		err := rows.Scan(fields...)		
		
		if err != nil {
			return nil, err
		}
		
		tableRow := t.newTableInstance()
		
		for fId := range fieldValues {
			value := *fieldValues[fId]
			if t.recid.Exists && fId == 1 {
				t.recid.Value = value.(int64)
				t.recid.IsSet = false
			} else {
				t.fieldValue[fId] = anytypeToStr(value)
				t.fieldValueSet[fId] = false
			}
		}
				
		retRows = append(retRows, tableRow)
		counter++
	}
	
	rows.Close()
	
	return retRows, nil	
}

func (t DbTable) buildInsertStr() (string, error) {
	var stmtStr		string
	var stmtFields 	string
	var stmtValues	string
	
	stmtStr = "INSERT INTO " + t.tableName + " "
	
	for fId := range t.fieldNames {
		if t.fieldValueSet[fId] {
			if len(stmtFields) == 0 {
				stmtFields = stmtFields + "("
				stmtValues = stmtValues + "("
			} else {
				stmtFields = stmtFields + ","
				stmtValues = stmtValues + ","
			}
			
			stmtFields = stmtFields + t.fieldNames[fId]
			stmtValues = stmtValues + toStmtStr(t.fieldValue[fId], t.fieldTypes[fId])
		}
	}
	
	if t.recid.Exists && !t.recid.AutoInc {
		stmtFields = stmtFields + "recid, "
		stmtValues = stmtFields + toStmtStr(anytypeToStr(t.recid.Value), "BIGINT")
	}
	
	if len(stmtFields) == 0 {
		return "", fmt.Errorf("No fields set!")
	} else {
		stmtFields = stmtFields + ")"
		stmtValues = stmtValues + ")"
	}
	
	stmtStr = stmtStr + stmtFields + " VALUES " + stmtValues
	
	return stmtStr, nil
}

// Builds anx executes an insert statement from the set field values
func (t DbTable) DoInsert(dbc *DbConnection) (int64, error) {
	stmtStr, err := t.buildInsertStr()
	
	if err != nil {
		return 0, err
	}
	
	return dbc.Exec(stmtStr)
}
