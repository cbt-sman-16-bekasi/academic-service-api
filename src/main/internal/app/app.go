package app

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/storage"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user"
	"go.uber.org/fx"
)

// NewApp creates the main FX application with all modules
func NewApp() *fx.App {
	// Initialize pre-fx dependencies (env, logger, timezone)
	initializeApp()

	return fx.New(
		// Core infrastructure
		CoreModule,

		// Feature modules
		auth.Module,
		class.Module,
		curriculum.Module,
		dapodik.Module,
		exam.Module,
		school.Module,
		spp.Module,
		storage.Module,
		student.Module,
		teacher.Module,
		user.Module,
	)
}

// Run starts the FX application
func Run() {
	app := NewApp()
	app.Run()
}
