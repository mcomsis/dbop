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
/*
type zzint interface {
	AddBtoA()
}

type zz struct {
	A int
	B int
}

func (z *zz) AddBtoA () {
	z.A = z.A + z.B
}

func doall(zi zzint) {
	zi.AddBtoA()
}

func main() {
	var z zz
	z.A = 1
	z.B = 2
	
	doall(z)
	
	fmt.Printf("%i",z.A)
*/
func main() {
	usersTable := newUsersTable()
	ret := dbop.DoStuff(&usersTable)
	fmt.Printf("returned %i\n", ret)
	con, err := sql.Open("mymysql", dbString)
	
	if err != nil {
		fmt.Printf("Error str %s", err)
	}
	
	_, err = con.Exec("insert into Users (name, registered, role, rating) values ('zaraza', '2012-12-10 22:03:27', 1, 1.1)")
	
	//fmt.Printf("package %n", i)
	
}

