package teacher

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, teacherHandler *handler.TeacherHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	teacher := academic.Group("/teacher")
	teacher.Use(jwt.AuthMiddleware())
	teacher.Use(jwt.SchoolScopeMiddleware())
	{
		// Teacher CRUD
		teacher.GET("/all",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheTeachers, cache.TtlDuration),
			teacherHandler.GetAllTeacher)
		teacher.GET("/detail/:id",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"),
			teacherHandler.GetTeacher)
		teacher.POST("/create",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			teacherHandler.CreateTeacher)
		teacher.PUT("/update/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "update"),
			teacherHandler.UpdateTeacher)
		teacher.DELETE("/delete/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "delete"),
			teacherHandler.DeleteTeacher)

		// Teacher Subject-Class mapping
		teacher.GET("/:teacherId/subject-class",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			teacherHandler.GetTeacherSubjectClassList)
		teacher.GET("/subject-class/detail/:id",
			jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"),
			teacherHandler.GetDetailTeacherSubject)
		teacher.POST("/subject-class/create",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			teacherHandler.CreateTeacherSubject)
		teacher.PUT("/subject-class/update/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "update"),
			teacherHandler.UpdateTeacherSubject)
		teacher.DELETE("/subject-class/delete/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "delete"),
			teacherHandler.DeleteTeacherSubject)
	}
}
