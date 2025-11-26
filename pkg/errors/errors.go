package errors

import "net/http"

// AppError represents an application error with context
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Predefined error types
var (
	ErrNotFound = &AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}

	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Authentication required",
		Status:  http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:    "FORBIDDEN",
		Message: "Access denied",
		Status:  http.StatusForbidden,
	}

	ErrBadRequest = &AppError{
		Code:    "BAD_REQUEST",
		Message: "Invalid request",
		Status:  http.StatusBadRequest,
	}

	ErrInternal = &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}

	ErrConflict = &AppError{
		Code:    "CONFLICT",
		Message: "Resource already exists",
		Status:  http.StatusConflict,
	}

	ErrRateLimitExceeded = &AppError{
		Code:    "RATE_LIMIT_EXCEEDED",
		Message: "Too many requests",
		Status:  http.StatusTooManyRequests,
	}
)

// New creates a new AppError with a custom message
func New(base *AppError, message string) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: message,
		Status:  base.Status,
	}
}

// Wrap wraps an error with an AppError
func Wrap(base *AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Status:  base.Status,
		Err:     err,
	}
}
