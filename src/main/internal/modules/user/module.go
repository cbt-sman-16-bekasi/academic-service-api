package user

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideUserRepository creates repository with injected DB
func ProvideUserRepository(db *gorm.DB) *repository.UserRepository {
	return repository.NewUserRepository(db)
}

// ProvideUserService creates service with injected dependencies
func ProvideUserService(repo *repository.UserRepository) *service.UserService {
	return service.NewUserService(repo)
}

// ProvideUserHandler creates handler with injected service
func ProvideUserHandler(svc *service.UserService) *handler.UserHandler {
	return handler.NewUserHandler(svc)
}

// RegisterUserRoutes registers user routes with DI
func RegisterUserRoutes(rg *types.RouterGroup, h *handler.UserHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all user module dependencies
var Module = fx.Module("user",
	fx.Provide(
		ProvideUserRepository,
		ProvideUserService,
		ProvideUserHandler,
	),
	fx.Invoke(RegisterUserRoutes),
)
