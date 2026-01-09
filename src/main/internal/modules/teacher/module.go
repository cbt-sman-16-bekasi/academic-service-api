package teacher

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/service"
	userService "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideTeacherRepository creates repository with injected DB
func ProvideTeacherRepository(db *gorm.DB) *repository.TeacherRepository {
	return repository.NewTeacherRepository(db)
}

// ProvideTeacherService creates service with injected dependencies
func ProvideTeacherService(repo *repository.TeacherRepository, userSvc *userService.UserService) *service.TeacherService {
	return service.NewTeacherService(repo, userSvc)
}

// ProvideTeacherHandler creates handler with injected service
func ProvideTeacherHandler(svc *service.TeacherService) *handler.TeacherHandler {
	return handler.NewTeacherHandler(svc)
}

// RegisterTeacherRoutes registers routes with DI
func RegisterTeacherRoutes(rg *types.RouterGroup, h *handler.TeacherHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all teacher module dependencies
var Module = fx.Module("teacher",
	fx.Provide(
		ProvideTeacherRepository,
		ProvideTeacherService,
		ProvideTeacherHandler,
	),
	fx.Invoke(RegisterTeacherRoutes),
)
