//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		name         string
		request      map[string]any
		responseCode int
	}{
		{
			name: "Success to create user",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "password",
			},
			responseCode: http.StatusCreated,
		},
		{
			name: "Empty password",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "password length under 8 characters",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "passwor",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "Illegal Password",
			request: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "passwor",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "Illegal Email Format",
			request: map[string]any{
				"name":     "Test",
				"email":    "test",
				"password": "password",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "Empty Name",
			request: map[string]any{
				"email":    "test",
				"password": "password",
			},
			responseCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.request)
			if err != nil {
				assert.FailNow(t, "fail to marshal json", err)
			}

			resp, err := http.Post(testServer.URL+"/v1/users/signup", "application/json", bytes.NewBuffer(body))
			if err != nil {
				assert.FailNow(t, "fail to request", err)
			}

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					assert.Fail(t, "failed to close body", err)
				}
			}()

			assert.Equal(t, tt.responseCode, resp.StatusCode)

			if resp.StatusCode == http.StatusCreated {
				payload, err := io.ReadAll(resp.Body)
				if err != nil {
					assert.FailNow(t, "fail to read body", err)
				}

				var got struct {
					Status string `json:"status"`
				}

				require.NoError(t, json.Unmarshal(payload, &got))
				assert.Equal(t, "ACTIVE", got.Status)
			} else {
				problem := readProblemResponse(t, resp)
				assert.Equal(t, tt.responseCode, problem.Status)
				assert.Equal(t, "VALIDATION_ERROR", problem.Type)
			}

			err = testDb.Cleanup()
			if err != nil {
				assert.Fail(t, "fail to cleanup testDb", err)
			}
		})
	}
}

func TestSignup_DuplicateRequest(t *testing.T) {
	request := map[string]any{
		"name":     "Test",
		"email":    "test@example.com",
		"password": "password",
	}

	body, err := json.Marshal(request)
	if err != nil {
		assert.FailNow(t, "fail to marshal json", err)
	}

	resp, err := http.Post(testServer.URL+"/v1/users/signup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		assert.FailNow(t, "fail to request", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			assert.Fail(t, "failed to close body", err)
		}
	}()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp2, err := http.Post(testServer.URL+"/v1/users/signup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		assert.FailNow(t, "fail to request", err)
	}

	defer func() {
		err := resp2.Body.Close()
		if err != nil {
			assert.Fail(t, "failed to close body", err)
		}
	}()

	problem := readProblemResponse(t, resp2)
	assert.Equal(t, http.StatusInternalServerError, problem.Status)
	assert.Equal(t, "INTERNAL_ERROR", problem.Type)

	err = testDb.Cleanup()
	if err != nil {
		assert.Fail(t, "fail to cleanup testDb", err)
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name         string
		request      map[string]any
		responseCode int
	}{
		{
			name: "Success login",
			request: map[string]any{
				"email":    "login@example.com",
				"password": "password",
			},
			responseCode: http.StatusOK,
		},
		{
			name: "User not found",
			request: map[string]any{
				"email":    "unknown@example.com",
				"password": "password",
			},
			responseCode: http.StatusUnauthorized,
		},
		{
			name: "Wrong password",
			request: map[string]any{
				"email":    "login@example.com",
				"password": "wrongpass",
			},
			responseCode: http.StatusUnauthorized,
		},
		{
			name: "Missing password",
			request: map[string]any{
				"email": "login@example.com",
			},
			responseCode: http.StatusBadRequest,
		},
		{
			name: "Invalid email",
			request: map[string]any{
				"email":    "invalid",
				"password": "password",
			},
			responseCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "User not found" {
				signupBody, err := json.Marshal(map[string]any{
					"name":     "Login User",
					"email":    "login@example.com",
					"password": "password",
				})
				if err != nil {
					assert.FailNow(t, "fail to marshal signup json", err)
				}

				signupResp, err := http.Post(testServer.URL+"/v1/users/signup", "application/json", bytes.NewBuffer(signupBody))
				if err != nil {
					assert.FailNow(t, "fail to signup", err)
				}
				_ = signupResp.Body.Close()
			}

			body, err := json.Marshal(tt.request)
			if err != nil {
				assert.FailNow(t, "fail to marshal login json", err)
			}

			resp, err := http.Post(testServer.URL+"/v1/users/login", "application/json", bytes.NewBuffer(body))
			if err != nil {
				assert.FailNow(t, "fail to request", err)
			}

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					assert.Fail(t, "failed to close body", err)
				}
			}()

			assert.Equal(t, tt.responseCode, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				payload, err := io.ReadAll(resp.Body)
				if err != nil {
					assert.FailNow(t, "fail to read body", err)
				}

				var got struct {
					Token     string `json:"token"`
					ExpiresAt string `json:"expiresAt"`
				}

				require.NoError(t, json.Unmarshal(payload, &got))
				assert.NotEmpty(t, got.Token)
				assert.NotEmpty(t, got.ExpiresAt)
			} else {
				problem := readProblemResponse(t, resp)
				assert.Equal(t, tt.responseCode, problem.Status)
				if tt.responseCode == http.StatusBadRequest {
					assert.Equal(t, "VALIDATION_ERROR", problem.Type)
				}
				if tt.responseCode == http.StatusUnauthorized {
					assert.Equal(t, "INVALID_CREDENTIAL", problem.Type)
				}
			}

			err = testDb.Cleanup()
			if err != nil {
				assert.Fail(t, "fail to cleanup testDb", err)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	// Helper: signup + login to obtain a valid JWT.
	getToken := func(t *testing.T, email string) string {
		t.Helper()

		signupBody, err := json.Marshal(map[string]any{
			"name":     "List User",
			"email":    email,
			"password": "password",
		})
		require.NoError(t, err)

		signupResp, err := http.Post(testServer.URL+"/v1/users/signup", "application/json", bytes.NewBuffer(signupBody))
		require.NoError(t, err)
		_ = signupResp.Body.Close()

		loginBody, err := json.Marshal(map[string]any{
			"email":    email,
			"password": "password",
		})
		require.NoError(t, err)

		loginResp, err := http.Post(testServer.URL+"/v1/users/login", "application/json", bytes.NewBuffer(loginBody))
		require.NoError(t, err)
		defer func() { _ = loginResp.Body.Close() }()

		var loginOut struct {
			Token string `json:"token"`
		}
		payload, err := io.ReadAll(loginResp.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(payload, &loginOut))

		return loginOut.Token
	}

	t.Run("Success with valid JWT returns user list", func(t *testing.T) {
		token := getToken(t, "listtest@example.com")

		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		payload, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got struct {
			Users []struct {
				Id        string `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Status    string `json:"status"`
				CreatedAt string `json:"created_at"`
			} `json:"users"`
			Total  int `json:"total"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}
		require.NoError(t, json.Unmarshal(payload, &got))
		assert.GreaterOrEqual(t, got.Total, 1)
		assert.Equal(t, 20, got.Limit)
		assert.Equal(t, 0, got.Offset)

		// Verify password_hash is not exposed.
		assert.NotContains(t, string(payload), "password_hash")

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("Pagination with limit and offset", func(t *testing.T) {
		token := getToken(t, "paginate@example.com")

		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users?limit=5&offset=0", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		payload, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}
		require.NoError(t, json.Unmarshal(payload, &got))
		assert.Equal(t, 5, got.Limit)
		assert.Equal(t, 0, got.Offset)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("Empty user list returns users:[] total:0", func(t *testing.T) {
		token := getToken(t, "empty@example.com")

		err := testDb.Cleanup()
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// After cleanup users table is empty; token is still valid.
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		payload, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got struct {
			Users []any `json:"users"`
			Total int   `json:"total"`
		}
		require.NoError(t, json.Unmarshal(payload, &got))
		assert.Equal(t, 0, got.Total)
		assert.Empty(t, got.Users)
	})

	t.Run("Missing Authorization header returns 401", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users", nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		problem := readProblemResponse(t, resp)
		assert.Equal(t, "UNAUTHORIZED", problem.Type)
	})

	t.Run("Invalid JWT returns 401", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer this.is.invalid")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		problem := readProblemResponse(t, resp)
		assert.Equal(t, "UNAUTHORIZED", problem.Type)
	})

	t.Run("limit=200 out of range returns 400", func(t *testing.T) {
		token := getToken(t, "limitcheck@example.com")

		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users?limit=200", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		problem := readProblemResponse(t, resp)
		assert.Equal(t, "VALIDATION_ERROR", problem.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("limit=0 out of range returns 400", func(t *testing.T) {
		token := getToken(t, "limitcheck2@example.com")

		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/v1/users?limit=0", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		problem := readProblemResponse(t, resp)
		assert.Equal(t, "VALIDATION_ERROR", problem.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})
}

type problemResponse struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}

func readProblemResponse(t *testing.T, resp *http.Response) problemResponse {
	t.Helper()

	assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "application/problem+json"))

	payload, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var got problemResponse
	require.NoError(t, json.Unmarshal(payload, &got))
	assert.NotEmpty(t, got.Type)
	assert.NotEmpty(t, got.Title)

	return got
}
