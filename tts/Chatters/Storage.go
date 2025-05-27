package chatters

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

    _ "github.com/mattn/go-sqlite3"
)

const DBLocation string = "ttsdata.db"

func EstablishDBConnection() *sql.DB {
	if !checkDBExistance() {
		fmt.Println("Creating DB")
		InitDB()
	}
	var db, err = sql.Open("sqlite3", DBLocation)
	if err != nil {
		panic(fmt.Sprintf("Error while opening DB. Original error: %v\n", err.Error()))
	}
	return db
}

func InitDB() {
	var f, createErr = os.Create(DBLocation)
	if createErr != nil {
		panic(fmt.Sprintf("Error while creating DB file. Original error: %v\n", createErr.Error()))
	}
	f.Close()

	var db, err = sql.Open("sqlite3", DBLocation)
	if err != nil {
		panic(fmt.Sprintf("Error while opening DB. Original error: %v\n", err.Error()))
	}
	defer db.Close()

	var _, initErr = db.Exec(DBInitSequence)
	if initErr != nil {
		panic(fmt.Sprintf("Error while creating DB file. Original error: %v\n", initErr.Error()))
	}
}

func checkDBExistance() bool {
	var _, err = os.Stat(DBLocation)
	return !errors.Is(err, os.ErrNotExist)
}

const DBInitSequence string = `
CREATE TABLE Chatters (
	Id INTEGER PRIMARY KEY AUTOINCREMENT, 
	Username TEXT,
	ELO INTEGER
);
`
