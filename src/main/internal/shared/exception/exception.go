package exception

import "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"

// ExceptionResponse represents the parsed exception response
type ExceptionResponse struct {
	Code    int
	Message string
	Errors  interface{}
}

// Handle handles panic recovery and returns appropriate response data
func Handle(err any) ExceptionResponse {
	// Check for NotFoundExceptionStruct
	if e, ok := err.(NotFoundExceptionStruct); ok {
		return ExceptionResponse{
			Code:    e.ErrorCode,
			Message: e.Error,
			Errors:  nil,
		}
	}

	// Check for BadRequestExceptionStruct
	if e, ok := err.(BadRequestExceptionStruct); ok {
		return ExceptionResponse{
			Code:    e.ErrorCode,
			Message: e.Error,
			Errors:  nil,
		}
	}

	// Check for InternalServerExceptionStruct
	if e, ok := err.(InternalServerExceptionStruct); ok {
		return ExceptionResponse{
			Code:    e.ErrorCode,
			Message: e.Error,
			Errors:  nil,
		}
	}

	// Check for standard error interface
	if e, ok := err.(error); ok {
		return ExceptionResponse{
			Code:    response.ServerError,
			Message: "Internal server error",
			Errors:  e.Error(),
		}
	}

	// Default fallback
	return ExceptionResponse{
		Code:    response.ServerError,
		Message: "Internal server error",
		Errors:  err,
	}
}
