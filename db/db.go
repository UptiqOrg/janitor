package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(connString string) (*sql.DB, error) {
	log.Print("Connecting to database")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
