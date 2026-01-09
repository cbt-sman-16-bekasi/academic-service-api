package user

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	user := academic.Group("/user")
	user.Use(jwt.AuthMiddleware())
	user.Use(jwt.SchoolScopeMiddleware())
	{
		// Roles
		user.GET("/roles",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			userHandler.GetAllRoles)

		// User CRUD
		user.GET("/all",
			jwt.RequirePermission([]string{"ADMIN"}, "list"),
			userHandler.GetAllUser)
		user.GET("/detail/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "read"),
			userHandler.GetDetailUser)
		user.POST("/create",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			userHandler.CreateUser)
		user.PUT("/update/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "update"),
			userHandler.UpdateUser)
		user.POST("/reset-password/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "update"),
			userHandler.ResetPassword)
	}
}
