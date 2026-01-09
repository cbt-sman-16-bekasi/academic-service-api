package exception

// NotFoundExceptionStruct represents a not found error
type NotFoundExceptionStruct struct {
	ErrorCode int
	Error     string
}

// NewNotFoundException creates a new NotFoundExceptionStruct
func NewNotFoundException(code int, error string) NotFoundExceptionStruct {
	return NotFoundExceptionStruct{
		ErrorCode: code,
		Error:     error,
	}
}
