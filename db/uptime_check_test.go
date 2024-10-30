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

var websiteId, err = uuid.Parse("e4505e8c-f83b-42c0-b6ff-dd497899149a")

func InsertDummyUptimeChecks(testDb *sql.DB) error {
	_, err := testDb.Exec("DELETE FROM uptime_checks")
	if err != nil {
		return err
	}

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

func TestUptimeCheck(t *testing.T) {
	g := Goblin(t)
	ctx := context.Background()

	postgresC, dbConnString, err := tests.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup postgres container: %s", err)
	}

	testDb, err := ConnectDB(dbConnString)
	if err != nil {
		t.Fatalf("Error connecting to database: %s", err)
	}

	g.Describe("GetExpiredUptimeChecks", func() {
		if err := InsertDummyUptimeChecks(testDb); err != nil {
			t.Fatalf("Error inserting entries in the uptime_checks table: %s", err)
		}

		g.It("Should return a list of expired uptime checks", func() {
			expiredUptimeChecks, err := GetExpiredUptimeChecks(testDb)
			if err != nil {
				g.Fail(err)
			}
			g.Assert(len(expiredUptimeChecks)).Equal(1)
			g.Assert(time.Since(expiredUptimeChecks[0].CreatedAt) > 7*24*time.Hour).IsTrue()
		})

		g.It("Should not return error because all inputs are valid", func() {
			_, err := GetExpiredUptimeChecks(testDb)
			if err != nil {
				log.Fatalf("Error running test GetExpiredUptimeChecks: %s", err)
			}
			g.Assert(err).IsNil("should not return an error")
		})

		g.It("Should handle empty results", func() {
			_, err := testDb.Exec("DELETE FROM uptime_checks")
			g.Assert(err).IsNil()

			expiredUptimeChecks, err := GetExpiredUptimeChecks(testDb)
			g.Assert(err).IsNil()
			g.Assert(len(expiredUptimeChecks)).Equal(0)
		})
	})

	g.Describe("DeleteUptimeChecksOneByOne", func() {
		if err := InsertDummyUptimeChecks(testDb); err != nil {
			t.Fatalf("Error inserting entries in the uptime_checks table: %s", err)
		}

		g.It("Should delete the specified uptime checks", func() {
			// Static list of uptime checks to delete
			checksToDelete := []UptimeCheck{
				{
					ID:           uuid.New(),
					WebsiteID:    websiteId,
					Status:       "down",
					StatusCode:   500,
					ResponseTime: 300,
					CreatedAt:    time.Now().Add(-10 * 24 * time.Hour),
				},
			}

			// Insert the static checks into the database
			for _, check := range checksToDelete {
				_, err := testDb.Exec(`
					INSERT INTO uptime_checks (id, website_id, status, status_code, response_time, created_at)
					VALUES ($1, $2, $3, $4, $5, $6)`,
					check.ID, check.WebsiteID, check.Status, check.StatusCode, check.ResponseTime, check.CreatedAt)
				if err != nil {
					g.Fail(err)
				}
			}

			affectedRows, err := DeleteUptimeChecksOneByOne(testDb, checksToDelete)
			if err != nil {
				g.Fail(err)
			}
			g.Assert(affectedRows).Equal(int64(len(checksToDelete)))

			// Verify that the checks have been deleted
			for _, check := range checksToDelete {
				var count int
				err := testDb.QueryRow("SELECT COUNT(*) FROM uptime_checks WHERE id = $1", check.ID).Scan(&count)
				if err != nil {
					g.Fail(err)
				}
				g.Assert(count).Equal(0)
			}
		})

		g.It("Should handle empty input without error", func() {
			affectedRows, err := DeleteUptimeChecksOneByOne(testDb, []UptimeCheck{})
			g.Assert(err).IsNil()
			g.Assert(affectedRows).Equal(int64(0))
		})

		g.It("Should return an error if the database operation fails", func() {
			// Simulate a failure by closing the database connection
			testDb.Close()

			checksToDelete := []UptimeCheck{
				{
					ID:           uuid.New(),
					WebsiteID:    uuid.New(),
					Status:       "down",
					StatusCode:   500,
					ResponseTime: 300,
					CreatedAt:    time.Now().Add(-10 * 24 * time.Hour),
				},
			}

			_, err := DeleteUptimeChecksOneByOne(testDb, checksToDelete)
			g.Assert(err).IsNotNil()
		})
	})

	g.Describe("GetExpiredUptimeChecks and DeleteUptimeChecksOneByOne", func() {
		if err := InsertDummyUptimeChecks(testDb); err != nil {
			t.Fatalf("Error inserting entries in the uptime_checks table: %s", err)
		}
		g.It("Should delete the specified uptime checks", func() {
			expiredUptimeChecks, err := GetExpiredUptimeChecks(testDb)
			if err != nil {
				g.Fail(err)
			}
			g.Assert(len(expiredUptimeChecks)).Equal(1)

			affectedRows, err := DeleteUptimeChecksOneByOne(testDb, expiredUptimeChecks)
			if err != nil {
				g.Fail(err)
			}
			g.Assert(affectedRows).Equal(int64(len(expiredUptimeChecks)))

			remainingChecks, err := GetExpiredUptimeChecks(testDb)
			if err != nil {
				g.Fail(err)
			}
			g.Assert(len(remainingChecks)).Equal(0)
		})

		g.It("Should handle empty input without error", func() {
			affectedRows, err := DeleteUptimeChecksOneByOne(testDb, []UptimeCheck{})
			g.Assert(err).IsNil()
			g.Assert(affectedRows).Equal(int64(0))
		})

		g.It("Should return an error if the database operation fails", func() {
			// Simulate a failure by closing the database connection
			testDb.Close()

			expiredUptimeChecks, err := GetExpiredUptimeChecks(testDb)
			if err != nil {
				g.Fail(err)
			}

			_, err = DeleteUptimeChecksOneByOne(testDb, expiredUptimeChecks)
			g.Assert(err).IsNotNil()
		})
	})

	defer tests.Teardown(ctx, postgresC)
}
