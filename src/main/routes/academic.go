package routes

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/controllers"
	"github.com/gin-gonic/gin"
)

func academicRoutes(gr *gin.RouterGroup) {
	schoolController := controllers.NewSchoolController()

	academic := gr.Group("/academic")
	academic.GET("/school", schoolController.GetSchool)
}
