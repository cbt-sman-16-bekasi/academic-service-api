package routes

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/controllers"
	"github.com/gin-gonic/gin"
)

func academicRoutes(gr *gin.RouterGroup) {
	schoolController := controllers.NewSchoolController()

	academic := gr.Group("/academic")
	academic.GET("/school", schoolController.GetSchool)
	academic.GET("/class-code", schoolController.GetAllClassCode)
	academic.GET("/subjects", schoolController.GetAllSubject)

	class := academic.Group("/class")
	class.GET("/all", schoolController.GetAllClass)
	class.GET("/detail/:id", schoolController.GetDetailClass)
	class.POST("/create", schoolController.CreateNewClass)
	class.PUT("/update/:id", schoolController.UpdateClass)
	class.DELETE("/delete/:id", schoolController.DeleteClass)

	classSubject := class.Group("/subject")
	classSubject.GET("/all", schoolController.GetAllSubject)
	classSubject.GET("/detail/:id", schoolController.GetClassSubject)
	classSubject.POST("/create", schoolController.CreateClassSubject)
	classSubject.PUT("/update/:id", schoolController.ModifyClassSubject)
	classSubject.DELETE("/delete/:id", schoolController.DeleteClassSubject)

	student := academic.Group("/student")
	student.GET("/all", schoolController.GetAllStudent)
	student.GET("/detail/:id", schoolController.GetStudent)
	student.POST("/create", schoolController.CreateStudent)
	student.PUT("/update/:id", schoolController.UpdateStudent)
	student.DELETE("/delete/:id", schoolController.DeleteStudent)
}
