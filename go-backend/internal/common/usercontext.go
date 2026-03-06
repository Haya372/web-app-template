package common

import "context"

type userIDContextKey struct{}

// WithUserID returns a new context carrying the authenticated user's ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// UserIDFromContext extracts the user ID stored by WithUserID.
// Returns an empty string if no user ID is present.
func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(userIDContextKey{}).(string)

	return id
}
