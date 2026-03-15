package http

import (
	"net/http"
)

const problemContentType = "application/problem+json"

// writeJSONError serialises an arbitrary error as a JSON string.
// Used by the generated-handler error callback for request-parse failures.
func writeJSONError(w http.ResponseWriter, err error) error {
	body := `{"type":"VALIDATION_ERROR","title":"validation error","status":400,"detail":"` + err.Error() + `"}`
	_, werr := w.Write([]byte(body))

	return werr
}
