package middleware

import (
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	reqctx "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware logs request completion with structured fields
// including correlation ID for distributed tracing
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		rc := reqctx.FromGin(c)
		cfg := config.GetApp()

		// Build log event with correlation IDs
		logEvent := log.Info()
		if len(c.Errors) > 0 {
			logEvent = log.Error()
		}

		logEvent.
			Str("correlationId", rc.CorrelationID).
			Str("traceId", rc.TraceID).
			Str("requestId", rc.RequestID).
			Str("app", cfg.Name).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("ip", c.ClientIP()).
			Str("userAgent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Dur("latency", time.Since(start))

		// Add school context if available
		if rc.SchoolCode != "" {
			logEvent.Str("schoolCode", rc.SchoolCode)
		}
		if rc.UserID > 0 {
			logEvent.Uint("userId", rc.UserID)
		}

		// Add errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logEvent.Str("error", err.Error())
			}
			logEvent.Msg("Request failed")
		} else {
			logEvent.Msg("Request completed")
		}
	}
}
