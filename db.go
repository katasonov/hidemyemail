package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var g_conn_string string

func addEmailToDatabase(uid string, email string) error {
	db, err := sql.Open("mysql", g_conn_string)
	if err != nil {
		return err
	}
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO email VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	_, err = stmtIns.Exec(uid, email)
	if err != nil {
		return err
	}
	return nil
}

func getEmailByUidFromDatabase(uid string) (string, error) {
	email := ""
	db, err := sql.Open("mysql", g_conn_string)
	if err != nil {
		return email, err
	}
	defer db.Close()
	// Prepare statement for inserting data
	stmtSel, err := db.Prepare("SELECT email FROM email WHERE uid=?") // ? = placeholder
	if err != nil {
		return email, err
	}
	defer stmtSel.Close() // Close the statement when we leave main() / the program terminates
	err = stmtSel.QueryRow(uid).Scan(&email)
	if err != nil {
		return email, err
	}
	return email, nil
}

func getUidByEmailFromDatabase(email string) (string, error) {
	uid := ""
	db, err := sql.Open("mysql", g_conn_string)
	if err != nil {
		return uid, err
	}
	defer db.Close()
	// Prepare statement for inserting data
	stmtSel, err := db.Prepare("SELECT uid FROM email WHERE email=?") // ? = placeholder
	if err != nil {
		return uid, err
	}
	defer stmtSel.Close() // Close the statement when we leave main() / the program terminates
	err = stmtSel.QueryRow(email).Scan(&uid)
	if err != nil {
		return uid, err
	}
	return uid, nil
}
