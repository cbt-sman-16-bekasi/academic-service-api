package dapodik

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, dapodikHandler *handler.DapodikHandler) {
	// Use /academic prefix to maintain backward compatibility
	academic := router.Group("/academic")
	dapodik := academic.Group("/dapodik")
	dapodik.Use(jwt.AuthMiddleware())
	dapodik.Use(jwt.SchoolScopeMiddleware())
	{
		// Config routes
		dapodik.GET("/config",
			jwt.RequirePermission([]string{"ADMIN"}, "read"),
			dapodikHandler.GetDapodikConfig)
		dapodik.POST("/config",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			dapodikHandler.SaveDapodikConfig)

		// Test connection
		dapodik.POST("/test-connection",
			jwt.RequirePermission([]string{"ADMIN"}, "read"),
			dapodikHandler.TestDapodikConnection)

		// Sync routes
		dapodik.POST("/sync",
			jwt.RequirePermission([]string{"ADMIN"}, "create"),
			dapodikHandler.TriggerSync)
		dapodik.GET("/history",
			jwt.RequirePermission([]string{"ADMIN"}, "list"),
			dapodikHandler.GetSyncHistory)
		dapodik.GET("/history/:id",
			jwt.RequirePermission([]string{"ADMIN"}, "read"),
			dapodikHandler.GetSyncHistoryDetail)
		dapodik.GET("/summary",
			jwt.RequirePermission([]string{"ADMIN"}, "read"),
			dapodikHandler.GetSyncSummary)

		// Preview routes - fetch data from Dapodik without syncing
		preview := dapodik.Group("/preview")
		{
			preview.GET("/sekolah",
				jwt.RequirePermission([]string{"ADMIN"}, "read"),
				dapodikHandler.PreviewSekolah)
			preview.GET("/ptk",
				jwt.RequirePermission([]string{"ADMIN"}, "read"),
				dapodikHandler.PreviewPTK)
			preview.GET("/rombel",
				jwt.RequirePermission([]string{"ADMIN"}, "read"),
				dapodikHandler.PreviewRombonganBelajar)
			preview.GET("/peserta-didik",
				jwt.RequirePermission([]string{"ADMIN"}, "read"),
				dapodikHandler.PreviewPesertaDidik)
			preview.GET("/pengguna",
				jwt.RequirePermission([]string{"ADMIN"}, "read"),
				dapodikHandler.PreviewPengguna)
		}
	}
}
