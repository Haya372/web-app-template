package vo

type ErrorCode string

const (
	ValidationErrorCode = ErrorCode("VALIDATION_ERROR")
)

type Error interface {
	Error() string
	Code() ErrorCode
	Message() string
	Details() map[string]any
}

type baseError struct {
	code    ErrorCode
	message string
	err     error
	details map[string]any
}

func (e *baseError) Error() string {
	return e.err.Error()
}

func (e *baseError) Code() ErrorCode {
	return e.code
}

func (e *baseError) Message() string {
	return e.message
}

func (e *baseError) Details() map[string]any {
	return e.details
}

func NewValidationError(message string, details map[string]any, err error) Error {
	return &baseError{
		code:    ValidationErrorCode,
		message: message,
		err:     err,
		details: details,
	}
}
