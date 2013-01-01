package main 

import (
		"fmt"		
		"./dbop"
)

const dbString = "test/root/"

func newUsersTable() dbop.DbTable {
	var usersTable 	dbop.DbTable
	var recid		[2]bool
	
	fieldNames := make([]string, 6)
	fieldTypes := make([]string, 6)
	
	fieldNames[0] = "name"
	fieldNames[1] = "registered"
	fieldNames[2] = "role"
	fieldNames[3] = "rating"
	fieldNames[4] = "yr"
	
	fieldTypes[0] = "VARCHAR"
	fieldTypes[1] = "DATETIME"
	fieldTypes[2] = "SMALLINT"
	fieldTypes[3] = "FLOAT"
	fieldTypes[4] = "YEAR"
	
	recid[0] = true
	recid[1] = true
	
	usersTable.InitTable("Users", fieldNames, fieldTypes, recid)
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

