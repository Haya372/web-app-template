package db

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BuildDSN returns the database DSN to use for connecting.
// In CI/production: DATABASE_DSN is set directly and used as-is.
// In local development: DB_PORT (and optionally DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
// override individual components. Defaults: host=localhost, port=55432, user/pass=postgres, db=backend.
func BuildDSN() string {
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		return dsn
	}

	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "55432")
	user := envOr("DB_USER", "postgres")
	pass := envOr("DB_PASSWORD", "postgres")
	name := envOr("DB_NAME", "backend")

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, net.JoinHostPort(host, port), name)
}

func NewDbPool(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, BuildDSN())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
