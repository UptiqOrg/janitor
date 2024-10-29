package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database %s", err)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
