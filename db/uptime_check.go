package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func GetExpiredUptimeChecks(db *sql.DB) ([]UptimeCheck, error) {
	var expiredUptimeChecks []UptimeCheck
	log.Print("Getting expired uptime checks")
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
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
