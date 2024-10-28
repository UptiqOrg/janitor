package main

import (
	"janitor/db"
	"janitor/tests"
	"testing"

	. "github.com/franela/goblin"
)

func TestMain(t *testing.T) {
	g := Goblin(t)

	postgresC, dbConnString, err := tests.Setup(t)
	if err != nil {
		t.Fatalf("failed to setup postgres container: %s", err)
		return
	}

	db.ConnectDB(dbConnString)

	g.Describe("getExpiredUptimeChecks", func() {
		g.It("should return a list of expired uptime checks", func() {
			g.Assert(1).Equal(1)
		})
	})

	tests.Teardown(&postgresC)
}
