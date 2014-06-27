package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func addEmailToDatabase(uid string, email string) error {
	db, err := sql.Open("mysql", "hidemyemail:Avk241083@/hidemyemaildb")
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
	db, err := sql.Open("mysql", "hidemyemail:Avk241083@/hidemyemaildb")
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

