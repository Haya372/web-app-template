package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const defaultJWTTTLMinutes = 60

var (
	errMissingJWTSecret = errors.New("AUTH_JWT_SECRET is required")
	errInvalidJWTTTL    = errors.New("AUTH_JWT_TTL_MINUTES must be positive int")
)

type jwtConfig struct {
	secret []byte
	ttl    time.Duration
}

type jwtServiceImpl struct {
	tracer trace.Tracer
	logger common.Logger
	config jwtConfig
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type jwtClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func (g *jwtServiceImpl) GenerateUserAccessToken(
	ctx context.Context,
	user entity.User,
) (*service.UserAccessToken, error) {
	ctx, span := g.tracer.Start(ctx, "Generate")
	defer span.End()

	now := time.Now().UTC()
	expiresAt := now.Add(g.config.ttl)

	header := jwtHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	}
	claims := jwtClaims{
		Subject:   user.Id().String(),
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
	}

	headerSegment, err := encodeJWTSection(header)
	if err != nil {
		g.logger.Error(ctx, "failed to encode jwt header", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	claimsSegment, err := encodeJWTSection(claims)
	if err != nil {
		g.logger.Error(ctx, "failed to encode jwt claims", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	signingInput := fmt.Sprintf("%s.%s", headerSegment, claimsSegment)
	signature := signHS256(g.config.secret, signingInput)

	return &service.UserAccessToken{
		Value:     fmt.Sprintf("%s.%s", signingInput, signature),
		ExpiresAt: expiresAt,
	}, nil
}

func NewJwtService() (service.JwtService, error) {
	config, err := loadJWTConfig()
	if err != nil {
		return nil, err
	}

	return &jwtServiceImpl{
		tracer: otel.Tracer("JwtService"),
		logger: common.NewLogger(),
		config: config,
	}, nil
}

func loadJWTConfig() (jwtConfig, error) {
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		return jwtConfig{}, errMissingJWTSecret
	}

	ttlMinutes := defaultJWTTTLMinutes

	if rawTTL := os.Getenv("AUTH_JWT_TTL_MINUTES"); rawTTL != "" {
		parsed, err := strconv.Atoi(rawTTL)
		if err != nil || parsed <= 0 {
			return jwtConfig{}, fmt.Errorf("%w: got %q", errInvalidJWTTTL, rawTTL)
		}

		ttlMinutes = parsed
	}

	return jwtConfig{
		secret: []byte(secret),
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

func encodeJWTSection(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func signHS256(secret []byte, payload string) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
