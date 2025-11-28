package errors

import "fmt"

// CustomError represents a structured internal error with a code and message.
type CustomError struct {
	Code    int
	Message string
	Err     error
}

// Error returns the string representation of the error.
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *CustomError) Unwrap() error {
	return e.Err
}

// New creates a new CustomError.
func New(code int, msg string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: msg,
	}
}

// Wrap wraps an existing error into a CustomError.
func Wrap(err error, customErr *CustomError) *CustomError {
	return &CustomError{
		Code:    customErr.Code,
		Message: customErr.Message,
		Err:     err,
	}
}

// Is checks if the target error matches the custom error code.
func Is(err error, target *CustomError) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code == target.Code
	}
	return false
}

// As checks if the error is a CustomError and assigns it to the target.
func As(err error, target **CustomError) bool {
	if customErr, ok := err.(*CustomError); ok {
		*target = customErr
		return true
	}
	return false
}

// General / Internal Errors (10xxxx)
var (
	ErrInternalServer = New(100001, "Internal server error occurred.")
	ErrNotImplemented = New(100002, "Feature not implemented.")
)

// Authentication Errors (20xxxx)
var (
	ErrUnauthorized = New(200001, "User is not authorized.")
	ErrTokenInvalid = New(200002, "Authentication token is invalid.")
	ErrTokenExpired = New(200003, "Authentication token has expired.")
)

// Authorization Errors (30xxxx)
var (
	ErrForbidden               = New(300001, "Access to resource is forbidden.")
	ErrInsufficientPermissions = New(300002, "User has insufficient permissions.")
)

// Validation / Input Errors (40xxxx)
var (
	ErrInvalidInput         = New(400001, "Input parameters are invalid.")
	ErrMissingRequiredField = New(400002, "Required field is missing.")
)

// Infrastructure Errors (50xxxx)
var (
	ErrDatabase           = New(500001, "Database operation failed.")
	ErrRecordNotFound     = New(500002, "Requested record was not found.")
	ErrDuplicateRecord    = New(500003, "Record already exists.")
	ErrNetwork            = New(500004, "Network communication failed.")
	ErrTimeout            = New(500005, "Operation timed out.")
	ErrExternalAPI        = New(500006, "External API returned an error.")
	ErrExternalAPIParsing = New(500007, "Failed to parse response from external API.")
)

// Domain / Business Logic Errors (60xxxx)
var (
	ErrPaperDownload = New(600001, "Failed to download paper.")
)
