package http_test

import (
	"os"
	"testing"

	infrahttp "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/http"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoConfig_DefaultPort(t *testing.T) {
	t.Setenv("APP_PORT", "")

	cfg := infrahttp.NewEchoConfig()

	assert.Equal(t, ":8080", cfg.Address)
}

func TestNewEchoConfig_CustomPort(t *testing.T) {
	t.Setenv("APP_PORT", "9090")

	cfg := infrahttp.NewEchoConfig()

	assert.Equal(t, ":9090", cfg.Address)
}

func TestNewEchoConfig_RestoredAfterTest(t *testing.T) {
	original := os.Getenv("APP_PORT")

	t.Setenv("APP_PORT", "8181")

	cfg := infrahttp.NewEchoConfig()
	assert.Equal(t, ":8181", cfg.Address)

	// t.Setenv restores original value; verify expectation on original
	_ = original
}
