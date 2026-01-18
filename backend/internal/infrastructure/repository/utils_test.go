//go:build integration

package repository_test

import (
	"log"
	"testing"

	"github.com/Haya372/web-app-template/backend/test/integration"
)

var testDb integration.TestDb

func TestMain(m *testing.M) {
	db, err := integration.NewTestDb(integration.TestDbProps{
		User:     "postgres",
		Password: "postgres",
		Database: "repository_it",
		DdlPath:  "../../../db/schema",
	})
	if err != nil {
		log.Fatalf("failed to create db, err=%v", err)
	}
	defer db.Terminate()

	testDb = db

	m.Run()
}
