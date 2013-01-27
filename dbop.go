// Package dbop includes the most commonly used db operations - select, insert, delete
// in an easy to use manner without the need to use or know anything about the underlying 
// sql queries.
package dbop

import (
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	"strconv"
)

func toStmtStr(value string, valueType string) string {
	switch valueType {
	case "BIT", "TINYINT", "BOOL", "BOOLEAN", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "SERIAL", "DECIMAL", "DEC", "FLOAT", "DOUBLE", "YEAR":
		return value

	case "DATE", "DATETIME", "TIMESTAMP", "TIME", "CHAR", "VARCHAR", "BINARY", "VARBINARY", "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT", "ENUM", "SET":
		return "'" + value + "'"
	}

	return ""
}

func anytypeToStr(value interface{}) string {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(value.(int64), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case float32, float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(value.(bool))
	default:
		return fmt.Sprintf("%s", value)
	}

	return ""
}

type RecId struct {
	Value   uint64
	Exists  bool
	AutoInc bool
	IsSet   bool
}

type DbUpdateField struct {
	FieldName string
	Value     string
}

// Defines the database table base type
type DbTable struct {
	tableName     string
	fieldNames    []string
	fieldTypes    []string
	fieldValue    []string
	fieldValueSet []bool
	recid         RecId
}

// Returns the recid value of the table and the IsSet value. IsSet will be true if the 
// value has been set and not read, it will be false if the value is empty or has been read from db
// The method will panic if recid does not exist for this table.
func (t DbTable) RecId() (uint64, bool) {
	if !t.recid.Exists {
		panic("Rec id doesn't exist for this table")
	}

	return t.recid.Value, t.recid.IsSet
}

// Sets the value for recid to be used when executing db operations. This value will be included in ther
// where clause. If 0 is passed in, the method will clear the recid value. This method also must be used
// if recid field is not set as auto_increment in the database and must be maintained in the application.
// The method will panic if recid does not exist for this table.
func (t *DbTable) SetRecId(recId uint64) {
	if !t.recid.Exists {
		panic("Rec id dosn't exist for this table")
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
	var tbl DbTable
	var fieldNames []string
	var fieldTypes []string
	var recid [2]bool

	fieldNames = make([]string, len(t.fieldNames))
	copy(fieldNames, t.fieldNames)

	fieldTypes = make([]string, len(t.fieldTypes))
	copy(fieldTypes, t.fieldTypes)

	recid[0] = t.recid.Exists
	recid[1] = t.recid.AutoInc

	tbl.InitTable(t.tableName, fieldNames, fieldTypes, recid)

	return tbl
}

// Initiates the base type with info from a specific table in the database.
// recid is a field and a unique index for that field that can be created on the table
// for easily performing DoUpdate and DoDelete operations after selecting a single record
func (t *DbTable) InitTable(tableName string, fieldNames []string, fieldTypes []string, recid [2]bool) {
	t.tableName = tableName
	t.fieldNames = fieldNames
	t.fieldTypes = fieldTypes
	t.fieldValue = make([]string, len(fieldTypes))
	t.fieldValueSet = make([]bool, len(fieldTypes))
	t.recid.Exists = recid[0]
	if t.recid.Exists {
		t.recid.AutoInc = recid[1]
	}
}

// Resets the table variable for initiating as a different database table
func (t *DbTable) ResetTable() {
	t.tableName = ""
	t.fieldNames = nil
	t.fieldTypes = nil
	t.recid.AutoInc = false
	t.recid.Exists = false
	t.recid.Value = 0
}

// Returns a slice of all the field names for the initiated table
func (t DbTable) GetFieldNameList() []string {
	return t.fieldNames
}

// Returns a slice of all the field types for the initiated table
func (t DbTable) GetFieldTypeList() []string {
	return t.fieldTypes
}

// Return field db type from a field name
func (t DbTable) GetFieldType(fieldName string) string {
	for fId, name := range t.fieldNames {
		if name == fieldName {
			return t.fieldTypes[fId]
		}
	}

	panic("Field not found")
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
			return true
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

	t.recid.Value = 0
	t.recid.IsSet = false
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

	if fieldName == "recid" {
		t.recid.Value = 0
		t.recid.IsSet = false
		return true
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
		panic(err)
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

func (t DbTable) buildSelectStr(firstonly bool) (string, error) {
	var selectStr string
	var whereStr string
	var hasWhere bool

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
		hasWhere = true
	}

	if hasWhere {
		selectStr = selectStr + " WHERE " + whereStr
	}

	if firstonly {
		selectStr = selectStr + " LIMIT 1"
	}

	// for debugging	
	fmt.Printf("%v\n", selectStr)

	return selectStr, nil
}

// Builds and executes a select statement based on the field values that have been set using
// using the SetFieldValue() function. Will return true if successfull and will populate the 
// field values for the variable called from. All fields returned by the db will be populated,
// but fields will be considered not set. Selects only the first line from the table.
func (t *DbTable) DoSelectFirstonly(dbc *DbConnection) error {
	var offset int

	queryStr, err := t.buildSelectStr(true)

	if err != nil {
		return err
	}

	row := dbc.connection.QueryRow(queryStr)

	fieldCount := len(t.fieldNames)

	if t.recid.Exists {
		fieldCount = fieldCount + 1
	}

	fields := make([]interface{}, fieldCount)
	fieldValues := make([]*interface{}, fieldCount)

	for fId := range fields {
		fields[fId] = &fieldValues[fId]
	}

	err = row.Scan(fields...)

	if err != nil {
		return err
	}

	if t.recid.Exists {
		offset = 1
	} else {
		offset = 0
	}

	for fId := range fieldValues {
		value := *fieldValues[fId]
		if t.recid.Exists && fId == 0 {
			int64value := value.(int64)
			t.recid.Value = uint64(int64value)
			t.recid.IsSet = false
		} else {
			t.fieldValue[fId-offset] = anytypeToStr(value)
			t.fieldValueSet[fId-offset] = false
		}
	}

	return nil
}

// Builds and executes a select statement based on the field values that have been set using
// using the SetFieldValue() function. Will return a slice of DbTable objects that represent
// the selected table rows. If no lines are found, will return a 0 sized slice. Will return 
// nil value and an error if a problem was encountered.
func (t DbTable) DoSelect(dbc *DbConnection) ([]DbTable, error) {
	var retRows []DbTable
	var counter int
	var offset int

	queryStr, err := t.buildSelectStr(false)

	if err != nil {
		return nil, err
	}

	rows, err := dbc.connection.Query(queryStr)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		fieldCount := len(t.fieldNames)

		if t.recid.Exists {
			fieldCount++
		}

		fields := make([]interface{}, fieldCount)
		fieldValues := make([]*interface{}, fieldCount)

		for fId := range fields {
			fields[fId] = &fieldValues[fId]
		}

		err = rows.Scan(fields...)

		if err != nil {
			return nil, err
		}

		tableRow := t.newTableInstance()

		if t.recid.Exists {
			offset = 1
		} else {
			offset = 0
		}

		for fId := range fieldValues {
			value := *fieldValues[fId]
			if t.recid.Exists && fId == 0 {
				int64value := value.(int64)
				tableRow.recid.Value = uint64(int64value)
				tableRow.recid.IsSet = false
			} else {
				tableRow.fieldValue[fId-offset] = anytypeToStr(value)
				tableRow.fieldValueSet[fId-offset] = false
			}
		}

		retRows = append(retRows, tableRow)
		counter++
	}

	rows.Close()

	return retRows, nil
}

func (t DbTable) buildInsertStr() (string, error) {
	var stmtStr string
	var stmtFields string
	var stmtValues string

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

	// for debugging
	fmt.Printf("%v\n", stmtStr)

	return stmtStr, nil
}

// Builds and executes an insert statement from the set field values
func (t *DbTable) DoInsert(dbc *DbConnection) error {
	stmtStr, err := t.buildInsertStr()

	if err != nil {
		return err
	}

	rows, err := dbc.Exec(stmtStr)

	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("Something went wrong, insert affected %v rows.", rows)
	}

	if t.recid.Exists && t.recid.AutoInc {
		err = t.DoSelectFirstonly(dbc)
	}

	if err != nil {
		return err
	}

	return nil
}

func (t DbTable) buildDeleteStr(useRecId bool) (string, error) {
	var hasWhere bool
	var whereStr string

	deleteStr := "DELETE FROM " + t.tableName

	if useRecId {
		whereStr = t.tableName + ".recid = " + toStmtStr(anytypeToStr(t.recid.Value), "BIGINT")
		hasWhere = true
	} else {
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
	}

	if !hasWhere {
		return "", fmt.Errorf("Delete must have a where clause")
	}

	deleteStr = deleteStr + " WHERE " + whereStr

	// for debugging
	fmt.Printf("%v\n", deleteStr)

	return deleteStr, nil
}

// Deletes the selected record. If no record has previously been selected (recid has no value or it
// has been set manually), an error will be returned. DoDelete will only work for tables that have
// the recid field. To make sure a record has been selected, check if it has a recid. 
func (t *DbTable) DoDelete(dbcon *DbConnection) error {
	if !t.recid.Exists || t.recid.IsSet || t.recid.Value == 0 {
		return fmt.Errorf("No record has been selected, cant DoDelete()!")
	}

	deleteStr, err := t.buildDeleteStr(true)

	if err != nil {
		return err
	}

	rows, err := dbcon.Exec(deleteStr)

	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("Something went wrong, %i lines deleted. SQL query: %s", rows, deleteStr)
	}

	t.ClearFields()

	return nil
}

// Deletes all lines using the values that have been set in the fields. This will also include the recid
// if it exists and has been set. This can be used for deleting records in bulk or deleting a specific 
// record by its recid without first selecting it. Will return the number of rows deleted or an error if
// something went wrong. If no rows fit the criteria, 0 and no error will be returned. 
func (t *DbTable) DoDeleteWhere(dbcon *DbConnection) (int64, error) {
	deleteStr, err := t.buildDeleteStr(false)

	if err != nil {
		return 0, err
	}

	rows, err := dbcon.Exec(deleteStr)

	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (t DbTable) buildUpdateStr(useRecId bool, whereFields []DbUpdateField) (string, error) {
	var hasWhere bool
	var hasSet bool
	var whereStr string
	var setStr string

	if useRecId && (!t.recid.Exists || t.recid.Value == 0) {
		return "", fmt.Errorf("Record has not been selected or the table does not use recid field.")
	}

	queryStr := "UPDATE " + t.tableName + " SET "

	for fId, isSet := range t.fieldValueSet {
		if isSet {
			if len(setStr) > 0 {
				setStr = setStr + ", "
			}

			setStr = setStr + "`" + t.fieldNames[fId] + "`" + " = " + toStmtStr(t.fieldValue[fId], t.fieldTypes[fId])
			hasSet = true
		}
	}

	if !hasSet {
		return "", fmt.Errorf("No fields have been set for update!")
	}

	if useRecId {
		whereStr = t.tableName + ".recid = " + anytypeToStr(t.recid.Value)
		hasWhere = true
	} else {
		if whereFields == nil || len(whereFields) == 0 {
			return "", fmt.Errorf("Missing where conditions for update.")
		}

		for _, field := range whereFields {
			if len(whereStr) != 0 {
				whereStr = whereStr + ", "
			}

			whereStr = whereStr + t.tableName + "." + field.FieldName + " = " + toStmtStr(field.Value, t.GetFieldType(field.FieldName))
			hasWhere = true
		}
	}

	if !hasWhere {
		return "", fmt.Errorf("No condictions in the WHERE clause. recid is not used and condictions not passed in.")
	}

	queryStr = queryStr + setStr + " WHERE " + whereStr

	// for debugging
	fmt.Printf("%v\n", queryStr)

	return queryStr, nil
}

// Updates the selected with the values set for fields. Cannot be used for tables that don't have 
// recid. A record must be selected before the DoUpdate can be called. 
func (t *DbTable) DoUpdate(dbcon *DbConnection) error {
	if !t.recid.Exists {
		return fmt.Errorf("This table does not have recid.")
	}

	if t.recid.IsSet {
		return fmt.Errorf("Record must be selected, setting the recid value will not work.")
	}

	if t.recid.Value == 0 {
		return fmt.Errorf("Record has not been selected")
	}

	queryStr, err := t.buildUpdateStr(true, nil)

	if err != nil {
		return err
	}

	rows, err := dbcon.Exec(queryStr)

	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("Something went wrong, %v lines where updated.", rows)
	}

	for fId, isSet := range t.fieldValueSet {
		if isSet {
			t.fieldValueSet[fId] = false
		}
	}

	return nil
}

// This table will update all records that meet the criteria specified by the whereFields slice.
// Values to be updated must be set via the SetFieldValue() method. Updated row count will be returned.
func (t *DbTable) DoUpdateWhere(dbcon *DbConnection, whereFields []DbUpdateField) (int64, error) {
	if whereFields == nil {
		return 0, fmt.Errorf("whereFields value can't be nil.")
	}

	if len(whereFields) == 0 {
		return 0, fmt.Errorf("At least one field must be specified in the where clause.")
	}

	queryStr, err := t.buildUpdateStr(false, whereFields)

	if err != nil {
		return 0, err
	}

	rows, err := dbcon.Exec(queryStr)

	if err != nil {
		return rows, err
	}

	t.ClearFields()

	return rows, nil
}