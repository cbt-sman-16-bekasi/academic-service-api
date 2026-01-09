package response

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewFieldError creates a new field validation error
func NewFieldError(field, code, message string) FieldError {
	return FieldError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}

// NewErrorDetail creates a new error detail
func NewErrorDetail(code, message string) ErrorDetail {
	return ErrorDetail{
		Code:    code,
		Message: message,
	}
}

// WithTarget adds a target to the error detail
func (e ErrorDetail) WithTarget(target string) ErrorDetail {
	e.Target = target
	return e
}

// Common error codes
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeBadRequest    = "BAD_REQUEST"
	ErrCodeServiceError  = "SERVICE_ERROR"
	ErrCodeDatabaseError = "DATABASE_ERROR"
	ErrCodeExternalError = "EXTERNAL_SERVICE_ERROR"
	ErrCodeRateLimited   = "RATE_LIMITED"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodeMissingField  = "MISSING_FIELD"
	ErrCodeInvalidFormat = "INVALID_FORMAT"
	ErrCodeDuplicate     = "DUPLICATE"
	ErrCodeExpired       = "EXPIRED"
)
