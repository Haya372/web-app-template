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

type errorImpl struct {
	code    ErrorCode
	message string
	err     error
	details map[string]any
}

func (e *errorImpl) Error() string {
	return e.err.Error()
}

func (e *errorImpl) Code() ErrorCode {
	return e.code
}

func (e *errorImpl) Message() string {
	return e.message
}

func (e *errorImpl) Details() map[string]any {
	return e.details
}

func NewValidationError(message string, details map[string]any, err error) Error {
	return &errorImpl{
		code:    ValidationErrorCode,
		message: message,
		err:     err,
		details: details,
	}
}
