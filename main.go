package main 

import (
		"fmt"		
		"./dbop"
)

const dbString = "test/root/"

func newUsersTable() dbop.DbTable {
	var usersTable dbop.DbTable
	
	fieldNames := make([]string, 5)
	fieldTypes := make([]string, 5)
	
	fieldNames[0] = "id"
	fieldNames[1] = "name"
	fieldNames[2] = "registered"
	fieldNames[3] = "role"
	fieldNames[4] = "rating"
	
	fieldTypes[0] = "BIGINT"
	fieldTypes[1] = "VARCHAR"
	fieldTypes[2] = "DATETIME"
	fieldTypes[3] = "SMALLINT"
	fieldTypes[4] = "FLOAT"
	
	usersTable.InitTable("Users", fieldNames, fieldTypes)
	return usersTable
}

func main() {
	var dbcon dbop.DbConnection
	dbcon.Open(dbString)
	
	usersTable := newUsersTable()
	
	ret := dbop.DoStuff(&usersTable)
	fmt.Printf("returned %s\n", ret)
	
	rows := dbcon.Exec("insert into Users (name, registered, role, rating) values ('zaraza', '2012-12-10 22:03:27', 1, 1.1)")
	fmt.Printf("affected rows %s\n", rows)
}

