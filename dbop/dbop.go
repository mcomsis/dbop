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

/*
type DateTime struct {
	day		int
	month	int
	year	int
	hour	int
	minute	int
	second	int
}
*/
/*
func (dateTime DateTime) ToStmtStr() string {
	return fmt.Sprintf("'%i-%i-%i %i:%i:%i'", dateTime.year, dateTime.month, dateTime.day, dateTime.hour, dateTime.minute, dateTime.second) //2012-12-10 22:03:27
}
*/

type DbTable struct {
	tableName 	string
	fieldNames	[]string
	fieldTypes	[]string
	fieldValue 	[]string
}

func (t *DbTable) InitTable(tableName string, fieldNames []string, fieldTypes []string) {
	t.tableName = tableName
	t.fieldNames = fieldNames
	t.fieldTypes = fieldTypes
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
			return true;
		}
	}
	
	return false
}

type DbConnection struct {
	

}

func DoStuff(t Table) int {
	return 1
}

