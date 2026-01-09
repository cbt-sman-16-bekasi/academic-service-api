package auth

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideAuthRepository creates repository with injected DB
func ProvideAuthRepository(db *gorm.DB) *repository.AuthRepository {
	return repository.NewAuthRepository(db)
}

// ProvideAuthService creates service with injected dependencies
func ProvideAuthService(repo *repository.AuthRepository) *service.AuthService {
	return service.NewAuthService(repo)
}

// ProvideAuthHandler creates handler with injected service
func ProvideAuthHandler(svc *service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(svc)
}

// RegisterAuthRoutes registers routes with DI
func RegisterAuthRoutes(rg *types.RouterGroup, h *handler.AuthHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all auth module dependencies
var Module = fx.Module("auth",
	fx.Provide(
		ProvideAuthRepository,
		ProvideAuthService,
		ProvideAuthHandler,
	),
	fx.Invoke(RegisterAuthRoutes),
)
