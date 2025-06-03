package routes

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/observer"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yon-module/yon-framework/logger"
	"github.com/yon-module/yon-framework/server"
	"os"
)

func init() {
	logger.Log.Info().Msg("[] Yon server running use port " + os.Getenv("yon.server.port"))
	server.AddRoute(func(gr *gin.RouterGroup) {
		gr.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	})

	server.AddRoutes(academicRoutes)
	server.AddRoutes(cbtRouter)
	observer.RegisterEvent()
}
