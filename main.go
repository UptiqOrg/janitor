package main

import (
	"database/sql"
	"janitor/db"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type UptimeCheck struct {
	ID           uuid.UUID
	WebsiteID    uuid.UUID
	Status       string
	StatusCode   int
	ResponseTime int
	CreatedAt    time.Time
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
		log.Fatal().Err(err).Msg("Error loading environment variables")
	}
}

func deleteUptimeChecksOneByOne(db *sql.DB, checks []UptimeCheck) (int64, error) {
	var totalAffected int64
	for _, check := range checks {
		result, err := db.Exec(`
            DELETE FROM uptime_checks
            WHERE id = $1
        `, check.ID)
		if err != nil {
			return totalAffected, err
		}
		affected, _ := result.RowsAffected()
		totalAffected += affected
	}
	return totalAffected, nil
}

func GetExpiredUptimeChecks(db *sql.DB) ([]UptimeCheck, error) {
	var expiredUptimeChecks []UptimeCheck
	log.Print("Getting expired uptime checks")
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	query := `
		SELECT id, created_at
		FROM uptime_checks
		WHERE created_at <= $1
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, sevenDaysAgo)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching id and created_at")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var check UptimeCheck
		err := rows.Scan(&check.ID, &check.CreatedAt)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning row")
		}
		expiredUptimeChecks = append(expiredUptimeChecks, check)
	}
	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating rows")
		return nil, err
	}

	return expiredUptimeChecks, nil
}

func main() {
	db, err := db.ConnectDB(os.Getenv("SECRET_XATA_PG_ENDPOINT"))

	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	expiredUptimeChecks, err := GetExpiredUptimeChecks(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting expired uptime checks")
	}
	log.Printf("Expired uptime checks: %v", expiredUptimeChecks)

	defer db.Close()
}
