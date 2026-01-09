package class

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, classHandler *handler.ClassHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	class := academic.Group("/class")
	class.Use(jwt.AuthMiddleware())
	class.Use(jwt.SchoolScopeMiddleware())
	{
		// Class CRUD
		class.GET("/all",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheClass, cache.TtlDuration),
			classHandler.GetAllClass)
		class.GET("/detail/:id",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"),
			classHandler.GetDetailClass)
		class.POST("/create",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			classHandler.CreateNewClass)
		class.PUT("/update/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "update"),
			classHandler.UpdateClass)
		class.DELETE("/delete/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "delete"),
			classHandler.DeleteClass)

		// Class members
		class.GET("/:classId/member",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			classHandler.MemberOfClass)
		class.POST("/member/add",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			classHandler.AddMemberOfClass)
		class.DELETE("/member/:id/delete",
			jwt.RequirePermission([]string{"ADMIN"}, "delete"),
			classHandler.DeleteMemberOfClass)
		class.DELETE("/batch/delete",
			jwt.RequirePermission([]string{"ADMIN"}, "delete"),
			classHandler.DeleteMemberBatchOfClass)
	}
}
