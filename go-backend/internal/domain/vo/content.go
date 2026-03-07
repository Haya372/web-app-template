package vo

import "errors"

type Content string

var (
	errIllegalContent = errors.New("illegal content")
)

// NewContent validates raw and returns a Content value object.
// Returns a ValidationError if raw is empty.
func NewContent(raw string) (*Content, error) {
	if raw == "" {
		return nil, NewValidationError("content is required", nil, errIllegalContent)
	}

	content := Content(raw)

	return &content, nil
}
