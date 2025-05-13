package routes

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/controllers"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/bucket"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	request2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/redisstore"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/server/response"
	"io"
	"time"
)

func academicRoutes(gr *gin.RouterGroup) {
	schoolController := controllers.NewSchoolController()
	classController := controllers.NewClassController()
	studentController := controllers.NewStudentController()
	teacherController := controllers.NewTeacherController()
	userController := controllers.NewUserController()
	examController := controllers.NewExamController()
	curriculumController := controllers.NewCurriculumController()

	academic := gr.Group("/academic")

	academic.GET("/school", schoolController.GetSchool)
	academic.GET("/:bucketName/:folder/:objectName/download", func(c *gin.Context) {
		bucketName := c.Param("bucketName")
		objectName := c.Param("objectName")
		folder := c.Param("folder")

		minioBucket := bucket.NewMinio()
		object, err := minioBucket.RetrieveObject(bucketName, folder+"/"+objectName)
		if err != nil {
			panic(err)
		}
		defer object.Close()

		// Set header agar bisa di-download
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		c.Header("Content-Disposition", "attachment; filename="+objectName)
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		c.Header("Expires", "0")

		// Stream file ke response
		if _, err := io.Copy(c.Writer, object); err != nil {
			response.ErrorResponse(500, "Failed download file", nil)
			return
		}
	})
	academic.PUT("/school/update", jwt.AuthMiddleware(), jwt.RequirePermission([]string{"ADMIN"}, "update"), schoolController.ModifySchool)
	academic.POST("/auth/login", schoolController.AuthLogin)
	academic.GET("/recalculate", examController.ExamSessionRecalculate)
	gr.POST("/auth/cbt/login", schoolController.AuthCBTLogin)
	gr.POST("/cbt/token/validate", jwt.AuthMiddleware(), examController.ValidateToken)
	gr.POST("/auth/change-password", jwt.AuthMiddleware(), schoolController.ChangePassword)
	gr.POST("/auth/change-profile", jwt.AuthMiddleware(), schoolController.ChangeProfile)
	gr.POST("/cbt/exam/submit", jwt.AuthMiddleware(), examController.SubmitExamSession)
	gr.POST("/upload/base64", jwt.AuthMiddleware(), func(context *gin.Context) {
		var request request2.UploadBase64Request
		if err := context.ShouldBindJSON(&request); err != nil {
			panic(err)
		}

		minioCof := bucket.NewMinio()
		info, url := minioCof.UploadViaBase64(request.FileData, time.Now().Format("20060102"))

		response.SuccessResponse("Success", map[string]interface{}{
			"info": info,
			"url":  url,
		}).Json(context)
	})

	curriculumRoute := academic.Group("/curriculum").Use(jwt.AuthMiddleware())
	{
		curriculumRoute.GET("/subject/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheSubjects, redisstore.TtlDuration), curriculumController.GetAllSubject)
		curriculumRoute.GET("/subject/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), curriculumController.GetSubject)
		curriculumRoute.POST("/subject/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), curriculumController.CreateSubject)
		curriculumRoute.PUT("/subject/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), curriculumController.UpdateSubject)
		curriculumRoute.DELETE("/subject/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), curriculumController.DeleteSubject)
	}

	masterAcademic := academic.Use(jwt.AuthMiddleware(), jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"))
	{
		masterAcademic.GET("/dashboard", schoolController.GetDashboard)
		masterAcademic.GET("/class-code", schoolController.GetAllClassCode)
		masterAcademic.GET("/subjects", schoolController.GetAllSubject)
	}

	class := academic.Group("/class")
	masterClass := class.Use(jwt.AuthMiddleware())
	{
		masterClass.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheClass, redisstore.TtlDuration), classController.GetAllClass)
		masterClass.GET("/:classId/member", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), classController.MemberOfClass)
		masterClass.GET("/member/add", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), classController.AddMemberOfClass)
		masterClass.GET("/member/:id/delete", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), classController.DeleteMemberOfClass)
		masterClass.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), classController.GetDetailClass)
		masterClass.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), classController.CreateNewClass)
		masterClass.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), classController.UpdateClass)
		masterClass.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), classController.DeleteClass)
	}

	classSubject := class.Group("/subject").Use(jwt.AuthMiddleware())
	{
		classSubject.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), schoolController.GetAllClassSubject)
		classSubject.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), schoolController.GetClassSubject)
		classSubject.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), schoolController.CreateClassSubject)
		classSubject.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), schoolController.ModifyClassSubject)
		classSubject.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), schoolController.DeleteClassSubject)
	}

	student := academic.Group("/student").Use(jwt.AuthMiddleware())
	{
		student.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheStudents, redisstore.TtlDuration), studentController.GetAllStudent)
		student.GET("/template/download", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), studentController.DownloadTemplate)
		student.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN"}, "read"), studentController.GetStudent)
		student.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), studentController.CreateStudent)
		student.POST("/template/upload", jwt.RequirePermission([]string{"ADMIN"}, "create"), studentController.UploadStudent)
		student.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), studentController.UpdateStudent)
		student.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), studentController.DeleteStudent)
	}

	teacher := academic.Group("/teacher").Use(jwt.AuthMiddleware())
	{
		teacher.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheTeachers, redisstore.TtlDuration), teacherController.GetAllTeacher)
		teacher.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), teacherController.GetTeacher)
		teacher.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), teacherController.CreateTeacher)
		teacher.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), teacherController.UpdateTeacher)
		teacher.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), teacherController.DeleteTeacher)

		teacher.GET("/:teacherId/class-subject/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), teacherController.GetTeacherSubjectClassList)
		teacher.GET("/class-subject/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), teacherController.GetDetailTeacherSubject)
		teacher.POST("/class-subject/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), teacherController.CreateTeacherSubject)
		teacher.PUT("/class-subject/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), teacherController.UpdateTeacherSubject)
		teacher.DELETE("/class-subject/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), teacherController.DeleteTeacherSubject)
	}

	user := academic.Group("/user").Use(jwt.AuthMiddleware())
	{
		user.GET("/roles", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), userController.GetAllRoles)
	}

	bank := academic.Group("/bank").Use(jwt.AuthMiddleware())
	{
		bank.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheBankQuestion, redisstore.TtlDuration), examController.GetAllBankQuestion)
		bank.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetDetailMasterBankQuestion)
		bank.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetDetailMasterBankQuestionSubject)
		bank.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.CreateMasterBankQuestion)
		bank.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), examController.UpdateMasterBankQuestion)
		bank.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), examController.DeleteMasterBankQuestion)
		bank.GET("/question/:code", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.GetQuestionByBankQuestionCode)
		bank.GET("/question/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetQuestionByBankQuestion)
		bank.POST("/question/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.CreateBankQuestion)
		bank.PUT("/question/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), examController.UpdateBankQuestion)
		bank.DELETE("/question/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), examController.DeleteBankQuestion)
		bank.POST("/:masterBankId/question/template/upload", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.UploadBankQuestion)
	}

	exam := academic.Group("/exam")
	{
		examRoute := exam.Use(jwt.AuthMiddleware())
		examRoute.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheExam, redisstore.TtlDuration), examController.GetAllExam)
		examRoute.GET("/member/:examCode", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.GetExamMember)
		examRoute.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetDetailExam)
		examRoute.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.CreateExam)
		examRoute.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), examController.UpdateExam)
		examRoute.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), examController.DeleteExam)
		examRoute.GET("/:examId/question", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.GetAllExamQuestion)
		examRoute.GET("/:examId/question/template/download", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.DownloadTemplateQuestion)
		examRoute.POST("/:examId/question/template/upload", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.UploadQuestion)
	}

	examQuestion := exam.Group("/question").Use(jwt.AuthMiddleware())
	{
		examQuestion.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetDetailExamQuestion)
		examQuestion.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.CreateExamQuestion)
		examQuestion.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), examController.UpdateExamQuestion)
		examQuestion.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), examController.DeleteExamQuestion)
	}

	examSession := exam.Group("/session").Use(jwt.AuthMiddleware())
	{
		examSession.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheExamSession, redisstore.TtlDuration), examController.GetAllExamSession)
		examSession.GET("/report", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheExamSessionReport, redisstore.TtlDuration), examController.ExamSessionReport)
		examSession.GET("/answer/student", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.ExamSessionAnswerResultStudent)
		examSession.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetExamSession)
		examSession.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.CreateExamSession)
		examSession.POST("/correction/answer/student", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.ExamSessionAnswerStudentCorrection)
		examSession.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), examController.UpdateExamSession)
		examSession.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), examController.DeleteExamSession)
		examSession.GET("/attendance", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.GetAttendance)
		examSession.GET("/attendance/download", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.DownloadAttendance)
		examSession.GET("/member/:sessionId", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.ExamSessionMember)
	}

	examSessionToken := exam.Group("/session/token").Use(jwt.AuthMiddleware())
	{
		examSessionToken.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examController.GetAllExamSessionToken)
		examSessionToken.POST("/generate", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examController.CreateExamSessionToken)
	}

	typeExam := exam.Group("/type-exam").Use(jwt.AuthMiddleware())
	{
		typeExam.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			redisstore.CacheMiddleware(redisstore.CacheTypeExam, redisstore.TtlDuration), examController.GetAllTypeExam)
		typeExam.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examController.GetDetailTypeExam)
		typeExam.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), examController.CreateTypeExam)
		typeExam.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), examController.ModifyTypeExam)
		typeExam.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), examController.DeleteTypeExam)
	}
}
