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
	"strings"
	"time"

	"github.com/Haya372/web-app-template/go-backend/internal/common"
	"github.com/Haya372/web-app-template/go-backend/internal/domain/entity"
	"github.com/Haya372/web-app-template/go-backend/internal/usecase/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultJWTTTLMinutes = 60
	jwtPartsCount        = 3
)

var (
	errMissingJWTSecret = errors.New("AUTH_JWT_SECRET is required")
	errInvalidJWTTTL    = errors.New("AUTH_JWT_TTL_MINUTES must be positive int")
	errInvalidJWTFormat = errors.New("invalid JWT format")
	errInvalidSignature = errors.New("invalid JWT signature")
	errTokenExpired     = errors.New("JWT token has expired")
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

func (g *jwtServiceImpl) ValidateToken(ctx context.Context, token string) (*service.TokenClaims, error) {
	_, span := g.tracer.Start(ctx, "ValidateToken")
	defer span.End()

	parts := strings.Split(token, ".")
	if len(parts) != jwtPartsCount {
		span.RecordError(errInvalidJWTFormat)
		span.SetStatus(codes.Error, errInvalidJWTFormat.Error())

		return nil, errInvalidJWTFormat
	}

	signingInput := fmt.Sprintf("%s.%s", parts[0], parts[1])
	expectedSig := signHS256(g.config.secret, signingInput)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		span.RecordError(errInvalidSignature)
		span.SetStatus(codes.Error, errInvalidSignature.Error())

		return nil, errInvalidSignature
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("decode claims: %w", err)
	}

	var claims jwtClaims
	if err = json.Unmarshal(claimsJSON, &claims); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	if time.Now().UTC().Unix() > claims.ExpiresAt {
		span.RecordError(errTokenExpired)
		span.SetStatus(codes.Error, errTokenExpired.Error())

		return nil, errTokenExpired
	}

	return &service.TokenClaims{UserId: claims.Subject}, nil
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
