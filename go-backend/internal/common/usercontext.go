package common

import "context"

type userIDContextKey struct{}

// WithUserId returns a new context carrying the authenticated user's ID.
func WithUserId(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// UserIdFromContext extracts the user ID stored by WithUserId.
// Returns an empty string if no user ID is present.
func UserIdFromContext(ctx context.Context) string {
	id, _ := ctx.Value(userIDContextKey{}).(string)

	return id
}
