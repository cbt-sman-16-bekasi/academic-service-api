package academicyear

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module provides all academic year dependencies
var Module = fx.Options(
	fx.Provide(
		ProvideAcademicYearRepository,
		ProvideAcademicYearService,
		ProvideAcademicYearHandler,
	),
	fx.Invoke(RegisterAcademicYearRoutes),
)

// ProvideAcademicYearRepository creates a new repository
func ProvideAcademicYearRepository(db *gorm.DB) *repository.AcademicYearRepository {
	return repository.NewAcademicYearRepository(db)
}

// ProvideAcademicYearService creates a new service
func ProvideAcademicYearService(repo *repository.AcademicYearRepository) *service.AcademicYearService {
	return service.NewAcademicYearService(repo)
}

// ProvideAcademicYearHandler creates a new handler
func ProvideAcademicYearHandler(svc *service.AcademicYearService) *handler.AcademicYearHandler {
	return handler.NewAcademicYearHandler(svc)
}

// RegisterAcademicYearRoutes registers all routes
func RegisterAcademicYearRoutes(router *gin.Engine, h *handler.AcademicYearHandler) {
	RegisterRoutes(router, h)
}
