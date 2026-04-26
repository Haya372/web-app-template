package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Content string

// maxContentLength is the maximum number of Unicode runes allowed in a Content VO.
// The DB column is text (no DB-level length constraint); this limit is a business rule.
const maxContentLength = 10000

var (
	errIllegalContent = errors.New("illegal content")
)

// NewContent validates raw and returns a Content value object.
// Leading/trailing whitespace is trimmed before validation.
// Returns a ValidationError if the trimmed value is empty or exceeds maxContentLength runes.
func NewContent(raw string) (*Content, error) {
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" {
		return nil, NewValidationError("content is required", nil, errIllegalContent)
	}

	if utf8.RuneCountInString(trimmed) > maxContentLength {
		return nil, NewValidationError(
			fmt.Sprintf("content must be at most %d characters long", maxContentLength),
			map[string]any{"max_length": maxContentLength},
			errIllegalContent,
		)
	}

	content := Content(trimmed)

	return &content, nil
}

func (c Content) String() string {
	return string(c)
}
