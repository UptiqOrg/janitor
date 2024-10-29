package db

import (
	"context"
	"database/sql"
	"janitor/tests"
	"log"
	"testing"
	"time"

	. "github.com/franela/goblin"
	"github.com/google/uuid"
)

func InsertDummyUptimeChecks(testDb *sql.DB) error {
	websiteId, err := uuid.Parse("e4505e8c-f83b-42c0-b6ff-dd497899149a")
	if err != nil {
		log.Fatalf("Unable to parse website ID: %v", err)
	}

	uptimeChecks := []UptimeCheck{
		{
			ID:           uuid.New(),
			WebsiteID:    websiteId,
			Status:       "up",
			StatusCode:   200,
			ResponseTime: 120,
		},
		{
			ID:           uuid.New(),
			WebsiteID:    websiteId,
			Status:       "down",
			StatusCode:   404,
			ResponseTime: 250,
		},
	}

	for _, check := range uptimeChecks {
		_, err := testDb.Exec(`
			INSERT INTO uptime_checks (id, website_id, status, status_code, response_time)
			VALUES ($1, $2, $3, $4, $5)`,
			check.ID, check.WebsiteID, check.Status, check.StatusCode, check.ResponseTime)
		if err != nil {
			log.Fatalf("Error inserting values uptime_check: %v", err)
		}
	}

	expiredItem := &UptimeCheck{
		ID:           uuid.New(),
		WebsiteID:    websiteId,
		Status:       "down",
		StatusCode:   500,
		ResponseTime: 300,
		CreatedAt:    time.Now().Add(-10 * 24 * time.Hour),
	}

	_, err = testDb.Exec(`
		INSERT INTO uptime_checks (id, website_id, status, status_code, response_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		expiredItem.ID, expiredItem.WebsiteID, expiredItem.Status, expiredItem.StatusCode, expiredItem.ResponseTime, expiredItem.CreatedAt)
	if err != nil {
		log.Fatalf("Error inserting values uptime_check: %v", err)
	}

	return nil
}

func TestGetExpiredUptimeChecks(t *testing.T) {
	g := Goblin(t)
	ctx := context.Background()

	postgresC, dbConnString, err := tests.Setup(ctx)
	if err != nil {
		log.Fatalf("Failed to setup postgres container: %s", err)
	}

	testDb, err := ConnectDB(dbConnString)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	if err := InsertDummyUptimeChecks(testDb); err != nil {
		log.Fatalf("Error inserting entries in the uptime_checks table: %s", err)
	}

	g.Describe("GetExpiredUptimeChecks", func() {
		g.It("Should return a list of expired uptime checks", func() {
			expiredUptimeChecks, err := GetExpiredUptimeChecks(testDb)
			if err != nil {
				log.Fatalf("Error running test GetExpiredUptimeChecks: %s", err)
			}
			g.Assert(len(expiredUptimeChecks)).Equal(1)
		})
	})

	defer tests.Teardown(ctx, postgresC)
}
