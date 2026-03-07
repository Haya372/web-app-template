//go:build integration

package reader_test

import (
	"log"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/test/integration"
)

var testDb integration.TestDb

func TestMain(m *testing.M) {
	db, err := integration.NewTestDb(integration.TestDbProps{
		User:      "postgres",
		Password:  "postgres",
		Database:  "reader_it",
		DbDirPath: "../../../db",
		Schema:    "reader_it",
	})
	if err != nil {
		log.Fatalf("failed to create db, err=%v", err)
	}

	defer db.Terminate()

	testDb = db

	m.Run()
}
