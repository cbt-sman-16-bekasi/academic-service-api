package student

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/service"
	userService "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideStudentRepository creates repository with injected DB
func ProvideStudentRepository(db *gorm.DB) *repository.StudentRepository {
	return repository.NewStudentRepository(db)
}

// ProvideStudentService creates service with injected dependencies
func ProvideStudentService(repo *repository.StudentRepository, userSvc *userService.UserService) *service.StudentService {
	return service.NewStudentService(repo, userSvc)
}

// ProvideStudentHandler creates handler with injected service
func ProvideStudentHandler(svc *service.StudentService) *handler.StudentHandler {
	return handler.NewStudentHandler(svc)
}

// RegisterStudentRoutes registers student routes with DI
func RegisterStudentRoutes(rg *types.RouterGroup, h *handler.StudentHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all student module dependencies
var Module = fx.Module("student",
	fx.Provide(
		ProvideStudentRepository,
		ProvideStudentService,
		ProvideStudentHandler,
	),
	fx.Invoke(RegisterStudentRoutes),
)
