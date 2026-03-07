package vo

// Permission represents a fine-grained access right using the "<resource>:<action>" format.
type Permission string

const (
	PermissionUsersList   Permission = "users:list"
	PermissionUsersCreate Permission = "users:create"
)
