package service

import (
	"fmt"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/client"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/entity"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/repository"
	"github.com/rs/zerolog/log"
)

// SyncExecutor handles the actual synchronization with Dapodik
type SyncExecutor struct {
	repo      *repository.DapodikRepository
	client    *client.DapodikClient
	history   *entity.DapodikSyncHistory
	syncStats SyncStats
}

// SyncStats tracks synchronization statistics
type SyncStats struct {
	TotalStudents  int
	TotalTeachers  int
	TotalClasses   int
	SyncedStudents int
	SyncedTeachers int
	SyncedClasses  int
	Errors         []string
}

// SyncExecutorConfig holds configuration for creating a SyncExecutor
type SyncExecutorConfig struct {
	DapodikURL  string
	APIKey      string
	NPSN        string
	UseProxy    bool // Toggle: true = via proxy, false = direct
	ProxyURL    string
	ProxyAPIKey string
}

// NewSyncExecutor creates a new sync executor
// npsn is required for all Dapodik API calls
func NewSyncExecutor(repo *repository.DapodikRepository, cfg SyncExecutorConfig, history *entity.DapodikSyncHistory) *SyncExecutor {
	return &SyncExecutor{
		repo: repo,
		client: client.NewDapodikClientWithConfig(client.DapodikClientConfig{
			BaseURL:     cfg.DapodikURL,
			APIKey:      cfg.APIKey,
			NPSN:        cfg.NPSN,
			UseProxy:    cfg.UseProxy,
			ProxyURL:    cfg.ProxyURL,
			ProxyAPIKey: cfg.ProxyAPIKey,
		}),
		history: history,
	}
}

// Execute runs the synchronization process
func (s *SyncExecutor) Execute() error {
	log.Info().
		Uint("historyID", s.history.ID).
		Str("syncType", string(s.history.SyncType)).
		Msg("🚀 Starting Dapodik sync execution")

	// Update status to in progress
	s.history.Status = entity.SyncStatusInProgress
	_ = s.repo.UpdateHistory(s.history)

	var finalErr error

	// Execute sync steps
	if err := s.syncSekolah(); err != nil {
		s.syncStats.Errors = append(s.syncStats.Errors, "Sekolah: "+err.Error())
		log.Error().Err(err).Msg("Failed to sync sekolah")
	}

	if err := s.syncPTK(); err != nil {
		s.syncStats.Errors = append(s.syncStats.Errors, "PTK: "+err.Error())
		log.Error().Err(err).Msg("Failed to sync PTK")
	}

	if err := s.syncRombonganBelajar(); err != nil {
		s.syncStats.Errors = append(s.syncStats.Errors, "Rombel: "+err.Error())
		log.Error().Err(err).Msg("Failed to sync rombongan belajar")
	}

	if err := s.syncPesertaDidik(); err != nil {
		s.syncStats.Errors = append(s.syncStats.Errors, "PD: "+err.Error())
		log.Error().Err(err).Msg("Failed to sync peserta didik")
		finalErr = err
	}

	// Update history with results
	now := time.Now()
	s.history.CompletedAt = &now
	s.history.TotalStudents = s.syncStats.TotalStudents
	s.history.TotalTeachers = s.syncStats.TotalTeachers
	s.history.TotalClasses = s.syncStats.TotalClasses
	s.history.SyncedStudents = s.syncStats.SyncedStudents
	s.history.SyncedTeachers = s.syncStats.SyncedTeachers
	s.history.SyncedClasses = s.syncStats.SyncedClasses

	if len(s.syncStats.Errors) > 0 {
		s.history.Status = entity.SyncStatusFailed
		s.history.ErrorMessage = fmt.Sprintf("Errors: %v", s.syncStats.Errors)
	} else {
		s.history.Status = entity.SyncStatusSuccess
	}

	_ = s.repo.UpdateHistory(s.history)

	log.Info().
		Uint("historyID", s.history.ID).
		Str("status", string(s.history.Status)).
		Int("totalStudents", s.syncStats.TotalStudents).
		Int("syncedStudents", s.syncStats.SyncedStudents).
		Int("totalTeachers", s.syncStats.TotalTeachers).
		Int("syncedTeachers", s.syncStats.SyncedTeachers).
		Int("totalClasses", s.syncStats.TotalClasses).
		Int("syncedClasses", s.syncStats.SyncedClasses).
		Msg("✅ Dapodik sync completed")

	return finalErr
}

// syncSekolah syncs school data
func (s *SyncExecutor) syncSekolah() error {
	log.Info().Msg("📊 Syncing sekolah data...")

	sekolah, err := s.client.GetSekolah()
	if err != nil {
		return fmt.Errorf("failed to fetch sekolah: %w", err)
	}

	log.Info().
		Str("nama", sekolah.Nama).
		Str("npsn", sekolah.NPSN).
		Msg("Fetched sekolah data")

	// TODO: Update school data in local database
	// For now, just log the data

	return nil
}

// syncPTK syncs teacher data
func (s *SyncExecutor) syncPTK() error {
	log.Info().Msg("👨‍🏫 Syncing PTK (teachers) data...")

	teachers, total, err := s.client.GetPTK()
	if err != nil {
		return fmt.Errorf("failed to fetch PTK: %w", err)
	}

	s.syncStats.TotalTeachers = total
	log.Info().Int("total", total).Int("fetched", len(teachers)).Msg("Fetched PTK data")

	// Process teachers
	for _, ptk := range teachers {
		// TODO: Implement actual teacher sync logic
		// - Check if teacher exists by NUPTK or NIK
		// - Create or update teacher record
		// - Create or update user account
		log.Debug().
			Str("nama", ptk.Nama).
			Str("nuptk", ptk.NUPTK).
			Str("jabatan", ptk.JabatanPTKStr).
			Msg("Processing PTK")

		s.syncStats.SyncedTeachers++
	}

	return nil
}

// syncRombonganBelajar syncs class data
func (s *SyncExecutor) syncRombonganBelajar() error {
	log.Info().Msg("🏫 Syncing rombongan belajar (classes) data...")

	classes, total, err := s.client.GetRombonganBelajar()
	if err != nil {
		return fmt.Errorf("failed to fetch rombongan belajar: %w", err)
	}

	s.syncStats.TotalClasses = total
	log.Info().Int("total", total).Int("fetched", len(classes)).Msg("Fetched rombongan belajar data")

	// Process classes
	for _, rombel := range classes {
		// TODO: Implement actual class sync logic
		// - Check if class exists
		// - Create or update class record
		// - Sync class members
		log.Debug().
			Str("nama", rombel.Nama).
			Str("tingkat", rombel.TingkatPendidikanStr).
			Str("waliKelas", rombel.PTKIDStr).
			Int("jumlahSiswa", len(rombel.AnggotaRombel)).
			Msg("Processing rombel")

		s.syncStats.SyncedClasses++
	}

	return nil
}

// syncPesertaDidik syncs student data
func (s *SyncExecutor) syncPesertaDidik() error {
	log.Info().Msg("👨‍🎓 Syncing peserta didik (students) data...")

	students, total, err := s.client.GetPesertaDidik()
	if err != nil {
		return fmt.Errorf("failed to fetch peserta didik: %w", err)
	}

	s.syncStats.TotalStudents = total
	log.Info().Int("total", total).Int("fetched", len(students)).Msg("Fetched peserta didik data")

	// Process students
	for _, pd := range students {
		// TODO: Implement actual student sync logic
		// - Check if student exists by NISN
		// - Create or update student record
		// - Create or update user account
		// - Update class membership
		log.Debug().
			Str("nama", pd.Nama).
			Str("nisn", pd.NISN).
			Str("kelas", pd.NamaRombel).
			Msg("Processing peserta didik")

		s.syncStats.SyncedStudents++
	}

	return nil
}

// =============================================================================
// Preview Functions (for testing without actual sync)
// =============================================================================

// PreviewPengguna fetches users for preview
func (s *SyncExecutor) PreviewPengguna() ([]client.Pengguna, int, error) {
	return s.client.GetPengguna()
}

// PreviewSekolah fetches school for preview
func (s *SyncExecutor) PreviewSekolah() (*client.Sekolah, error) {
	return s.client.GetSekolah()
}

// PreviewRombonganBelajar fetches classes for preview
func (s *SyncExecutor) PreviewRombonganBelajar() ([]client.RombonganBelajar, int, error) {
	return s.client.GetRombonganBelajar()
}

// PreviewPTK fetches teachers for preview
func (s *SyncExecutor) PreviewPTK() ([]client.PTK, int, error) {
	return s.client.GetPTK()
}

// PreviewPesertaDidik fetches students for preview
func (s *SyncExecutor) PreviewPesertaDidik() ([]client.PesertaDidik, int, error) {
	return s.client.GetPesertaDidik()
}
