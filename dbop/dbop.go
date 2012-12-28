package dbop

import (
		"fmt"
		"database/sql"
		_ "github.com/ziutek/mymysql/godrv"
)

type DbType interface {
	ToStmtStr() string
}

type Table interface {
	InitTable(tableName string, fieldNames []string, fieldTypes []string)
	ResetTable()
	GetTableName() string
	GetFieldNameList() []string
	GetFieldTypeList() []string
	GetFieldValue(fieldName string) string
	SetFieldValue(fieldName string, fieldValue string) bool
}

func ToStmtStr(value string, valueType string) string {
	switch valueType {
		case "BIT", "TINYINT", "BOOL", "BOOLEAN", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "SERIAL", "DECIMAL", "DEC", "FLOAT", "DOUBLE", "YEAR":
			return value
			
		case "DATE","DATETIME", "TIMESTAPM", "TIME", "CHAR", "VARCHAR", "BINARY", "VARBINARY", "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT", "ENUM", "SET":
			return fmt.Sprintf("'%s'",value)
	}
	
	return ""
}


type DbTable struct {
	tableName 		string
	fieldNames		[]string
	fieldTypes		[]string
	fieldValue 		[]string
	fieldValueSet 		[]bool
}

func (t *DbTable) InitTable(tableName string, fieldNames []string, fieldTypes []string) {
	t.tableName = tableName
	t.fieldNames = fieldNames
	t.fieldTypes = fieldTypes
	t.fieldValueSet = make([]bool, len(fieldTypes))
}

func (t *DbTable) ResetTable() {
	t.tableName = ""
	t.fieldNames = nil
	t.fieldTypes = nil
}

func (t DbTable) GetFieldNameList() []string {
	return t.fieldNames
}

func (t DbTable) GetFieldTypeList() []string {
	return t.fieldTypes
}


func (t DbTable) GetTableName() string {
	return t.tableName
}

func (t DbTable) GetFieldValue(fieldName string) string {
	for fId, fn := range t.fieldNames {
		if fieldName == fn {
			return t.fieldValue[fId]
		}
	}
	
	return ""
}

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

func (t *DbTable) ClearFields() {
	for fId = 1; fId <= len(t.fieldValue); fId++ {
		t.fieldValue[fId] = ""
		t.fieldValueSet[fId] = false
	}
}

type DbConnection struct {
	connection *sql.DB
}

func (dbc *DbConnection) Open(connectionStr string) {
	con, err := sql.Open("mymysql", connectionStr)
	
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	
	dbc.connection = con
}

func (dbc *DbConnection) Exec(queryStr string) int64 {
	result, err := dbc.connection.Exec(queryStr)
	
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	
	return rowsAffected
}


func DoStuff(t Table) int {
	return 1
}

