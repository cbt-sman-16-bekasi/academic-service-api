package curriculum

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideSubjectRepository creates repository with injected DB
func ProvideSubjectRepository(db *gorm.DB) *repository.SubjectRepository {
	return repository.NewSubjectRepository(db)
}

// ProvideSubjectService creates service with injected dependencies
func ProvideSubjectService(repo *repository.SubjectRepository) *service.SubjectService {
	return service.NewSubjectService(repo)
}

// ProvideSubjectHandler creates handler with injected service
func ProvideSubjectHandler(svc *service.SubjectService) *handler.SubjectHandler {
	return handler.NewSubjectHandler(svc)
}

// RegisterCurriculumRoutes registers routes with DI
func RegisterCurriculumRoutes(rg *types.RouterGroup, h *handler.SubjectHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all curriculum module dependencies
var Module = fx.Module("curriculum",
	fx.Provide(
		ProvideSubjectRepository,
		ProvideSubjectService,
		ProvideSubjectHandler,
	),
	fx.Invoke(RegisterCurriculumRoutes),
)
