package middleware

import "github.com/gin-gonic/gin"

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Client-ID, X-Request-ID, X-Correlation-ID, X-Trace-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
