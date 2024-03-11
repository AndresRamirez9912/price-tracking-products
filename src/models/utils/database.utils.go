package dbUtils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"price-tracking-products/src/constants"

	_ "github.com/lib/pq"
)

func OpenDBConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable", os.Getenv(constants.DB_USER), os.Getenv(constants.PASSWORD), os.Getenv(constants.DB_NAME), os.Getenv(constants.HOST))
	db, err := sql.Open(constants.DRIVER_NAME, connStr)
	if err != nil {
		log.Println("Error openning the DB connection", err)
		return nil, err
	}
	log.Println("Application connected to the DB")
	err = db.Ping()
	if err != nil {
		log.Println("Error Pinging the DB in openning function", err)
		return nil, err
	}
	log.Println("Successfully Ping to the DB")
	return db, nil
}
