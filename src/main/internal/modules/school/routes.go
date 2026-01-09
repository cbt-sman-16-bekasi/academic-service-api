package school

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, schoolHandler *handler.SchoolHandler) {
	academic := router.Group("/academic")

	// Public endpoints
	academic.GET("/load/config", cache.CacheMiddleware(cache.CacheConfigSchool, cache.TtlOneDay), schoolHandler.RetrieveConfigSchool)

	// Protected school endpoints
	academic.GET("/school", jwt.AuthMiddleware(), cache.CacheMiddleware(cache.CacheInformationSchool, cache.TtlOneDay), schoolHandler.GetSchool)
	academic.PUT("/school/update", jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware(), jwt.RequirePermission([]string{"ADMIN"}, "update"), schoolHandler.ModifySchool)

	// Master data endpoints
	master := academic.Group("").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware(), jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"))
	{
		master.GET("/dashboard", schoolHandler.GetDashboard)
		master.GET("/class-code", schoolHandler.GetAllClassCode)
		master.GET("/subjects", schoolHandler.GetAllSubject)
	}

	// Class-Subject endpoints
	classSubject := academic.Group("/class/subject").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		classSubject.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), schoolHandler.GetAllClassSubject)
		classSubject.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), schoolHandler.GetClassSubject)
		classSubject.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), schoolHandler.CreateClassSubject)
		classSubject.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), schoolHandler.ModifyClassSubject)
		classSubject.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), schoolHandler.DeleteClassSubject)
	}
}
