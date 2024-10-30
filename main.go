package main

import (
	"janitor/db"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func loadEnv() error {
	log.Print("Loading environment variables")
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}

func init() {
	if err := loadEnv(); err != nil {
		log.Fatal().Err(err).Msg("Error loading environment variables")
	}
}

func main() {
	dbConn, err := db.ConnectDB(os.Getenv("SECRET_XATA_PG_ENDPOINT"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	expiredUptimeChecks, err := db.GetExpiredUptimeChecks(dbConn)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting expired uptime checks")
	}
	log.Printf("Expired uptime checks: %v", len(expiredUptimeChecks))

	defer dbConn.Close()
}
