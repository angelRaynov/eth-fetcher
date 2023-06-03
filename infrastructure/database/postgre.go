package database

import (
	"database/sql"
	"fmt"
)

func Init() *sql.DB {
	// Connection parameters
	host := "postgres"
	port := 5432
	user := "postgres"
	password := "postgres"
	dbname := "transaction_data"

	// Create the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//todo: reconnect

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return nil
	}

	return db
}
