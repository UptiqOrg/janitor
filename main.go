package main

import (
	"janitor/db"
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

func init() {
	if err := loadEnv(); err != nil {
		log.Fatal().Err(err).Msg("error loading environment variables")
	}
}

func main() {
	db, err := db.ConnectDB(os.Getenv("SECRET_XATA_PG_ENDPOINT"))
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}

	defer db.Close()
}
