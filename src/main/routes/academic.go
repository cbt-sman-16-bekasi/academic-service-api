package routes

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/controllers"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/gin-gonic/gin"
)

func academicRoutes(gr *gin.RouterGroup) {
	schoolController := controllers.NewSchoolController()
	classController := controllers.NewClassController()
	studentController := controllers.NewStudentController()
	teacherController := controllers.NewTeacherController()
	userController := controllers.NewUserController()

	academic := gr.Group("/academic")

	academic.GET("/school", schoolController.GetSchool)
	academic.POST("/auth/login", schoolController.AuthLogin)

	masterAcademic := academic.Use(jwt.AuthMiddleware(), jwt.RequirePermission([]string{"ADMIN"}, "list"))
	{
		masterAcademic.GET("/dashboard", schoolController.GetDashboard)
		masterAcademic.GET("/class-code", schoolController.GetAllClassCode)
		masterAcademic.GET("/subjects", schoolController.GetAllSubject)
	}

	class := academic.Group("/class")
	masterClass := class.Use(jwt.AuthMiddleware())
	{
		masterClass.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), classController.GetAllClass)
		masterClass.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), classController.GetDetailClass)
		masterClass.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), classController.CreateNewClass)
		masterClass.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), classController.UpdateClass)
		masterClass.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), classController.DeleteClass)
	}

	classSubject := class.Group("/subject").Use(jwt.AuthMiddleware())
	{
		classSubject.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), schoolController.GetAllClassSubject)
		classSubject.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), schoolController.GetClassSubject)
		classSubject.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), schoolController.CreateClassSubject)
		classSubject.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), schoolController.ModifyClassSubject)
		classSubject.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), schoolController.DeleteClassSubject)
	}

	student := academic.Group("/student").Use(jwt.AuthMiddleware())
	{
		student.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), studentController.GetAllStudent)
		student.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), studentController.GetStudent)
		student.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), studentController.CreateStudent)
		student.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), studentController.UpdateStudent)
		student.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), studentController.DeleteStudent)
	}

	teacher := academic.Group("/teacher").Use(jwt.AuthMiddleware())
	{
		teacher.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), teacherController.GetAllTeacher)
		teacher.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), teacherController.GetTeacher)
		teacher.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), teacherController.CreateTeacher)
		teacher.POST("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), teacherController.UpdateTeacher)
		teacher.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), teacherController.DeleteTeacher)
	}

	user := academic.Group("/user").Use(jwt.AuthMiddleware())
	{
		user.GET("/roles", jwt.RequirePermission([]string{"ADMIN"}, "list"), userController.GetAllRoles)
	}

}
