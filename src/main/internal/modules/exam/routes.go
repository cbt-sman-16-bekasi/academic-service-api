package exam

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/apikey"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers all exam module routes with injected handlers
func RegisterRoutesWithDI(router *gin.RouterGroup, handlers ExamHandlers) {
	examHandler := handlers.ExamHandler
	questionHandler := handlers.QuestionHandler
	bankHandler := handlers.BankHandler
	sessionHandler := handlers.SessionHandler
	typeExamHandler := handlers.TypeExamHandler
	cbtHandler := handlers.CBTHandler

	academic := router.Group("/academic")

	// Recalculate endpoint (special)
	academic.GET("/recalculate", sessionHandler.ExamSessionRecalculate)

	// ========================================
	// Bank Question Routes
	// ========================================
	bank := academic.Group("/bank").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		bank.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheBankQuestion, cache.TtlDuration), bankHandler.GetAllBankQuestion)
		bank.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), bankHandler.GetDetailMasterBankQuestion)
		bank.GET("/detail/bank/subject", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), bankHandler.GetDetailMasterBankQuestionSubject)
		bank.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), bankHandler.CreateMasterBankQuestion)
		bank.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), bankHandler.UpdateMasterBankQuestion)
		bank.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), bankHandler.DeleteMasterBankQuestion)
		bank.GET("/question/:code", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), bankHandler.GetQuestionByBankQuestionCode)
		bank.GET("/question/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), bankHandler.GetQuestionByBankQuestion)
		bank.POST("/question/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), bankHandler.CreateBankQuestion)
		bank.PUT("/question/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), bankHandler.UpdateBankQuestion)
		bank.DELETE("/question/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), bankHandler.DeleteBankQuestion)
		bank.POST("/:masterBankId/question/template/upload", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), bankHandler.UploadBankQuestion)
	}

	// ========================================
	// Exam Routes
	// ========================================
	exam := academic.Group("/exam")
	{
		examRoute := exam.Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
		examRoute.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheExam, cache.TtlDuration), examHandler.GetAllExam)
		examRoute.GET("/member/:examCode", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examHandler.GetExamMember)
		examRoute.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), examHandler.GetDetailExam)
		examRoute.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examHandler.CreateExam)
		examRoute.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), examHandler.UpdateExam)
		examRoute.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), examHandler.DeleteExam)
		examRoute.GET("/:examId/question", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examHandler.GetAllExamQuestion)
		examRoute.POST("/question/bank/add", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), examHandler.AddExamQuestionFromBank)
		examRoute.GET("/:examId/question/template/download", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examHandler.DownloadTemplateQuestion)
		examRoute.POST("/:examId/question/template/upload", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), examHandler.UploadQuestion)
	}

	// ========================================
	// Exam Question Routes
	// ========================================
	examQuestion := exam.Group("/question").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		examQuestion.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), questionHandler.GetDetailExamQuestion)
		examQuestion.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), questionHandler.CreateExamQuestion)
		examQuestion.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), questionHandler.UpdateExamQuestion)
		examQuestion.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), questionHandler.DeleteExamQuestion)
	}

	// ========================================
	// Exam Session Routes
	// ========================================
	examSession := exam.Group("/session").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		examSession.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheExamSession, cache.TtlDuration), sessionHandler.GetAllExamSession)
		examSession.GET("/report", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheExamSessionReport, cache.TtlDuration), sessionHandler.ExamSessionReport)
		examSession.POST("/generate/report", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), sessionHandler.ExamSessionGenerateReport)
		examSession.GET("/answer/student", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), sessionHandler.ExamSessionAnswerResultStudent)
		examSession.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), sessionHandler.GetExamSession)
		examSession.POST("/create", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), sessionHandler.CreateExamSession)
		examSession.POST("/correction/answer/student", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), sessionHandler.ExamSessionAnswerStudentCorrection)
		examSession.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), sessionHandler.UpdateExamSession)
		examSession.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "delete"), sessionHandler.DeleteExamSession)
		examSession.GET("/attendance", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), sessionHandler.GetAttendance)
		examSession.GET("/attendance/download", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), sessionHandler.DownloadAttendance)
		examSession.GET("/member/:sessionId", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), sessionHandler.ExamSessionMember)
		examSession.POST("/reset", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), sessionHandler.ExamSessionReset)
		examSession.POST("/change/score", jwt.RequirePermission([]string{"ADMIN"}, "update"), sessionHandler.ExamSessionCorrectionScore)
		examSession.POST("/cheat/reset", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "update"), sessionHandler.ResetSuspicious)
	}

	// ========================================
	// Exam Session Token Routes
	// ========================================
	examSessionToken := exam.Group("/session/token").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		examSessionToken.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"), sessionHandler.GetAllExamSessionToken)
		examSessionToken.POST("/generate", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "create"), sessionHandler.CreateExamSessionToken)
	}

	// ========================================
	// Type Exam Routes
	// ========================================
	typeExam := exam.Group("/type-exam").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		typeExam.GET("/all", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "list"),
			cache.CacheMiddleware(cache.CacheTypeExam, cache.TtlDuration), typeExamHandler.GetAllTypeExam)
		typeExam.GET("/detail/:id", jwt.RequirePermission([]string{"ADMIN", "TEACHER"}, "read"), typeExamHandler.GetDetailTypeExam)
		typeExam.POST("/create", jwt.RequirePermission([]string{"ADMIN"}, "create"), typeExamHandler.CreateTypeExam)
		typeExam.PUT("/update/:id", jwt.RequirePermission([]string{"ADMIN"}, "update"), typeExamHandler.ModifyTypeExam)
		typeExam.DELETE("/delete/:id", jwt.RequirePermission([]string{"ADMIN"}, "delete"), typeExamHandler.DeleteTypeExam)
	}

	// ========================================
	// CBT Routes (Student-facing)
	// ========================================
	// CBT Login
	router.POST("/auth/cbt/login", apikey.ApiKeyMiddleware(), cbtHandler.AuthCBTLogin)

	// Protected CBT routes
	cbtProtected := router.Group("").Use(jwt.AuthMiddleware(), jwt.SchoolScopeMiddleware())
	{
		cbtProtected.POST("/cbt/token/validate", cbtHandler.ValidateToken)
		cbtProtected.POST("/cbt/retrieve-session", cbtHandler.RetrieveDetailSessionCbt)
		cbtProtected.POST("/cbt/exam/submit", cbtHandler.SubmitExamSession)
		cbtProtected.POST("/cbt/exam/answer/sync", cbtHandler.SyncAnswer)
		cbtProtected.POST("/cbt/exam/answer/latest", cbtHandler.LastAnswer)
		cbtProtected.POST("/cbt/exam/report", cbtHandler.ExamSessionSuspiciousActivityReport)
		cbtProtected.POST("/cbt/session", cbtHandler.SessionInfo)
	}

	// Upload endpoint
	cbtProtected.POST("/upload/base64", cbtHandler.UploadBase64)
}
