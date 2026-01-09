package middleware

import (
	reqctx "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/context"
	"github.com/gin-gonic/gin"
)

// CorrelationMiddleware extracts or generates correlation/trace/request IDs
// and stores them in the context for use throughout the request lifecycle.
//
// Headers processed:
//   - X-Correlation-ID: Client-provided correlation ID (UUID format recommended)
//   - X-Trace-ID: Distributed tracing ID
//   - X-Request-ID: Unique request identifier
//
// If headers are not provided, new IDs are generated automatically.
// All IDs are also set as response headers for client tracking.
func CorrelationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create request context with IDs
		rc := reqctx.NewRequestContext(c)

		// Store in gin context
		rc.SetToGin(c)

		c.Next()
	}
}
