//go:build integration

package http_test

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/test/integration"
)

var testDb integration.TestDb
var testServer *httptest.Server

func TestMain(m *testing.M) {
	if err := os.Setenv("AUTH_JWT_SECRET", "test-secret"); err != nil {
		log.Fatalf("failed to set AUTH_JWT_SECRET, err=%v", err)
	}
	if err := os.Setenv("AUTH_JWT_TTL_MINUTES", "60"); err != nil {
		log.Fatalf("failed to set AUTH_JWT_TTL_MINUTES, err=%v", err)
	}

	db, err := integration.NewTestDb(integration.TestDbProps{
		User:      "postgres",
		Password:  "postgres",
		Database:  "repository_it",
		DbDirPath: "../../../db",
	})
	if err != nil {
		log.Fatalf("failed to create db, err=%v", err)
	}

	server, err := integration.InitializeTestServer(context.Background(), db.Pool())
	if err != nil {
		log.Fatalf("failed to start test server, err=%v", err)
	}
	defer db.Terminate()
	defer server.Close()

	testDb = db
	testServer = server

	m.Run()
}
