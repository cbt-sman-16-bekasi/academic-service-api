package student

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, studentHandler *handler.StudentHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	student := academic.Group("/student")
	student.Use(jwt.AuthMiddleware())
	student.Use(jwt.SchoolScopeMiddleware())
	{
		// Student CRUD
		student.GET("/all",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheStudents, cache.TtlDuration),
			studentHandler.GetAllStudent)
		student.GET("/detail/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "read"),
			studentHandler.GetStudent)
		student.POST("/create",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			studentHandler.CreateStudent)
		student.PUT("/update/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "update"),
			studentHandler.UpdateStudent)
		student.DELETE("/delete/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "delete"),
			studentHandler.DeleteStudent)

		// Template upload/download
		student.GET("/template/download",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			studentHandler.DownloadTemplate)
		student.POST("/template/upload",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			studentHandler.UploadStudent)
	}
}
