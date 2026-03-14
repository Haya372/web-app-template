package vo

type ErrorCode string

const (
	ValidationErrorCode        = ErrorCode("VALIDATION_ERROR")
	InvalidCredentialErrorCode = ErrorCode("INVALID_CREDENTIAL")
	UnauthorizedErrorCode      = ErrorCode("UNAUTHORIZED")
	ForbiddenErrorCode         = ErrorCode("FORBIDDEN")
	InternalErrorCode          = ErrorCode("INTERNAL_ERROR")
	DuplicateEmailErrorCode    = ErrorCode("DUPLICATE_EMAIL")
)

func (c ErrorCode) Title() string {
	switch c {
	case ValidationErrorCode:
		return "validation error"
	case InvalidCredentialErrorCode:
		return "invalid credential"
	case UnauthorizedErrorCode:
		return "Unauthorized"
	case ForbiddenErrorCode:
		return "forbidden"
	case InternalErrorCode:
		return "internal server error"
	case DuplicateEmailErrorCode:
		return "duplicate email"
	default:
		return "application error"
	}
}

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

func NewUnauthorizedError(message string, details map[string]any, err error) error {
	return &baseError{
		status:  401,
		code:    InvalidCredentialErrorCode,
		message: message,
		err:     err,
		details: details,
	}
}

func NewForbiddenError(message string, details map[string]any, err error) error {
	return &baseError{
		status:  403,
		code:    ForbiddenErrorCode,
		message: message,
		err:     err,
		details: details,
	}
}

func NewDuplicateEmailError(err error) error {
	return &baseError{
		status:  409,
		code:    DuplicateEmailErrorCode,
		message: "email already registered",
		err:     err,
	}
}

