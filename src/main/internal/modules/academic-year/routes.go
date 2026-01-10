package academicyear

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all academic year routes
func RegisterRoutes(router *gin.Engine, h *handler.AcademicYearHandler) {
	academicYear := router.Group("/academic/academic-year")
	academicYear.Use(jwt.SchoolScopeMiddleware())
	{
		// List and options
		academicYear.GET("", h.GetAll)
		academicYear.GET("/options", h.GetOptions)
		academicYear.GET("/active", h.GetActive)

		// CRUD by ID
		academicYear.GET("/:id", h.GetByID)
		academicYear.PUT("/:id", h.Update)
		academicYear.DELETE("/:id", h.Delete)
		academicYear.POST("/:id/set-active", h.SetActive)

		// By semester ID
		academicYear.GET("/semester/:semesterId", h.GetBySemesterID)

		// Create
		academicYear.POST("", h.Create)
		academicYear.POST("/from-dapodik", h.CreateFromDapodik)
		academicYear.POST("/set-active", h.SetActiveBySemesterID)

		// Student history
		academicYear.GET("/student/:studentId/history", h.GetStudentHistory)
	}
}
