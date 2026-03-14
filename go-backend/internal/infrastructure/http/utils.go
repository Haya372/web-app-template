package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Haya372/web-app-template/go-backend/internal/domain/vo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

const problemContentType = "application/problem+json"

func writeProblem(c *echo.Context, status int, res problemDetails) error {
	c.Response().Header().Set(echo.HeaderContentType, problemContentType)

	return c.JSON(status, res)
}

// writeJSONError serialises an arbitrary error as a JSON string.
// Used by the generated-handler error callback for request-parse failures.
func writeJSONError(w http.ResponseWriter, err error) error {
	_, werr := w.Write([]byte(`{"type":"VALIDATION_ERROR","title":"validation error","status":400,"detail":"` + err.Error() + `"}`))

	return werr
}

type problemDetails struct {
	Type     vo.ErrorCode `json:"type"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail,omitempty"`
	Instance string       `json:"instance,omitempty"`
	Errors   any          `json:"errors,omitempty"`
}

func handleError(err error) (int, problemDetails) {
	// Validation errors are normalized first so bind/validate failures share one contract.
	validationIssues := collectValidationIssues(err)
	if validationIssues != nil {
		return http.StatusBadRequest, problemDetails{
			Type:   vo.ValidationErrorCode,
			Title:  vo.ValidationErrorCode.Title(),
			Status: http.StatusBadRequest,
			Detail: "invalid request parameters",
			Errors: validationIssues,
		}
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) && httpErr.Code == http.StatusBadRequest {
		return http.StatusBadRequest, problemDetails{
			Type:   vo.ValidationErrorCode,
			Title:  vo.ValidationErrorCode.Title(),
			Status: http.StatusBadRequest,
			Detail: "invalid request parameters",
		}
	}

	var baseErr vo.Error
	// Domain errors are mapped to stable HTTP problem fields via ErrorCode.
	if errors.As(err, &baseErr) {
		return baseErr.Status(), problemDetails{
			Type:   baseErr.Code(),
			Title:  baseErr.Code().Title(),
			Status: baseErr.Status(),
			Detail: baseErr.Message(),
			Errors: baseErr.Details(),
		}
	}

	// Unknown errors are always sanitized to a generic 500 response.
	return http.StatusInternalServerError, problemDetails{
		Type:   vo.InternalErrorCode,
		Title:  vo.InternalErrorCode.Title(),
		Status: http.StatusInternalServerError,
	}
}

func collectValidationIssues(err error) map[string][]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	return buildValidationIssueMap(validationErrors)
}

func buildValidationIssueMap(validationErrors validator.ValidationErrors) map[string][]string {
	issueMap := make(map[string][]string, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := strings.ToLower(fieldError.Field())
		issueMap[field] = append(issueMap[field], messageForValidationRule(fieldError))
	}

	if len(issueMap) == 0 {
		return nil
	}

	return issueMap
}

func messageForValidationRule(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	default:
		return "is invalid"
	}
}
