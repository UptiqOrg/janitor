package db

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func DeleteUptimeChecksBatch(db *sql.DB, checks []UptimeCheck) (int64, error) {
	if len(checks) == 0 {
		return 0, nil
	}

	var totalAffected int64
	batchSize := 20

	for i := 0; i < len(checks); i += batchSize {
		end := i + batchSize
		if end > len(checks) {
			end = len(checks)
		}

		var ids []interface{}
		for _, check := range checks[i:end] {
			ids = append(ids, check.ID)
		}

		log.Print("Deleting uptime checks with IDs: ", ids)

		query := `
    DELETE FROM uptime_checks
    WHERE id = ANY($1)
   `
		result, err := db.Exec(query, pq.Array(ids))
		if err != nil {
			return totalAffected, err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return totalAffected, err
		}

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
