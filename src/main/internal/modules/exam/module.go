package exam

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/service"
	studentRepository "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/repository"
	studentService "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// =============================================================================
// Repository Providers
// =============================================================================

// ProvideExamRepository creates exam repository with injected DB
func ProvideExamRepository(db *gorm.DB) *repository.ExamRepository {
	return repository.NewExamRepository(db)
}

// ProvideExamSessionRepository creates exam session repository with injected DB
func ProvideExamSessionRepository(db *gorm.DB) *repository.ExamSessionRepository {
	return repository.NewExamSessionRepository(db)
}

// ProvideTypeExamRepository creates type exam repository with injected DB
func ProvideTypeExamRepository(db *gorm.DB) *repository.TypeExamRepository {
	return repository.NewTypeExamRepository(db)
}

// =============================================================================
// Service Providers
// =============================================================================

// ProvideExamService creates exam service with injected repository
func ProvideExamService(repo *repository.ExamRepository) *service.ExamService {
	return service.NewExamService(repo)
}

// ProvideExamSessionService creates exam session service with injected repositories
func ProvideExamSessionService(repo *repository.ExamSessionRepository, studentRepo *studentRepository.StudentRepository) *service.ExamSessionService {
	return service.NewExamSessionService(repo, studentRepo)
}

// ProvideTypeExamService creates type exam service with injected repository
func ProvideTypeExamService(repo *repository.TypeExamRepository) *service.TypeExamService {
	return service.NewTypeExamService(repo)
}

// =============================================================================
// Handler Providers
// =============================================================================

// ProvideExamHandler creates exam handler with injected service
func ProvideExamHandler(svc *service.ExamService) *handler.ExamHandler {
	return handler.NewExamHandler(svc)
}

// ProvideBankHandler creates bank handler with injected service
func ProvideBankHandler(svc *service.ExamService) *handler.BankHandler {
	return handler.NewBankHandler(svc)
}

// ProvideQuestionHandler creates question handler with injected service
func ProvideQuestionHandler(svc *service.ExamService) *handler.QuestionHandler {
	return handler.NewQuestionHandler(svc)
}

// ProvideSessionHandler creates session handler with injected service
func ProvideSessionHandler(svc *service.ExamSessionService) *handler.SessionHandler {
	return handler.NewSessionHandler(svc)
}

// ProvideTypeExamHandler creates type exam handler with injected service
func ProvideTypeExamHandler(svc *service.TypeExamService) *handler.TypeExamHandler {
	return handler.NewTypeExamHandler(svc)
}

// ProvideCBTHandler creates CBT handler with injected services
func ProvideCBTHandler(examSessionSvc *service.ExamSessionService, studentSvc *studentService.StudentService) *handler.CBTHandler {
	return handler.NewCBTHandler(examSessionSvc, studentSvc)
}

// =============================================================================
// Route Registration
// =============================================================================

// ExamHandlers contains all exam handlers for route registration
type ExamHandlers struct {
	fx.In

	ExamHandler     *handler.ExamHandler
	BankHandler     *handler.BankHandler
	QuestionHandler *handler.QuestionHandler
	SessionHandler  *handler.SessionHandler
	TypeExamHandler *handler.TypeExamHandler
	CBTHandler      *handler.CBTHandler
}

// RegisterExamRoutes registers all exam routes with DI
func RegisterExamRoutes(rg *types.RouterGroup, handlers ExamHandlers) {
	RegisterRoutesWithDI(rg.RouterGroup, handlers)
}

// =============================================================================
// Module Definition
// =============================================================================

// Module provides all exam module dependencies
var Module = fx.Module("exam",
	fx.Provide(
		// Repositories
		ProvideExamRepository,
		ProvideExamSessionRepository,
		ProvideTypeExamRepository,
		// Services
		ProvideExamService,
		ProvideExamSessionService,
		ProvideTypeExamService,
		// Handlers
		ProvideExamHandler,
		ProvideBankHandler,
		ProvideQuestionHandler,
		ProvideSessionHandler,
		ProvideTypeExamHandler,
		ProvideCBTHandler,
	),
	fx.Invoke(RegisterExamRoutes),
)
