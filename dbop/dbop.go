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

func ToStmtStr(value string, valueType string) string {
	switch valueType {
		case "BIT", "TINYINT", "BOOL", "BOOLEAN", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "SERIAL", "DECIMAL", "DEC", "FLOAT", "DOUBLE", "YEAR":
			return value
			
		case "DATE","DATETIME", "TIMESTAPM", "TIME", "CHAR", "VARCHAR", "BINARY", "VARBINARY", "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT", "ENUM", "SET":
			return fmt.Sprintf("'%s'",value)
	}
	
	return ""
}

func AnytypeToStr(value interface{}) string {
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


type DbTable struct {
	tableName 		string
	fieldNames		[]string
	fieldTypes		[]string
	fieldValue 		[]string
	fieldValueSet   []bool
}

func (t *DbTable) InitTable(tableName string, fieldNames []string, fieldTypes []string) {
	t.tableName 	= tableName
	t.fieldNames 	= fieldNames
	t.fieldTypes 	= fieldTypes
	t.fieldValue 	= make([]string, len(fieldTypes))
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
	for fId := 0; fId < len(t.fieldValue); fId++ {
		t.fieldValue[fId] = ""
		t.fieldValueSet[fId] = false
	}
}

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
			whereStr = whereStr + t.tableName + "." + t.fieldNames[fId] + " = " + ToStmtStr(t.fieldValue[fId], t.fieldTypes[fId])
			
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
		t.fieldValue[fId] = AnytypeToStr(value)
		t.fieldValueSet[fId] = false
	}
	
	return true
}

func (t DbTable) DoSelect(dbc *DbConnection) ([]DbTable, error) {
	var retRows 	[]DbTable
	var counter		int
	
	rows, err := dbc.connection.Query(t.buildSelectStr(false))
	
	if err != nil {
		return nil, err
	}
	
	fields := make([]interface{}, len(t.fieldNames))
	fieldValues := make([]*interface{}, len(t.fieldNames))
	
	for fId := range fields {
		fields[fId] = &fieldValues[fId]
	}
		
	for rows.Next() {
		err := rows.Scan(fields...)		
		
		if err != nil {
			return nil, err
		}
		
		for fId := range fieldValues {
			value := *fieldValues[fId]
			t.fieldValue[fId] = AnytypeToStr(value)
			t.fieldValueSet[fId] = false						
		}
		
		newRetRows := make([]DbTable, len(retRows)+1)
		id := 0
		
		for id < len(retRows) && len(retRows) > 0 {
			newRetRows[id] = retRows[id] 
			id++
		}
		
		newRetRows[id] = t // TODO te kaut kas nenotiek
		
		retRows = newRetRows
		
		counter++
	}
	
	rows.Close()
	
	return retRows, nil	
}
