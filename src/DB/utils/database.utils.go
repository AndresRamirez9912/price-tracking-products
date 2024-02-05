package dbUtils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func OpenDBConnection() (*sql.DB, error) {
	connStr := "user=postgres password=45665482 dbname=Price-Tracker host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
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

func CreateTransaction(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error creating the transaction")
		return nil, err
	}
	return tx, nil
}

func CloseTransaction(tx *sql.Tx, err error) {
	if err != nil {
		tx.Rollback()
		log.Println("Rollback made: ", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Commit Failed: ", err)
	}
}
