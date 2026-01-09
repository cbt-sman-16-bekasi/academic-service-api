package spp

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, paymentHandler *handler.PaymentHandler) {
	spp := router.Group("/spp")
	spp.Use(jwt.AuthMiddleware())
	spp.Use(jwt.SchoolScopeMiddleware())
	{
		// Student payments
		spp.GET("/student/payments", paymentHandler.GetStudentPayments)

		// Admin routes
		spp.GET("/pending", paymentHandler.GetPendingPayments)
		spp.POST("/create", paymentHandler.CreatePayment)
		spp.POST("/:id/pay", paymentHandler.RecordPayment)
		spp.GET("/:id/detail", paymentHandler.GetPaymentDetail)
	}
}
