package middleware

import (
	"fmt"
	"runtime/debug"

	reqctx "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/context"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RecoveryMiddleware recovers from panics and returns proper error response
// with correlation ID for tracing
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				rc := reqctx.FromGin(c)
				stack := string(debug.Stack())

				// Log with correlation ID for tracing
				log.Error().
					Str("correlationId", rc.CorrelationID).
					Str("traceId", rc.TraceID).
					Str("requestId", rc.RequestID).
					Str("path", rc.Path).
					Str("method", rc.Method).
					Str("error", fmt.Sprintf("%v", err)).
					Str("stack", stack).
					Msg("Panic recovered")

				// Handle exception and send response
				resp := exception.Handle(err)

				// Send response with meta
				response.New(c).
					Error(resp.Code, resp.Message).
					WithErrors(resp.Errors).
					WithMeta().
					JSONAbort()
			}
		}()
		c.Next()
	}
}
