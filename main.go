package main 

import (
		"fmt"		
		"./dbop"
)

const dbString = "test/root/"

func newUsersTable() dbop.DbTable {
	var usersTable dbop.DbTable
	
	fieldNames := make([]string, 6)
	fieldTypes := make([]string, 6)
	
	fieldNames[0] = "id"
	fieldNames[1] = "name"
	fieldNames[2] = "registered"
	fieldNames[3] = "role"
	fieldNames[4] = "rating"
	fieldNames[5] = "yr"
	
	fieldTypes[0] = "BIGINT"
	fieldTypes[1] = "VARCHAR"
	fieldTypes[2] = "DATETIME"
	fieldTypes[3] = "SMALLINT"
	fieldTypes[4] = "FLOAT"
	fieldTypes[5] = "YEAR"
	
	usersTable.InitTable("Users", fieldNames, fieldTypes)
	return usersTable
}

func main() {
	var dbcon dbop.DbConnection
	dbcon.Open(dbString)
	
	usersTable := newUsersTable()
	
	//usersTable.SetFieldValue("id","19")
	//usersTable.SetFieldValue("name","zaraza")
	//usersTable.ClearField("id")
	
	//res := usersTable.DoSelectFirstonly(&dbcon)
	
	/*
	
	if res {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("nok\n") 
	}
	
	
	val := usersTable.GetFieldValue("registered")
	*/
	
	//var usersTableList []dbop.DbTable
	/*
	usersTableList, err := usersTable.DoSelect(&dbcon)
	
	if err != nil {
		panic (err) 
	}
	
	for id := range usersTableList {
		recid := usersTableList[id].GetFieldValue("name")
		fmt.Printf("name = %s, id = %s\n", recid, id)
	} 
	*/
	
	usersTable.SetFieldValue("name", "omgomg")
	r, err := usersTable.DoInsert(&dbcon)
	
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", r)
	}
	
	fmt.Printf("kaut kas\n") 
}

