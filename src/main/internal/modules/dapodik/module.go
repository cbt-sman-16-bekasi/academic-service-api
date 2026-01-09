package dapodik

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/service"
	schoolRepo "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideDapodikRepository creates repository with injected DB
func ProvideDapodikRepository(db *gorm.DB) *repository.DapodikRepository {
	return repository.NewDapodikRepository(db)
}

// ProvideDapodikService creates service with injected dependencies
// Requires school repository to get NPSN for Dapodik API calls
func ProvideDapodikService(repo *repository.DapodikRepository, schoolRepository *schoolRepo.SchoolRepository) *service.DapodikService {
	return service.NewDapodikService(repo, schoolRepository)
}

// ProvideDapodikHandler creates handler with injected service
func ProvideDapodikHandler(svc *service.DapodikService) *handler.DapodikHandler {
	return handler.NewDapodikHandler(svc)
}

// RegisterDapodikRoutes registers routes with DI
func RegisterDapodikRoutes(rg *types.RouterGroup, h *handler.DapodikHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all dapodik module dependencies
var Module = fx.Module("dapodik",
	fx.Provide(
		ProvideDapodikRepository,
		ProvideDapodikService,
		ProvideDapodikHandler,
	),
	fx.Invoke(RegisterDapodikRoutes),
)
