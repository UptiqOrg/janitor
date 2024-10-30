package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging database: %w", err)
	}

	return db, nil
}
