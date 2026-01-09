package class

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideClassRepository creates repository with injected DB
func ProvideClassRepository(db *gorm.DB) *repository.ClassRepository {
	return repository.NewClassRepository(db)
}

// ProvideClassService creates service with injected dependencies
func ProvideClassService(repo *repository.ClassRepository) *service.ClassService {
	return service.NewClassService(repo)
}

// ProvideClassHandler creates handler with injected service
func ProvideClassHandler(svc *service.ClassService) *handler.ClassHandler {
	return handler.NewClassHandler(svc)
}

// RegisterClassRoutes registers routes with DI
func RegisterClassRoutes(rg *types.RouterGroup, h *handler.ClassHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all class module dependencies
var Module = fx.Module("class",
	fx.Provide(
		ProvideClassRepository,
		ProvideClassService,
		ProvideClassHandler,
	),
	fx.Invoke(RegisterClassRoutes),
)
