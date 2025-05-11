package redisstore

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CacheMiddleware untuk cache GET endpoint
func CacheMiddleware(prefix string, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		dataClaims := jwt.GetDataClaims(c)
		key := prefix + c.Request.URL.RawQuery + dataClaims.Username
		val, err := Get(key)
		if err == nil && val != "" {
			c.Data(http.StatusOK, "application/json", []byte(val))
			c.Abort()
			return
		}

		// Capture response
		writer := &responseCapture{ResponseWriter: c.Writer, body: []byte{}}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			_ = Set(key, string(writer.body), ttl)
		}
	}
}

// Untuk capture response body
type responseCapture struct {
	gin.ResponseWriter
	body []byte
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
