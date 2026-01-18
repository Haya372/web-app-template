//go:build integration

package repository

import (
	"context"
	"log"
	"testing"

	"github.com/Haya372/go-template/backend/test/integration"
)

var testDb integration.TestDb

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, err := integration.NewTestDb(integration.TestDbProps{
		User:     "postgres",
		Password: "postgres",
		Database: "repository_it",
		DdlPath:  "../../../db/schema",
	})
	if err != nil {
		log.Fatalf("failed to create db, err=%v", err)
	}
	defer db.Container.Terminate(ctx)

	testDb = *db

	m.Run()
}
