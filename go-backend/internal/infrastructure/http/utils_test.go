//go:build integration

package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	stdHTTP "net/http"
	"net/http/httptest"
	"os"
	"testing"

	clientgen "github.com/Haya372/web-app-template/go-backend/test/integration/client/generated"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/Haya372/web-app-template/go-backend/test/integration"
	"github.com/stretchr/testify/require"
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
		Database:  "http_it",
		DbDirPath: "../../../db",
		Schema:    "http_it",
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

// newTestClient returns a ClientWithResponses pointed at the test server.
func newTestClient() *clientgen.ClientWithResponses {
	c, err := clientgen.NewClientWithResponses(testServer.URL)
	if err != nil {
		panic(err)
	}

	return c
}

// withBearerToken returns a RequestEditorFn that adds an Authorization header.
func withBearerToken(token string) clientgen.RequestEditorFn {
	return func(_ context.Context, req *stdHTTP.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)

		return nil
	}
}

// rawPost sends a POST request with a JSON body directly, bypassing typed-client
// email validation. Use this when testing server-side rejection of invalid input
// that the generated client would reject before sending.
func rawPost(t *testing.T, path string, body any) *stdHTTP.Response {
	t.Helper()

	data, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := stdHTTP.NewRequestWithContext(
		context.Background(),
		stdHTTP.MethodPost,
		testServer.URL+path,
		bytes.NewReader(data),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := stdHTTP.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// signupAndGetToken creates a user via signup, logs in, optionally assigns a
// role, and returns the JWT token and user ID.
// Pass an empty roleID to skip role assignment.
func signupAndGetToken(t *testing.T, email, roleID string) (token, userID string) {
	t.Helper()

	ctx := context.Background()
	c := newTestClient()

	signupResp, err := c.PostV1UsersSignupWithResponse(ctx, clientgen.SignupRequest{
		Name:     "Test User",
		Email:    openapi_types.Email(email),
		Password: "password",
	})
	require.NoError(t, err)
	require.Equal(t, 201, signupResp.StatusCode())

	loginResp, err := c.PostV1UsersLoginWithResponse(ctx, clientgen.LoginRequest{
		Email:    openapi_types.Email(email),
		Password: "password",
	})
	require.NoError(t, err)
	require.Equal(t, 200, loginResp.StatusCode())
	require.NotNil(t, loginResp.JSON200)

	token = loginResp.JSON200.Token
	userID = loginResp.JSON200.User.Id

	if roleID != "" {
		assignRoleSQL := "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
		_, execErr := testDb.Pool().Exec(ctx, assignRoleSQL, userID, roleID)
		require.NoError(t, execErr)
	}

	return token, userID
}
