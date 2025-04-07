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
	examController := controllers.NewExamController()

	academic := gr.Group("/academic")

	academic.GET("/school", schoolController.GetSchool)
	academic.PUT("/school/update", jwt.AuthMiddleware(), jwt.RequirePermission([]string{"ADMIN"}, "update"), schoolController.ModifySchool)
	academic.POST("/auth/login", schoolController.AuthLogin)
	gr.POST("/auth/cbt/login", schoolController.AuthCBTLogin)
	gr.POST("/cbt/token/validate", jwt.AuthMiddleware(), examController.ValidateToken)
	gr.POST("/cbt/exam/submit", jwt.AuthMiddleware(), examController.SubmitExamSession)

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
		teacher.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), teacherController.UpdateTeacher)
		teacher.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), teacherController.DeleteTeacher)
	}

	user := academic.Group("/user").Use(jwt.AuthMiddleware())
	{
		user.GET("/roles", jwt.RequirePermission([]string{"ADMIN"}, "list"), userController.GetAllRoles)
	}

	bank := academic.Group("/bank").Use(jwt.AuthMiddleware())
	{
		bank.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAllBankQuestion)
		bank.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), examController.GetDetailMasterBankQuestion)
		bank.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateMasterBankQuestion)
		bank.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.UpdateMasterBankQuestion)
		bank.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteMasterBankQuestion)
		bank.GET("/question/:code", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetQuestionByBankQuestionCode)
		bank.GET("/question/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), examController.GetQuestionByBankQuestion)
		bank.POST("/question/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateBankQuestion)
		bank.PUT("/question/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.UpdateBankQuestion)
		bank.DELETE("/question/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteBankQuestion)
	}

	exam := academic.Group("/exam")
	{
		examRoute := exam.Use(jwt.AuthMiddleware())
		examRoute.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAllExam)
		examRoute.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), examController.GetDetailExam)
		examRoute.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateExam)
		examRoute.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.UpdateExam)
		examRoute.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteExam)
		examRoute.GET("/:examId/question", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAllExamQuestion)
		examRoute.GET("/:examId/question/template/download", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.DownloadTemplateQuestion)
		examRoute.POST("/:examId/question/template/upload", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.UploadQuestion)
	}

	examQuestion := exam.Group("/question").Use(jwt.AuthMiddleware())
	{
		examQuestion.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), examController.GetDetailExamQuestion)
		examQuestion.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateExamQuestion)
		examQuestion.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.UpdateExamQuestion)
		examQuestion.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteExamQuestion)
	}

	examSession := exam.Group("/session").Use(jwt.AuthMiddleware())
	{
		examSession.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAllExamSession)
		examSession.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), examController.GetExamSession)
		examSession.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateExamSession)
		examSession.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.UpdateExamSession)
		examSession.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteExamSession)
		examSession.GET("/attendance", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAttendance)
		examSession.GET("/attendance/download", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.DownloadAttendance)
	}

	examSessionToken := exam.Group("/session/token").Use(jwt.AuthMiddleware())
	{
		examSessionToken.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAllExamSessionToken)
		examSessionToken.POST("/generate", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateExamSessionToken)
	}

	typeExam := exam.Group("/type-exam").Use(jwt.AuthMiddleware())
	{
		typeExam.GET("/all", jwt.RequirePermission([]string{"ADMIN"}, "list"), examController.GetAllTypeExam)
		typeExam.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), examController.GetDetailTypeExam)
		typeExam.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateTypeExam)
		typeExam.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.ModifyTypeExam)
		typeExam.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteTypeExam)
	}
}
