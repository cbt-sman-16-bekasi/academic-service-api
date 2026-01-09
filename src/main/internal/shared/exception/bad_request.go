package exception

// BadRequestExceptionStruct represents a bad request error
type BadRequestExceptionStruct struct {
	ErrorCode int
	Error     string
}

// NewBadRequestExceptionStruct creates a new BadRequestExceptionStruct
func NewBadRequestExceptionStruct(code int, error string) BadRequestExceptionStruct {
	return BadRequestExceptionStruct{
		ErrorCode: code,
		Error:     error,
	}
}
