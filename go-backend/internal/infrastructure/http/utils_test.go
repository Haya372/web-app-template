//go:build integration

package http_test

import (
	"context"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/Haya372/web-app-template/go-backend/test/integration"
)

var testDb integration.TestDb
var testServer *httptest.Server

func TestMain(m *testing.M) {
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
