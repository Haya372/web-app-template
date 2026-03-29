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

func TestListPosts(t *testing.T) {
	t.Run("Missing Authorization header returns 401", func(t *testing.T) {
		resp, err := newTestClient().GetV1PostsWithResponse(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON401)
		assert.Equal(t, "UNAUTHORIZED", resp.ApplicationproblemJSON401.Type)
	})

	t.Run("Invalid JWT returns 401", func(t *testing.T) {
		resp, err := newTestClient().GetV1PostsWithResponse(
			context.Background(), nil, withBearerToken("this.is.invalid"),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON401)
		assert.Equal(t, "UNAUTHORIZED", resp.ApplicationproblemJSON401.Type)
	})

	t.Run("Authenticated request with no posts returns 200 with empty list", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "listposts-empty@example.com", "")
		c := newTestClient()

		resp, err := c.GetV1PostsWithResponse(context.Background(), nil, withBearerToken(token))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		require.NotNil(t, resp.JSON200)
		assert.Equal(t, 0, resp.JSON200.Total)
		assert.NotNil(t, resp.JSON200.Posts)
		assert.Empty(t, resp.JSON200.Posts)
		assert.Equal(t, 20, resp.JSON200.Limit)
		assert.Equal(t, 0, resp.JSON200.Offset)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("Authenticated request with posts returns 200 with pagination fields echoed", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "listposts-data@example.com", "")
		c := newTestClient()

		// Create a post
		createResp, err := c.PostV1PostsWithResponse(
			context.Background(),
			clientgen.CreatePostRequest{Content: "hello list"},
			withBearerToken(token),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, createResp.StatusCode())

		limit := 5
		offset := 0
		resp, err := c.GetV1PostsWithResponse(
			context.Background(),
			&clientgen.GetV1PostsParams{Limit: &limit, Offset: &offset},
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		require.NotNil(t, resp.JSON200)
		assert.Equal(t, 1, resp.JSON200.Total)
		assert.Len(t, resp.JSON200.Posts, 1)
		assert.Equal(t, limit, resp.JSON200.Limit)
		assert.Equal(t, offset, resp.JSON200.Offset)
		assert.Equal(t, "hello list", resp.JSON200.Posts[0].Content)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("limit=0 returns 400", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "listposts-badlimit@example.com", "")
		c := newTestClient()

		limit := 0
		resp, err := c.GetV1PostsWithResponse(
			context.Background(),
			&clientgen.GetV1PostsParams{Limit: &limit},
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON400)
		assert.Equal(t, "VALIDATION_ERROR", resp.ApplicationproblemJSON400.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("limit=101 above max returns 400", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "listposts-maxlimit@example.com", "")
		c := newTestClient()

		limit := 101
		resp, err := c.GetV1PostsWithResponse(
			context.Background(),
			&clientgen.GetV1PostsParams{Limit: &limit},
			withBearerToken(token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON400)
		assert.Equal(t, "VALIDATION_ERROR", resp.ApplicationproblemJSON400.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("all posts appear regardless of authenticated user", func(t *testing.T) {
		tokenA, _ := signupAndGetToken(t, "multiuser-a@example.com", "")
		tokenB, _ := signupAndGetToken(t, "multiuser-b@example.com", "")
		c := newTestClient()

		// User A creates a post
		createResp, err := c.PostV1PostsWithResponse(
			context.Background(),
			clientgen.CreatePostRequest{Content: "Post by user A"},
			withBearerToken(tokenA),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, createResp.StatusCode())
		postID := createResp.JSON201.Id

		// User B fetches all posts — should include User A's post (no user filtering, Issue #58)
		resp, err := c.GetV1PostsWithResponse(context.Background(), nil, withBearerToken(tokenB))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())
		require.NotNil(t, resp.JSON200)
		require.GreaterOrEqual(t, resp.JSON200.Total, 1)

		found := false
		for _, p := range resp.JSON200.Posts {
			if p.Id == postID {
				found = true
				break
			}
		}
		assert.True(t, found, "User A's post should be visible to User B")

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("pagination with offset returns correct subset", func(t *testing.T) {
		token, _ := signupAndGetToken(t, "offset-user@example.com", "")
		c := newTestClient()

		for range 3 {
			_, err := c.PostV1PostsWithResponse(
				context.Background(),
				clientgen.CreatePostRequest{Content: "post content"},
				withBearerToken(token),
			)
			require.NoError(t, err)
		}

		offset := 2
		limit := 10
		resp, err := c.GetV1PostsWithResponse(
			context.Background(),
			&clientgen.GetV1PostsParams{Offset: &offset, Limit: &limit},
			withBearerToken(token),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())
		require.NotNil(t, resp.JSON200)
		assert.Equal(t, 3, resp.JSON200.Total)
		assert.Len(t, resp.JSON200.Posts, 1)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})
}

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

	t.Run("invalid JWT returns 401", func(t *testing.T) {
		resp, err := newTestClient().PostV1PostsWithResponse(
			context.Background(),
			clientgen.CreatePostRequest{Content: "test"},
			withBearerToken("invalid-token-string"),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		require.NotNil(t, resp.ApplicationproblemJSON401)
		assert.Equal(t, "UNAUTHORIZED", resp.ApplicationproblemJSON401.Type)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})

	t.Run("response includes userId matching authenticated user", func(t *testing.T) {
		token, userID := signupAndGetToken(t, "post-fields@example.com", "")
		c := newTestClient()

		resp, err := c.PostV1PostsWithResponse(
			context.Background(),
			clientgen.CreatePostRequest{Content: "field check"},
			withBearerToken(token),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode())
		require.NotNil(t, resp.JSON201)

		assert.Equal(t, userID, resp.JSON201.UserId.String())
		assert.NotZero(t, resp.JSON201.CreatedAt)

		err = testDb.Cleanup()
		require.NoError(t, err)
	})
}
