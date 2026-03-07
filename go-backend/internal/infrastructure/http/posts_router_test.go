//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		request      map[string]any
		withAuth     bool
		responseCode int
		problemType  string
	}{
		{
			name:  "Success with valid JWT and content returns 201",
			email: "post-success@example.com",
			request: map[string]any{
				"content": "Hello, world!",
			},
			withAuth:     true,
			responseCode: http.StatusCreated,
		},
		{
			name:  "Missing Authorization header returns 401",
			email: "post-noauth@example.com",
			request: map[string]any{
				"content": "Hello, world!",
			},
			withAuth:     false,
			responseCode: http.StatusUnauthorized,
			problemType:  "UNAUTHORIZED",
		},
		{
			name:  "Empty content returns 400",
			email: "post-empty@example.com",
			request: map[string]any{
				"content": "",
			},
			withAuth:     true,
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
		{
			name:         "Missing content field returns 400",
			email:        "post-missing@example.com",
			request:      map[string]any{},
			withAuth:     true,
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := signupAndGetToken(t, tt.email, "")

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/v1/posts", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			if tt.withAuth {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.responseCode, resp.StatusCode)

			if tt.responseCode == http.StatusCreated {
				payload, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var got struct {
					Id        string `json:"id"`
					UserId    string `json:"userId"`
					Content   string `json:"content"`
					CreatedAt string `json:"createdAt"`
				}
				require.NoError(t, json.Unmarshal(payload, &got))
				assert.NotEmpty(t, got.Id)
				assert.NotEmpty(t, got.UserId)
				assert.Equal(t, "Hello, world!", got.Content)
				assert.NotEmpty(t, got.CreatedAt)
			} else if tt.problemType != "" {
				problem := readProblemResponse(t, resp)
				assert.Equal(t, tt.problemType, problem.Type)
			}

			err = testDb.Cleanup()
			require.NoError(t, err)
		})
	}
}
