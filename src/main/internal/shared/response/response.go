package response

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	"github.com/gin-gonic/gin"
)

// Status constants
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFail    = "fail"
)

// Response is the standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ResponseBuilder provides fluent interface for building responses
type ResponseBuilder struct {
	response Response
	ctx      *gin.Context
}

// New creates a new ResponseBuilder
func New(c *gin.Context) *ResponseBuilder {
	return &ResponseBuilder{
		response: Response{},
		ctx:      c,
	}
}

// Success creates a success response builder
func (b *ResponseBuilder) Success(message string) *ResponseBuilder {
	b.response.Status = StatusSuccess
	b.response.Code = Success
	b.response.Message = message
	return b
}

// Error creates an error response builder
func (b *ResponseBuilder) Error(code int, message string) *ResponseBuilder {
	b.response.Status = StatusError
	b.response.Code = code
	b.response.Message = message
	return b
}

// Fail creates a fail response builder (validation failures)
func (b *ResponseBuilder) Fail(message string) *ResponseBuilder {
	b.response.Status = StatusFail
	b.response.Code = UnprocessableEntity
	b.response.Message = message
	return b
}

// WithData sets the response data
func (b *ResponseBuilder) WithData(data interface{}) *ResponseBuilder {
	b.response.Data = data
	return b
}

// WithErrors sets the error details
func (b *ResponseBuilder) WithErrors(errors interface{}) *ResponseBuilder {
	// Hide error details in production
	if config.GetApp().Env == "production" {
		b.response.Errors = nil
	} else {
		b.response.Errors = errors
	}
	return b
}

// WithMeta includes request metadata in response
func (b *ResponseBuilder) WithMeta() *ResponseBuilder {
	if b.ctx != nil {
		b.response.Meta = NewMeta(b.ctx)
	}
	return b
}

// WithCode overrides the response code
func (b *ResponseBuilder) WithCode(code int) *ResponseBuilder {
	b.response.Code = code
	return b
}

// Build returns the Response
func (b *ResponseBuilder) Build() Response {
	return b.response
}

// JSON sends the response as JSON
func (b *ResponseBuilder) JSON() {
	if b.ctx == nil {
		return
	}
	b.ctx.JSON(GetHTTPStatus(b.response.Code), b.response)
}

// JSONAbort sends the response as JSON and aborts
func (b *ResponseBuilder) JSONAbort() {
	if b.ctx == nil {
		return
	}
	b.ctx.JSON(GetHTTPStatus(b.response.Code), b.response)
	b.ctx.Abort()
}

// =============================================================================
// Quick Helper Functions
// =============================================================================

// OK sends a success response with data
func OK(c *gin.Context, message string, data interface{}) {
	New(c).Success(message).WithData(data).WithMeta().JSON()
}

// CreatedResponse sends a 201 created response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	New(c).Success(message).WithData(data).WithCode(Created).WithMeta().JSON()
}

// BadRequestError sends a 400 bad request response
func BadRequestError(c *gin.Context, message string, errors interface{}) {
	New(c).Error(BadRequest, message).WithErrors(errors).WithMeta().JSONAbort()
}

// UnauthorizedError sends a 401 unauthorized response
func UnauthorizedError(c *gin.Context, message string) {
	New(c).Error(Unauthorized, message).WithMeta().JSONAbort()
}

// ForbiddenError sends a 403 forbidden response
func ForbiddenError(c *gin.Context, message string) {
	New(c).Error(Forbidden, message).WithMeta().JSONAbort()
}

// NotFoundError sends a 404 not found response
func NotFoundError(c *gin.Context, message string) {
	New(c).Error(NotFound, message).WithMeta().JSONAbort()
}

// ValidationError sends a 422 validation error response
func ValidationError(c *gin.Context, errors interface{}) {
	New(c).Fail("Validation failed").WithErrors(errors).WithMeta().JSONAbort()
}

// InternalError sends a 500 internal server error response
func InternalError(c *gin.Context, message string, errors interface{}) {
	New(c).Error(ServerError, message).WithErrors(errors).WithMeta().JSONAbort()
}

// ConflictError sends a 409 conflict error response
func ConflictError(c *gin.Context, message string, errors interface{}) {
	New(c).Error(Conflict, message).WithErrors(errors).WithMeta().JSONAbort()
}

// TooManyRequestsError sends a 429 too many requests error response
func TooManyRequestsError(c *gin.Context, message string) {
	New(c).Error(TooManyRequests, message).WithMeta().JSONAbort()
}
