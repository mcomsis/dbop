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
	fieldTypes[3] = "DECIMAL"
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
	usersTable.SetFieldValue("name", "testing 5 five")
	usersTable.SetFieldValue("registered", "2012-05-14 09:44:12")
	usersTable.SetFieldValue("role", "5")
	usersTable.SetFieldValue("rating", "5.7")
	usersTable.SetFieldValue("yr", "2015")
	_, err := usersTable.DoInsert(&dbcon)
	
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	
	recid, _ := usersTable.RecId() 
	fmt.Printf("recid = %v\n", recid)
	
	usersTable.SetFieldValue("name", "testing 2 two")
	usersTable.SetFieldValue("registered", "2012-02-14 09:44:12")
	usersTable.SetFieldValue("role", "2")
	usersTable.SetFieldValue("rating", "2.7")
	usersTable.SetFieldValue("yr", "2012")
	_, err = usersTable.DoInsert(&dbcon)
	
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	
	recid, _ = usersTable.RecId() 
	fmt.Printf("recid = %v\n", recid)
	
	usersTable.SetFieldValue("name", "testing 3 three")
	usersTable.SetFieldValue("registered", "2012-03-14 09:44:12")
	usersTable.SetFieldValue("role", "3")
	usersTable.SetFieldValue("rating", "3.7")
	usersTable.SetFieldValue("yr", "2013")
	_, err = usersTable.DoInsert(&dbcon)
	
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	
	recid, _ = usersTable.RecId() 
	fmt.Printf("recid = %v\n", recid)
	
	usersTable.SetFieldValue("name", "testing 4 four")
	usersTable.SetFieldValue("registered", "2012-04-14 09:44:12")
	usersTable.SetFieldValue("role", "4")
	usersTable.SetFieldValue("rating", "4.7")
	usersTable.SetFieldValue("yr", "2014")
	_, err = usersTable.DoInsert(&dbcon)
	
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	
	recid, _ = usersTable.RecId() 
	fmt.Printf("recid = %v\n", recid)
	
	fmt.Printf("done\n")
	*/
	
	/*
	// testing DoDeleteWhere()
	usersTable.SetFieldValue("role", "3")
	i, err := usersTable.DoDeleteWhere(&dbcon)
	
	if err != nil {
		fmt.Printf("err=%v\n",err)
	} else {
		fmt.Printf("lines=%v\n",i)
	}
	*/
	
	// testing DoSelectFirstony()
	//usersTable.SetFieldValue("role","3")
	/*
	usersTable.SetRecId(7) 
	err := usersTable.DoSelectFirstonly(&dbcon)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	
	fmt.Printf("rating = %v\n",usersTable.GetFieldValue("rating"))
	
	usersTable.DoDelete(&dbcon)
	*/
}

