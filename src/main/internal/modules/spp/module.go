package spp

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvidePaymentRepository creates repository with injected DB
func ProvidePaymentRepository(db *gorm.DB) *repository.PaymentRepository {
	return repository.NewPaymentRepository(db)
}

// ProvidePaymentService creates service with injected dependencies
func ProvidePaymentService(repo *repository.PaymentRepository) *service.PaymentService {
	return service.NewPaymentService(repo)
}

// ProvidePaymentHandler creates handler with injected service
func ProvidePaymentHandler(svc *service.PaymentService) *handler.PaymentHandler {
	return handler.NewPaymentHandler(svc)
}

// RegisterSppRoutes registers routes with DI
func RegisterSppRoutes(rg *types.RouterGroup, h *handler.PaymentHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all spp module dependencies
var Module = fx.Module("spp",
	fx.Provide(
		ProvidePaymentRepository,
		ProvidePaymentService,
		ProvidePaymentHandler,
	),
	fx.Invoke(RegisterSppRoutes),
)
