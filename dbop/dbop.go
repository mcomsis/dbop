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

func (t DbTable) newTableInstance() DbTable {
	var tbl 		DbTable
	var fieldNames	[]string
	var fieldTypes 	[]string
	
	fieldNames = make([]string, len(t.fieldNames))
	copy(fieldNames, t.fieldNames)
	
	fieldTypes = make([]string, len(t.fieldTypes))
	copy(fieldTypes, t.fieldTypes)
	
	tbl.InitTable(t.tableName, fieldNames, fieldTypes)
	
	return tbl;
}

// Defines the database table base type
type DbTable struct {
	tableName 		string
	fieldNames		[]string
	fieldTypes		[]string
	fieldValue 		[]string
	fieldValueSet   []bool
}

// Initiates the base type with info from a specific table in the database.
func (t *DbTable) InitTable(tableName string, fieldNames []string, fieldTypes []string) {
	t.tableName 	= tableName
	t.fieldNames 	= fieldNames
	t.fieldTypes 	= fieldTypes
	t.fieldValue 	= make([]string, len(fieldTypes))
	t.fieldValueSet = make([]bool, len(fieldTypes))
}

// Resets the table variable for initiating as a different database table
func (t *DbTable) ResetTable() {
	t.tableName = ""
	t.fieldNames = nil
	t.fieldTypes = nil
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
func (dbc *DbConnection) Exec(queryStr string) int64 {
	result, err := dbc.connection.Exec(queryStr)
	
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	
	return rowsAffected
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
		t.fieldValue[fId] = anytypeToStr(value)
		t.fieldValueSet[fId] = false
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
			tableRow.fieldValue[fId] = anytypeToStr(value)
			tableRow.fieldValueSet[fId] = false						
		}
				
		retRows = append(retRows, tableRow)
		counter++
	}
	
	rows.Close()
	
	return retRows, nil	
}
