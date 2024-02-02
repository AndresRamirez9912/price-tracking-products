package dbUtils

import (
	"database/sql"
	"log"
)

func OpenDBConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", "45665482")
	if err != nil {
		log.Println("Error openning the DB connection", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error Pinging the DB in openning function", err)
		return nil, err
	}
	return db, nil
}

func CloseDBConnection(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Println("Error Pinging the DB in closing function", err)
		return
	}

	db.Close()
}
