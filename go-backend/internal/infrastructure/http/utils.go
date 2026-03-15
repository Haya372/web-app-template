package http

import (
	"encoding/json"
	"net/http"
)

const problemContentType = "application/problem+json"

// writeJSONError serialises an arbitrary error as a JSON string.
// Used by the generated-handler error callback for request-parse failures.
func writeJSONError(w http.ResponseWriter, err error) error {
	detail, merr := json.Marshal(err.Error())
	if merr != nil {
		detail = []byte(`"internal error"`)
	}

	body := `{"type":"VALIDATION_ERROR","title":"validation error","status":400,"detail":` + string(detail) + `}`
	_, werr := w.Write([]byte(body))

	return werr
}
