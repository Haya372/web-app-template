package service_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	infra_service "github.com/Haya372/web-app-template/go-backend/internal/infrastructure/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type jwtClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func TestJwtService_GenerateUserAccessToken(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	t.Setenv("AUTH_JWT_TTL_MINUTES", "5")

	svc, err := infra_service.NewJwtService()
	require.NoError(t, err)

	user, err := entity.NewUser("test@example.com", "password", "Test", time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	now := time.Now().UTC()
	token, err := svc.GenerateUserAccessToken(t.Context(), user)
	require.NoError(t, err)
	require.NotNil(t, token)

	parts := strings.Split(token.Value, ".")
	require.Len(t, parts, 3)

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)

	var header jwtHeader

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(headerBytes, &header))
	assert.Equal(t, "HS256", header.Algorithm)
	assert.Equal(t, "JWT", header.Type)

	var claims jwtClaims
	require.NoError(t, json.Unmarshal(payload, &claims))
	assert.Equal(t, user.Id().String(), claims.Subject)
	assert.GreaterOrEqual(t, claims.IssuedAt, now.Add(-time.Second).Unix())
	assert.LessOrEqual(t, claims.IssuedAt, time.Now().UTC().Add(time.Second).Unix())

	exp := time.Unix(claims.ExpiresAt, 0).UTC()
	assert.WithinDuration(t, exp, now.Add(5*time.Minute), 2*time.Second)
	assert.WithinDuration(t, exp, token.ExpiresAt.UTC(), time.Second)

	signingInput := strings.Join(parts[:2], ".")
	mac := hmac.New(sha256.New, []byte("test-secret"))
	_, _ = mac.Write([]byte(signingInput))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	assert.Equal(t, expectedSignature, parts[2])
}

func TestJwtService_NewJwtService_MissingSecret(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "")
	t.Setenv("AUTH_JWT_TTL_MINUTES", "")

	svc, err := infra_service.NewJwtService()
	require.Error(t, err)
	assert.Nil(t, svc)
}

func TestJwtService_NewJwtService_InvalidTTL(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	t.Setenv("AUTH_JWT_TTL_MINUTES", "-1")

	svc, err := infra_service.NewJwtService()
	require.Error(t, err)
	assert.Nil(t, svc)
}
