package routes

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/controllers"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/bucket"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	request2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/server/response"
	"time"
)

func cbtRouter(gr *gin.RouterGroup) {

	schoolController := controllers.NewSchoolController()
	examController := controllers.NewExamController()

	gr.POST("/auth/cbt/login", schoolController.AuthCBTLogin)
	gr.POST("/cbt/token/validate", jwt.AuthMiddleware(), examController.ValidateToken)
	gr.POST("/cbt/retrieve-session", jwt.AuthMiddleware(), examController.RetrieveDetailSessionCbt)
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
}
