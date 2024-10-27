package main

import (
	"database/sql"
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

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}

	defer db.Close()
}
