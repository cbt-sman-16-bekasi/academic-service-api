package exception

// InternalServerExceptionStruct represents an internal server error
type InternalServerExceptionStruct struct {
	ErrorCode int
	Error     string
}

// NewInternalServerExceptionStruct creates a new InternalServerExceptionStruct
func NewInternalServerExceptionStruct(code int, error string) InternalServerExceptionStruct {
	return InternalServerExceptionStruct{
		ErrorCode: code,
		Error:     error,
	}
}

// NewIntenalServerExceptionStruct is backward compatibility alias (typo in original framework)
func NewIntenalServerExceptionStruct(code int, error string) InternalServerExceptionStruct {
	return NewInternalServerExceptionStruct(code, error)
}
