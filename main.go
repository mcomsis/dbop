package main 

import (
		"fmt"		
		"./dbop"
		//"time"	
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
	usersTable.SetFieldValue("name", "testing 1 one")
	usersTable.SetFieldValue("registered", "2012-04-13 09:44:12")
	usersTable.SetFieldValue("role", "1")
	usersTable.SetFieldValue("rating", "0.7")
	usersTable.SetFieldValue("yr", "2011")
	_, err := usersTable.DoInsert(&dbcon)
	
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	
	recid, _ := usersTable.RecId() 
	fmt.Printf("recid = %v\n", recid)
	fmt.Printf("done\n")
	*/
	
	/*
	// testing DoDeleteWhere()
	usersTable.SetFieldValue("role", "1")
	i, err := usersTable.DoDeleteWhere(&dbcon)
	
	if err != nil {
		fmt.Printf("err=%v\n",err)
	} else {
		fmt.Printf("lines=%v\n",i)
	}
	*/
	
	// testing DoSelectFirstony()
	usersTable.SetFieldValue("role","3")
	tbls, err := usersTable.DoSelect(&dbcon)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	
	fmt.Printf("%v\n",len(tbls))
		
}

