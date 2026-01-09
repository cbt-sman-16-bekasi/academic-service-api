package curriculum

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, subjectHandler *handler.SubjectHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	curriculum := academic.Group("/curriculum")
	curriculum.Use(jwt.AuthMiddleware())
	curriculum.Use(jwt.SchoolScopeMiddleware())
	{
		subject := curriculum.Group("/subject")
		{
			subject.GET("/all",
				jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
				cache.CacheMiddleware(cache.CacheSubjects, cache.TtlDuration),
				subjectHandler.GetAllSubject)
			subject.GET("/detail/:id",
				jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"),
				subjectHandler.GetSubject)
			subject.POST("/create",
				jwt.RequirePermission([]string{"ADMIN"}, "create"),
				subjectHandler.CreateSubject)
			subject.PUT("/update/:id",
				jwt.RequirePermission([]string{"ADMIN"}, "update"),
				subjectHandler.UpdateSubject)
			subject.DELETE("/delete/:id",
				jwt.RequirePermission([]string{"ADMIN"}, "delete"),
				subjectHandler.DeleteSubject)
		}
	}
}
