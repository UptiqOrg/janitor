package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Print("Reading environment variables")

	err := godotenv.Load(".env")
	if err != nil {
		log.Error().Err(err).Msg("error loading .env file")

	}

	dbConnString := os.Getenv("SECRET_XATA_PG_ENDPOINT")
	log.Print(dbConnString)
}

func main() {

}
