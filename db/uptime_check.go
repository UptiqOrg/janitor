package db

import (
	"database/sql"
	"reflect"
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

func DeleteUptimeChecksOneByOne(db *sql.DB, checks []UptimeCheck) (int64, error) {
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
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	query := `
		SELECT *
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
		reflect.ValueOf(&check).Elem()
		if err := rows.Scan(&check.ID, &check.WebsiteID, &check.Status, &check.StatusCode, &check.ResponseTime, &check.CreatedAt); err != nil {
			log.Error().Err(err).Msg("Error scanning row")
			continue
		}
		log.Print(check)
		expiredUptimeChecks = append(expiredUptimeChecks, check)
	}
	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating rows")
		return nil, err
	}

	return expiredUptimeChecks, nil
}
