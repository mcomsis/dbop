package main 

import (
		"fmt"		
		"./dbop"
)

const dbString = "test/root/"

func newUsersTable() dbop.DbTable {
	var usersTable 	dbop.DbTable
	var recid		[2]bool
	
	fieldNames := make([]string, 5)
	fieldTypes := make([]string, 5)
	
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
	/*
	// testing DoInsert()
	usersTable.SetFieldValue("name", "testing 3 three")
	usersTable.SetFieldValue("registered", "2012-06-22 06:12:12")
	usersTable.SetFieldValue("role", "3")
	usersTable.SetFieldValue("rating", "2.7")
	usersTable.SetFieldValue("yr", "2011")
	usersTable.DoInsert(&dbcon)
	
	recid, _ := usersTable.RecId() 
	fmt.Printf("recid = %v\n", recid)
	*/
	
	// testing DoSelectFirstony()
	usersTable.SetFieldValue("role","3")
	err := usersTable.DoSelectFirstonly(&dbcon)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	
	fmt.Printf("%v\n",usersTable.GetFieldValue("name"))
}

