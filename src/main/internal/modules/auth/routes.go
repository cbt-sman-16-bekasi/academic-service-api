package auth

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/apikey"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, authHandler *handler.AuthHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	auth := academic.Group("/auth")
	{
		// Public login endpoint (with API key validation)
		auth.POST("/login", apikey.ApiKeyMiddleware(), authHandler.AuthLogin)

		// Protected endpoints
		protected := auth.Group("")
		protected.Use(jwt.AuthMiddleware())
		protected.Use(jwt.SchoolScopeMiddleware())
		{
			protected.POST("/change-password", authHandler.ChangePassword)
			protected.POST("/change-profile", authHandler.ChangeProfile)
		}
	}
}
