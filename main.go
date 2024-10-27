package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type UptimeCheck struct {
	ID        uuid.UUID
	WebSiteId uuid.UUID
	CreatedAt time.Time
}

func loadEnv() error {
	log.Print("Loading environment variables")
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}

func connectDB() (*sql.DB, error) {
	log.Print("Connecting to database")
	db, err := sql.Open("postgres", os.Getenv("SECRET_XATA_PG_ENDPOINT"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func init() {
	if err := loadEnv(); err != nil {
		log.Fatal().Err(err).Msg("error loading environment variables")
	}
}

	err := godotenv.Load(".env")
	if err != nil {
		log.Error().Err(err).Msg("error loading .env file")

	}

	dbConnString := os.Getenv("SECRET_XATA_PG_ENDPOINT")
	log.Print(dbConnString)
}

func main() {
	db, err := connectDB()

}
