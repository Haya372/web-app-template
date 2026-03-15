//go:build integration

package http_test

import (
	"context"
	"net/http"
	"testing"

	clientgen "github.com/Haya372/web-app-template/go-backend/test/integration/client/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		request      clientgen.CreatePostRequest
		withAuth     bool
		responseCode int
		problemType  string
	}{
		{
			name:         "Success with valid JWT and content returns 201",
			email:        "post-success@example.com",
			request:      clientgen.CreatePostRequest{Content: "Hello, world!"},
			withAuth:     true,
			responseCode: http.StatusCreated,
		},
		{
			name:         "Missing Authorization header returns 401",
			email:        "post-noauth@example.com",
			request:      clientgen.CreatePostRequest{Content: "Hello, world!"},
			withAuth:     false,
			responseCode: http.StatusUnauthorized,
			problemType:  "UNAUTHORIZED",
		},
		{
			name:         "Empty content returns 400",
			email:        "post-empty@example.com",
			request:      clientgen.CreatePostRequest{Content: ""},
			withAuth:     true,
			responseCode: http.StatusBadRequest,
			problemType:  "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := signupAndGetToken(t, tt.email, "")

			ctx := context.Background()
			c := newTestClient()

			var opts []clientgen.RequestEditorFn
			if tt.withAuth {
				opts = append(opts, withBearerToken(token))
			}

			resp, err := c.PostV1PostsWithResponse(ctx, tt.request, opts...)
			require.NoError(t, err)
			assert.Equal(t, tt.responseCode, resp.StatusCode())

			if resp.StatusCode() == http.StatusCreated {
				require.NotNil(t, resp.JSON201)
				assert.NotEmpty(t, resp.JSON201.Id)
				assert.NotEmpty(t, resp.JSON201.UserId)
				assert.Equal(t, "Hello, world!", resp.JSON201.Content)
				assert.NotZero(t, resp.JSON201.CreatedAt)
			} else if tt.responseCode == http.StatusUnauthorized {
				require.NotNil(t, resp.ApplicationproblemJSON401)
				assert.Equal(t, tt.problemType, resp.ApplicationproblemJSON401.Type)
			} else if tt.responseCode == http.StatusBadRequest {
				require.NotNil(t, resp.ApplicationproblemJSON400)
				assert.Equal(t, tt.problemType, resp.ApplicationproblemJSON400.Type)
			}

			err = testDb.Cleanup()
			require.NoError(t, err)
		})
	}
}
