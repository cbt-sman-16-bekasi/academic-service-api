package response

import (
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	reqctx "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/context"
	"github.com/gin-gonic/gin"
)

// Meta contains request metadata for tracing and debugging
type Meta struct {
	CorrelationID string `json:"correlationId"`
	TraceID       string `json:"traceId"`
	RequestID     string `json:"requestId"`
	Timestamp     string `json:"timestamp"`
	Path          string `json:"path"`
	Method        string `json:"method"`
	Latency       string `json:"latency,omitempty"`
	Version       string `json:"version,omitempty"`
}

// NewMeta creates Meta from gin.Context
func NewMeta(c *gin.Context) *Meta {
	rc := reqctx.FromGin(c)

	meta := &Meta{
		CorrelationID: rc.CorrelationID,
		TraceID:       rc.TraceID,
		RequestID:     rc.RequestID,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Path:          rc.Path,
		Method:        rc.Method,
	}

	// Add latency if start time is available
	if !rc.StartTime.IsZero() {
		meta.Latency = rc.Latency().String()
	}

	// Add version from config
	cfg := config.GetApp()
	if cfg != nil {
		meta.Version = "v1.0.0" // Could be from config or build info
	}

	return meta
}

// NewMetaFromContext creates Meta from RequestContext
func NewMetaFromContext(rc *reqctx.RequestContext) *Meta {
	meta := &Meta{
		CorrelationID: rc.CorrelationID,
		TraceID:       rc.TraceID,
		RequestID:     rc.RequestID,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Path:          rc.Path,
		Method:        rc.Method,
	}

	if !rc.StartTime.IsZero() {
		meta.Latency = rc.Latency().String()
	}

	return meta
}
