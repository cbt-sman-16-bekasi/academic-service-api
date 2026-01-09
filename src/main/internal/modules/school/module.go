package school

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideSchoolRepository creates repository with injected DB
func ProvideSchoolRepository(db *gorm.DB) *repository.SchoolRepository {
	return repository.NewSchoolRepository(db)
}

// ProvideSchoolService creates service with injected dependencies
func ProvideSchoolService(repo *repository.SchoolRepository) *service.SchoolService {
	return service.NewSchoolService(repo)
}

// ProvideSchoolHandler creates handler with injected service
func ProvideSchoolHandler(svc *service.SchoolService) *handler.SchoolHandler {
	return handler.NewSchoolHandler(svc)
}

// RegisterSchoolRoutes registers routes with DI
func RegisterSchoolRoutes(rg *types.RouterGroup, h *handler.SchoolHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all school module dependencies
var Module = fx.Module("school",
	fx.Provide(
		ProvideSchoolRepository,
		ProvideSchoolService,
		ProvideSchoolHandler,
	),
	fx.Invoke(RegisterSchoolRoutes),
)
