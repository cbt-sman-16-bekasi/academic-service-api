package context

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys
type contextKey string

const (
	KeyCorrelationID contextKey = "correlation_id"
	KeyTraceID       contextKey = "trace_id"
	KeyRequestID     contextKey = "request_id"
	KeyStartTime     contextKey = "start_time"
	KeySchoolCode    contextKey = "school_code"
	KeyUserID        contextKey = "user_id"
	KeyPath          contextKey = "path"
	KeyMethod        contextKey = "method"
)

// Header names
const (
	HeaderCorrelationID = "X-Correlation-ID"
	HeaderTraceID       = "X-Trace-ID"
	HeaderRequestID     = "X-Request-ID"
)

// RequestContext holds all request-scoped data
type RequestContext struct {
	CorrelationID string    `json:"correlationId"`
	TraceID       string    `json:"traceId"`
	RequestID     string    `json:"requestId"`
	StartTime     time.Time `json:"startTime"`
	Path          string    `json:"path"`
	Method        string    `json:"method"`
	SchoolCode    string    `json:"schoolCode,omitempty"`
	UserID        uint      `json:"userId,omitempty"`
}

// NewRequestContext creates a new RequestContext from gin.Context
func NewRequestContext(c *gin.Context) *RequestContext {
	correlationID := c.GetHeader(HeaderCorrelationID)
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	traceID := c.GetHeader(HeaderTraceID)
	if traceID == "" {
		traceID = generateShortID(16)
	}

	requestID := c.GetHeader(HeaderRequestID)
	if requestID == "" {
		requestID = generateShortID(8)
	}

	return &RequestContext{
		CorrelationID: correlationID,
		TraceID:       traceID,
		RequestID:     requestID,
		StartTime:     time.Now(),
		Path:          c.Request.URL.Path,
		Method:        c.Request.Method,
	}
}

// WithSchoolCode sets the school code
func (rc *RequestContext) WithSchoolCode(code string) *RequestContext {
	rc.SchoolCode = code
	return rc
}

// WithUserID sets the user ID
func (rc *RequestContext) WithUserID(id uint) *RequestContext {
	rc.UserID = id
	return rc
}

// Latency returns the time elapsed since request start
func (rc *RequestContext) Latency() time.Duration {
	return time.Since(rc.StartTime)
}

// ToContext adds RequestContext to context.Context
func (rc *RequestContext) ToContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, KeyCorrelationID, rc.CorrelationID)
	ctx = context.WithValue(ctx, KeyTraceID, rc.TraceID)
	ctx = context.WithValue(ctx, KeyRequestID, rc.RequestID)
	ctx = context.WithValue(ctx, KeyStartTime, rc.StartTime)
	ctx = context.WithValue(ctx, KeyPath, rc.Path)
	ctx = context.WithValue(ctx, KeyMethod, rc.Method)
	if rc.SchoolCode != "" {
		ctx = context.WithValue(ctx, KeySchoolCode, rc.SchoolCode)
	}
	if rc.UserID > 0 {
		ctx = context.WithValue(ctx, KeyUserID, rc.UserID)
	}
	return ctx
}

// SetToGin stores RequestContext in gin.Context
func (rc *RequestContext) SetToGin(c *gin.Context) {
	c.Set(string(KeyCorrelationID), rc.CorrelationID)
	c.Set(string(KeyTraceID), rc.TraceID)
	c.Set(string(KeyRequestID), rc.RequestID)
	c.Set(string(KeyStartTime), rc.StartTime)
	c.Set(string(KeyPath), rc.Path)
	c.Set(string(KeyMethod), rc.Method)

	// Set response headers for client tracking
	c.Header(HeaderCorrelationID, rc.CorrelationID)
	c.Header(HeaderTraceID, rc.TraceID)
	c.Header(HeaderRequestID, rc.RequestID)
}

// FromGin retrieves RequestContext from gin.Context
func FromGin(c *gin.Context) *RequestContext {
	rc := &RequestContext{
		CorrelationID: c.GetString(string(KeyCorrelationID)),
		TraceID:       c.GetString(string(KeyTraceID)),
		RequestID:     c.GetString(string(KeyRequestID)),
		Path:          c.Request.URL.Path,
		Method:        c.Request.Method,
	}

	if startTime, exists := c.Get(string(KeyStartTime)); exists {
		if t, ok := startTime.(time.Time); ok {
			rc.StartTime = t
		}
	}

	if schoolCode := c.GetString(string(KeySchoolCode)); schoolCode != "" {
		rc.SchoolCode = schoolCode
	}

	if userID, exists := c.Get(string(KeyUserID)); exists {
		if id, ok := userID.(uint); ok {
			rc.UserID = id
		}
	}

	return rc
}

// GetCorrelationID extracts correlation ID from gin.Context
func GetCorrelationID(c *gin.Context) string {
	return c.GetString(string(KeyCorrelationID))
}

// GetTraceID extracts trace ID from gin.Context
func GetTraceID(c *gin.Context) string {
	return c.GetString(string(KeyTraceID))
}

// GetRequestID extracts request ID from gin.Context
func GetRequestID(c *gin.Context) string {
	return c.GetString(string(KeyRequestID))
}

// generateShortID generates a random hex string
func generateShortID(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}
