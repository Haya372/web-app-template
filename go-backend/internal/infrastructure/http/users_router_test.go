//go:build integration

package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	clientgen "github.com/Haya372/web-app-template/go-backend/test/integration/client/generated"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// adminRoleID is the seeded role ID with full permissions including users:list.
const adminRoleID = "00000000-0000-0000-0000-000000000001"

func TestSignup(t *testing.T) {
	tests := []struct {
		name         string
		request      clientgen.SignupRequest
		responseCode int
		problemType  string
	}{
		{
			name: "Success to create user",
			request: clientgen.SignupRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Password: "password",
			},
			responseCode: http.StatusCreated,
		},
		{
			name: "Empty password",
			request: clientgen.SignupRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Password: "",
			},
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
		{
			name: "password length under 8 characters",
			request: clientgen.SignupRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Password: "passwor",
			},
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
		{
			name: "Empty Name",
			request: clientgen.SignupRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := newTestClient().PostV1UsersSignupWithResponse(context.Background(), tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.responseCode, resp.StatusCode())

			if resp.StatusCode() == http.StatusCreated {
				require.NotNil(t, resp.JSON201)
				assert.Equal(t, "ACTIVE", resp.JSON201.Status)
			} else {
				require.NotNil(t, resp.ApplicationproblemJSON400)
				assert.Equal(t, http.StatusBadRequest, resp.ApplicationproblemJSON400.Status)
				assert.Equal(t, tt.problemType, resp.ApplicationproblemJSON400.Type)
			}

			err = testDb.Cleanup()
			require.NoError(t, err)
		})
	}
}

func TestSignup_IllegalEmailFormat(t *testing.T) {
	// openapi_types.Email rejects invalid emails client-side, so send raw JSON.
	resp := rawPost(t, "/v1/users/signup", map[string]string{
		"name":     "Test",
		"email":    "test",
		"password": "password",
	})
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var problem clientgen.ProblemDetails
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&problem))
	assert.Equal(t, "VALIDATION_ERROR", problem.Type)

	err := testDb.Cleanup()
	require.NoError(t, err)
}

func TestSignup_DuplicateRequest(t *testing.T) {
	c := newTestClient()
	ctx := context.Background()

	request := clientgen.SignupRequest{
		Name:     "Test",
		Email:    openapi_types.Email("test@example.com"),
		Password: "password",
	}

	resp1, err := c.PostV1UsersSignupWithResponse(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp1.StatusCode())

	resp2, err := c.PostV1UsersSignupWithResponse(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp2.StatusCode())
	require.NotNil(t, resp2.ApplicationproblemJSON409)
	assert.Equal(t, "DUPLICATE_EMAIL", resp2.ApplicationproblemJSON409.Type)

	err = testDb.Cleanup()
	require.NoError(t, err)
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name         string
		setup        bool // whether to pre-create the user
		request      clientgen.LoginRequest
		responseCode int
		problemType  string
	}{
		{
			name:  "Success login",
			setup: true,
			request: clientgen.LoginRequest{
				Email:    "login@example.com",
				Password: "password",
			},
			responseCode: http.StatusOK,
		},
		{
			name:  "User not found",
			setup: false,
			request: clientgen.LoginRequest{
				Email:    "unknown@example.com",
				Password: "password",
			},
			responseCode: http.StatusUnauthorized,
			problemType:  "INVALID_CREDENTIAL",
		},
		{
			name:  "Wrong password",
			setup: true,
			request: clientgen.LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpass",
			},
			responseCode: http.StatusUnauthorized,
			problemType:  "INVALID_CREDENTIAL",
		},
		{
			name:  "Missing password",
			setup: true,
			request: clientgen.LoginRequest{
				Email: "login@example.com",
			},
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient()
			ctx := context.Background()

			if tt.setup {
				signupResp, err := c.PostV1UsersSignupWithResponse(ctx, clientgen.SignupRequest{
					Name:     "Login User",
					Email:    "login@example.com",
					Password: "password",
				})
				require.NoError(t, err)
				require.Equal(t, http.StatusCreated, signupResp.StatusCode())
			}

			resp, err := c.PostV1UsersLoginWithResponse(ctx, tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.responseCode, resp.StatusCode())

			if resp.StatusCode() == http.StatusOK {
				require.NotNil(t, resp.JSON200)
				assert.NotEmpty(t, resp.JSON200.Token)
				assert.NotZero(t, resp.JSON200.ExpiresAt)
			} else if tt.responseCode == http.StatusBadRequest {
				require.NotNil(t, resp.ApplicationproblemJSON400)
				assert.Equal(t, tt.problemType, resp.ApplicationproblemJSON400.Type)
			} else if tt.responseCode == http.StatusUnauthorized {
				require.NotNil(t, resp.ApplicationproblemJSON401)
				assert.Equal(t, tt.problemType, resp.ApplicationproblemJSON401.Type)
			}

			err = testDb.Cleanup()
			require.NoError(t, err)
		})
	}
}

func TestLogin_InvalidEmail(t *testing.T) {
	// openapi_types.Email rejects invalid emails client-side, so send raw JSON.
	resp := rawPost(t, "/v1/users/login", map[string]string{
		"email":    "invalid",
		"password": "password",
	})
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var problem clientgen.ProblemDetails
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&problem))
	assert.Equal(t, "VALIDATION_ERROR", problem.Type)

	err := testDb.Cleanup()
	require.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	t.Run("Success with valid JWT and admin role returns user list", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "listtest@example.com", adminRoleID)
		c := newTestClient()

		resp, err := c.GetV1UsersWithResponse(context.Background(), nil, withBearerToken(token))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		require.NotNil(t, resp.JSON200)
		assert.GreaterOrEqual(t, resp.JSON200.Total, 1)
		assert.Equal(t, 20, resp.JSON200.Limit)
		assert.Equal(t, 0, resp.JSON200.Offset)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("Pagination with limit and offset", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "paginate@example.com", adminRoleID)
		c := newTestClient()

		limit := 5
		offset := 0
		resp, err := c.GetV1UsersWithResponse(
			context.Background(),
			&clientgen.GetV1UsersParams{Limit: &limit, Offset: &offset},
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		require.NotNil(t, resp.JSON200)
		assert.Equal(t, 5, resp.JSON200.Limit)
		assert.Equal(t, 0, resp.JSON200.Offset)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("Missing Authorization header returns 401", func(t *testing.T) {
		resp, err := newTestClient().GetV1UsersWithResponse(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON401)
		assert.Equal(t, "UNAUTHORIZED", resp.ApplicationproblemJSON401.Type)
	})

	t.Run("Invalid JWT returns 401", func(t *testing.T) {
		resp, err := newTestClient().GetV1UsersWithResponse(
			context.Background(),
			nil,
			withBearerToken("this.is.invalid"),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON401)
		assert.Equal(t, "UNAUTHORIZED", resp.ApplicationproblemJSON401.Type)
	})

	t.Run("User without any role returns 403", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "norole@example.com", "") // no role assigned

		resp, err := newTestClient().GetV1UsersWithResponse(
			context.Background(),
			nil,
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON403)
		assert.Equal(t, "FORBIDDEN", resp.ApplicationproblemJSON403.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("limit=200 out of range returns 400", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "limitcheck@example.com", adminRoleID)
		c := newTestClient()

		limit := 200
		resp, err := c.GetV1UsersWithResponse(
			context.Background(),
			&clientgen.GetV1UsersParams{Limit: &limit},
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON400)
		assert.Equal(t, "VALIDATION_ERROR", resp.ApplicationproblemJSON400.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("limit=0 out of range returns 400", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "limitcheck2@example.com", adminRoleID)
		c := newTestClient()

		limit := 0
		resp, err := c.GetV1UsersWithResponse(
			context.Background(),
			&clientgen.GetV1UsersParams{Limit: &limit},
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON400)
		assert.Equal(t, "VALIDATION_ERROR", resp.ApplicationproblemJSON400.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})
}
