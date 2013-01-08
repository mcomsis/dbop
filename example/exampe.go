package main

import (
	"fmt"
	"github.com/mcomsis/dbop"
)

const dbString = "test/root/"

/*
# Example table. for the purposes of this example recid field is used.
# When creating tables with recid field, it always must be the first one.
CREATE TABLE `Users` (
  `recid` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `registered` datetime NOT NULL,
  `role` smallint(6) DEFAULT NULL,
  `rating` decimal(10,2) DEFAULT NULL,
  `yr` year(4) DEFAULT NULL,
  PRIMARY KEY (`recid`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  UNIQUE KEY `recid_UNIQUE` (`recid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1$$
*/

// this method will initialise the table. can be written differently as long as all the elements are there
func newUsersTable() dbop.DbTable {
	var usersTable dbop.DbTable
	var recid [2]bool

	// lentgh of fieldName and fieldType slices is the number of fields excluding recid if it is used	
	fieldNames := make([]string, 5)
	fieldTypes := make([]string, 5)

	// list of field names as they are in the sql table
	fieldNames[0] = "name"
	fieldNames[1] = "registered"
	fieldNames[2] = "role"
	fieldNames[3] = "rating"
	fieldNames[4] = "yr"

	// list of field types as they are in the sql table. sizes are not specified.
	fieldTypes[0] = "VARCHAR"
	fieldTypes[1] = "DATETIME"
	fieldTypes[2] = "SMALLINT"
	fieldTypes[3] = "DECIMAL"
	fieldTypes[4] = "YEAR"

	// recid field description - must be specified even if recid is not used  
	recid[0] = true // true if recid field is used for this table
	recid[1] = true // true if AUTO_INCREMENT property is set in sql. if recid is handled manually, should be set to false

	// InitTable() called with the value prepared
	usersTable.InitTable("Users", fieldNames, fieldTypes, recid)
	return usersTable
}

func main() {
	var dbcon dbop.DbConnection
	dbcon.Open(dbString)

	usersTable := newUsersTable()

	// Insert is very simple. set all the field values and call DoInsert() 
	fmt.Printf("===== DoInsert()\n")
	usersTable.SetFieldValue("name", "testing 1 one")
	usersTable.SetFieldValue("registered", "2012-05-14 09:44:12")
	usersTable.SetFieldValue("role", "1")
	usersTable.SetFieldValue("rating", "1.7")
	usersTable.SetFieldValue("yr", "2011")
	err := usersTable.DoInsert(&dbcon) // return error if something went wrong.

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}

	recid, _ := usersTable.RecId() // if using recid with AUTO_INCREMENT, the field will be populated after an insert
	// second paremer is a bool - 'true' means the recid value has been set manually
	// false and a recid of not zero means that a record has been selected from db 
	fmt.Printf("recid of the new record = %v\n", recid)

	// let's insert three more rows
	usersTable.SetFieldValue("name", "testing 2 two")
	usersTable.SetFieldValue("registered", "2012-02-14 09:44:12")
	usersTable.SetFieldValue("role", "2")
	usersTable.SetFieldValue("rating", "2.7")
	usersTable.SetFieldValue("yr", "2012")
	err = usersTable.DoInsert(&dbcon)

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}

	recid, _ = usersTable.RecId()
	fmt.Printf("recid of the new record = %v\n", recid)

	usersTable.SetFieldValue("name", "testing 3 three")
	usersTable.SetFieldValue("registered", "2012-03-14 09:44:12")
	usersTable.SetFieldValue("role", "3")
	usersTable.SetFieldValue("rating", "3.7")
	usersTable.SetFieldValue("yr", "2013")
	err = usersTable.DoInsert(&dbcon)

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
	err = usersTable.DoInsert(&dbcon)

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}

	recid, _ = usersTable.RecId()
	fmt.Printf("recid of the new record = %v\n", recid)

	// Selecting records is quock and easy. Set the field value and call DoSelectFirstonly()
	// before selecting a record, let's clear the value from the last insert
	fmt.Printf("===== DoSelectFirstonly()\n")

	usersTable.ClearFields() // wasn't needed on iserts as all the field values were manually set

	usersTable.SetFieldValue("role", "3")      // let's select the record with role=3
	err = usersTable.DoSelectFirstonly(&dbcon) // selects the first row with the set value, populates fields from db

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	fmt.Printf("name field of the selected record = %v\n", usersTable.GetFieldValue("name"))

	// now let's update the yr field to 2011, we will have two records with the same yr field value
	fmt.Printf("===== DoUpdate()\n")
	usersTable.SetFieldValue("yr", "2011")
	err = usersTable.DoUpdate(&dbcon)

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	usersTable.ClearFields()

	// now we can select a bunch of records with DoSelect()
	fmt.Printf("===== DoSelect()\n")
	usersTable.SetFieldValue("yr", "2011")

	rowList, err := usersTable.DoSelect(&dbcon) // returns a slice of all the selected rows

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	for _, tbl := range rowList {
		fmt.Printf("name = %v\n", tbl.GetFieldValue("name"))
	}

	usersTable.ClearFields()

	// now we can do a bulk update. bulk update is a little clumsy, suggestions welcome
	fmt.Printf("===== DoUpdate()\n")
	usersTable.SetFieldValue("rating", "3.3")

	// we have to manually prepare the WHERE list
	var whereField dbop.DbUpdateField
	whereFields := make([]dbop.DbUpdateField, 1)
	whereField.FieldName = "yr"
	whereField.Value = "2011"
	whereFields[0] = whereField

	lines, err := usersTable.DoUpdateWhere(&dbcon, whereFields)

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	fmt.Printf("updated lines - %v\n", lines) // should be two

	usersTable.ClearFields()

	// now let's quickly delete something, DoDelete can only be executed on a previously selected record
	// so we are doing a select first
	fmt.Printf("===== DoDelete()\n")
	usersTable.SetFieldValue("role", "1")
	err = usersTable.DoSelectFirstonly(&dbcon)

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	err = usersTable.DoDelete(&dbcon)

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	usersTable.ClearFields()

	// let's do a bulk delete. does not require a select, will delete all lines that fit
	fmt.Printf("===== DoDeleteWhere()\n")
	usersTable.SetFieldValue("rating", "3.3")
	lines, err = usersTable.DoDeleteWhere(&dbcon)

	if err != nil {
		fmt.Printf("error = %v\n", err)
	}

	fmt.Printf("updated lines - %v\n", lines)
}