package db_test

import (
	"testing"

	"github.com/Haya372/web-app-template/go-backend/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
)

func TestBuildDSN_UsesDatabaseDSNWhenSet(t *testing.T) {
	t.Setenv("DATABASE_DSN", "postgres://custom:secret@dbhost:5432/mydb?sslmode=require")
	t.Setenv("DB_PORT", "9999") // should be ignored

	dsn := db.BuildDSN()

	assert.Equal(t, "postgres://custom:secret@dbhost:5432/mydb?sslmode=require", dsn)
}

func TestBuildDSN_BuildsFromComponentsWhenDSNNotSet(t *testing.T) {
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")

	dsn := db.BuildDSN()

	assert.Equal(t, "postgres://postgres:postgres@localhost:55432/backend?sslmode=disable", dsn)
}

func TestBuildDSN_UsesCustomDBPort(t *testing.T) {
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "55532")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")

	dsn := db.BuildDSN()

	assert.Equal(t, "postgres://postgres:postgres@localhost:55532/backend?sslmode=disable", dsn)
}

func TestBuildDSN_UsesAllCustomComponents(t *testing.T) {
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("DB_HOST", "myhost")
	t.Setenv("DB_PORT", "5555")
	t.Setenv("DB_USER", "myuser")
	t.Setenv("DB_PASSWORD", "mypass")
	t.Setenv("DB_NAME", "mydb")

	dsn := db.BuildDSN()

	assert.Equal(t, "postgres://myuser:mypass@myhost:5555/mydb?sslmode=disable", dsn)
}
