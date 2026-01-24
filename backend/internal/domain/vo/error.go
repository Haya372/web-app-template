package vo

type ErrorCode string

const (
	ValidationErrorCode = ErrorCode("VALIDATION_ERROR")
	InternalErrorCode   = ErrorCode("INTERNAL_ERROR")
)

type Error interface {
	Status() int
	Error() string
	Code() ErrorCode
	Message() string
	Details() map[string]any
}

type baseError struct {
	status  int
	code    ErrorCode
	message string
	err     error
	details map[string]any
}

func (e *baseError) Status() int {
	return e.status
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

func NewValidationError(message string, details map[string]any, err error) error {
	return &baseError{
		status:  400,
		code:    ValidationErrorCode,
		message: message,
		err:     err,
		details: details,
	}
}
