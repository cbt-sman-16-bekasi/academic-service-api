package cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheMiddleware untuk cache GET endpoint dengan school scope
// Cache key sekarang include school_code untuk multi-tenant isolation
// Note: Untuk menggunakan school scope, pastikan AuthMiddleware dipanggil sebelum CacheMiddleware
func CacheMiddleware(prefix string, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil school_code dari context untuk cache isolation per sekolah
		// school_code bisa di-set oleh:
		// 1. SchoolScopeMiddleware (recommended)
		// 2. Manual set di controller
		schoolCode := getSchoolCodeFromContext(c)

		// Build cache key dengan school scope
		// Format: prefix:school:school_code:query_params
		key := fmt.Sprintf("%s:school:%s:%s", prefix, schoolCode, c.Request.URL.RawQuery)

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

// getSchoolCodeFromContext extracts school_code from gin context
// It tries multiple sources to find the school code
func getSchoolCodeFromContext(c *gin.Context) string {
	// Try to get from "school_code" key (set by SchoolScopeMiddleware)
	if sc, exists := c.Get("school_code"); exists {
		if scStr, ok := sc.(string); ok && scStr != "" {
			return scStr
		}
	}

	// Try to get from "claims" - extract school_code field using reflection-free approach
	if claims, exists := c.Get("claims"); exists {
		// Claims is a struct with SchoolCode field
		// Use type switch to handle it without importing jwt package
		switch v := claims.(type) {
		case interface{ GetSchoolCode() string }:
			return v.GetSchoolCode()
		default:
			// Try using fmt to extract - claims struct has SchoolCode field
			// This is a fallback for struct without method
			_ = v // avoid unused variable
		}
	}

	return ""
}

// responseCapture captures response body for caching
type responseCapture struct {
	gin.ResponseWriter
	body []byte
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
