package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/entity"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/repository"
	schoolRepo "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type DapodikService struct {
	repo       *repository.DapodikRepository
	schoolRepo *schoolRepo.SchoolRepository
}

// NewDapodikService creates service with injected repositories
func NewDapodikService(repo *repository.DapodikRepository, schoolRepository *schoolRepo.SchoolRepository) *DapodikService {
	return &DapodikService{
		repo:       repo,
		schoolRepo: schoolRepository,
	}
}

// GetConfig mengambil konfigurasi DAPODIK berdasarkan school code
func (s *DapodikService) GetConfig(c *gin.Context) *dto.DapodikConfigResponse {
	claims := jwt.GetDataClaims(c)

	config, err := s.repo.FindConfigBySchoolCode(claims.SchoolCode)
	if err != nil {
		return nil
	}

	return &dto.DapodikConfigResponse{
		ID:          config.ID,
		DapodikURL:  config.DapodikURL,
		APIKey:      maskAPIKey(config.APIKey),
		UseProxy:    config.UseProxy,
		ProxyURL:    config.ProxyURL,
		ProxyAPIKey: maskAPIKey(config.ProxyAPIKey),
		IsActive:    config.IsActive,
		LastSyncAt:  config.LastSyncAt,
		CreatedAt:   config.CreatedAt,
		UpdatedAt:   config.UpdatedAt,
	}
}

// SaveConfig menyimpan atau update konfigurasi DAPODIK
func (s *DapodikService) SaveConfig(c *gin.Context, req dto.SaveDapodikConfigRequest) *dto.DapodikConfigResponse {
	claims := jwt.GetDataClaims(c)

	config, err := s.repo.FindConfigBySchoolCode(claims.SchoolCode)

	if err != nil {
		// Create new config
		config = &entity.DapodikConfig{
			SchoolCode:  claims.SchoolCode,
			DapodikURL:  normalizeURL(req.DapodikURL),
			APIKey:      req.APIKey,
			UseProxy:    req.UseProxy,
			ProxyURL:    req.ProxyURL,
			ProxyAPIKey: req.ProxyAPIKey,
			IsActive:    true,
		}
		s.repo.CreateConfig(config)
	} else {
		// Update existing config
		config.DapodikURL = normalizeURL(req.DapodikURL)
		config.APIKey = req.APIKey
		config.UseProxy = req.UseProxy
		config.ProxyURL = req.ProxyURL
		config.ProxyAPIKey = req.ProxyAPIKey
		s.repo.UpdateConfig(config)
	}

	return &dto.DapodikConfigResponse{
		ID:          config.ID,
		DapodikURL:  config.DapodikURL,
		APIKey:      maskAPIKey(config.APIKey),
		UseProxy:    config.UseProxy,
		ProxyURL:    config.ProxyURL,
		ProxyAPIKey: maskAPIKey(config.ProxyAPIKey),
		IsActive:    config.IsActive,
		LastSyncAt:  config.LastSyncAt,
		CreatedAt:   config.CreatedAt,
		UpdatedAt:   config.UpdatedAt,
	}
}

// TestConnection melakukan test koneksi ke API DAPODIK (langsung atau via proxy)
func (s *DapodikService) TestConnection(req dto.TestConnectionRequest) *dto.TestConnectionResponse {
	startTime := time.Now()

	// Determine target URL based on UseProxy flag
	var targetURL string
	var viaProxy bool
	dapodikURL := normalizeURL(req.DapodikURL)

	if req.UseProxy && req.ProxyURL != "" {
		// Use proxy: send request to proxy with target URL as query parameter
		targetURL = fmt.Sprintf("%s/proxy?target=%s", req.ProxyURL, dapodikURL)
		viaProxy = true
	} else {
		// Direct connection to Dapodik
		targetURL = dapodikURL
		viaProxy = false
	}

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	httpReq, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return &dto.TestConnectionResponse{
			Success: false,
			Message: "Gagal membuat request: " + err.Error(),
			Latency: "0ms",
		}
	}

	// Set Dapodik API Key header
	httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)

	// Set proxy API key if using proxy
	if viaProxy && req.ProxyAPIKey != "" {
		httpReq.Header.Set("X-API-Key", req.ProxyAPIKey)
	}

	resp, err := httpClient.Do(httpReq)
	latency := time.Since(startTime)

	if err != nil {
		return &dto.TestConnectionResponse{
			Success: false,
			Message: "Koneksi gagal: " + err.Error(),
			Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return &dto.TestConnectionResponse{
			Success: false,
			Message: "API Key tidak valid atau tidak memiliki akses",
			Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
		}
	}

	if resp.StatusCode >= 400 {
		return &dto.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Server DAPODIK mengembalikan error: %d", resp.StatusCode),
			Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
		}
	}

	message := "Koneksi langsung berhasil! Server DAPODIK dapat dijangkau."
	if viaProxy {
		message = "Koneksi via proxy berhasil! Server DAPODIK dapat dijangkau."
	}

	return &dto.TestConnectionResponse{
		Success: true,
		Message: message,
		Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
	}
}

// TriggerSync memulai proses sinkronisasi dengan DAPODIK
func (s *DapodikService) TriggerSync(c *gin.Context, req dto.TriggerSyncRequest) *dto.DapodikSyncHistoryResponse {
	claims := jwt.GetDataClaims(c)

	// Validasi sync type
	syncType := entity.SyncType(req.SyncType)
	if syncType != entity.SyncTypeUpdate && syncType != entity.SyncTypeFullReset {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Tipe sinkronisasi tidak valid. Gunakan UPDATE atau FULL_RESET"))
	}

	// Cek apakah ada config
	_, err := s.repo.FindConfigBySchoolCode(claims.SchoolCode)
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Konfigurasi DAPODIK belum diatur"))
	}

	// Cek apakah ada sinkronisasi yang sedang berjalan
	runningSync, _ := s.repo.FindRunningSync(claims.SchoolCode)
	if runningSync != nil && runningSync.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Masih ada sinkronisasi yang sedang berjalan"))
	}

	// Buat history baru
	now := time.Now()
	history := entity.DapodikSyncHistory{
		SchoolCode:        claims.SchoolCode,
		SyncType:          syncType,
		Status:            entity.SyncStatusPending,
		StartedAt:         &now,
		TriggeredByUserId: claims.Id,
		TriggeredByName:   claims.Username,
	}
	_ = s.repo.CreateHistory(&history)

	// Get config for API credentials
	config, _ := s.repo.FindConfigBySchoolCode(claims.SchoolCode)

	// Get school NPSN
	school, err := s.schoolRepo.FindSchoolByCode(claims.SchoolCode)
	if err != nil || school.NPSN == "" {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "NPSN sekolah tidak ditemukan"))
	}

	// Start async sync process
	go func() {
		executor := NewSyncExecutor(s.repo, SyncExecutorConfig{
			DapodikURL:  config.DapodikURL,
			APIKey:      config.APIKey,
			NPSN:        school.NPSN,
			UseProxy:    config.UseProxy,
			ProxyURL:    config.ProxyURL,
			ProxyAPIKey: config.ProxyAPIKey,
		}, &history)
		if err := executor.Execute(); err != nil {
			// Error already logged and saved in history
		}

		// Update last sync time in config
		now := time.Now()
		config.LastSyncAt = &now
		_ = s.repo.UpdateConfig(config)
	}()

	return s.mapHistoryToResponse(&history)
}

// GetSyncHistory mengambil riwayat sinkronisasi dengan pagination
func (s *DapodikService) GetSyncHistory(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)

	// Set filter untuk school_code
	filter := map[string]interface{}{}
	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	page := database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetRequest(&request).
		SetModel([]entity.DapodikSyncHistory{}).
		FindAllPaging()

	return page
}

// GetSyncHistoryDetail mengambil detail satu riwayat sinkronisasi
func (s *DapodikService) GetSyncHistoryDetail(c *gin.Context, id uint) *dto.DapodikSyncHistoryResponse {
	claims := jwt.GetDataClaims(c)

	history, err := s.repo.FindHistoryByIDAndSchool(id, claims.SchoolCode)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Riwayat sinkronisasi tidak ditemukan"))
	}

	return s.mapHistoryToResponse(history)
}

// GetSyncSummary mengambil ringkasan sinkronisasi
func (s *DapodikService) GetSyncSummary(c *gin.Context) *dto.SyncSummaryResponse {
	claims := jwt.GetDataClaims(c)

	totalSync := s.repo.CountBySchool(claims.SchoolCode)
	totalSuccess := s.repo.CountBySchoolAndStatus(claims.SchoolCode, entity.SyncStatusSuccess)
	totalFailed := s.repo.CountBySchoolAndStatus(claims.SchoolCode, entity.SyncStatusFailed)

	lastSync, _ := s.repo.FindLastSyncBySchool(claims.SchoolCode)

	var lastSyncAt *time.Time
	var lastSyncStatus string
	if lastSync != nil && lastSync.ID != 0 {
		lastSyncAt = lastSync.StartedAt
		lastSyncStatus = string(lastSync.Status)
	}

	return &dto.SyncSummaryResponse{
		TotalSync:      int(totalSync),
		LastSyncAt:     lastSyncAt,
		LastSyncStatus: lastSyncStatus,
		TotalSuccess:   int(totalSuccess),
		TotalFailed:    int(totalFailed),
	}
}

// Helper functions

func (s *DapodikService) mapHistoryToResponse(history *entity.DapodikSyncHistory) *dto.DapodikSyncHistoryResponse {
	duration := ""
	if history.StartedAt != nil && history.CompletedAt != nil {
		dur := history.CompletedAt.Sub(*history.StartedAt)
		duration = formatDuration(dur)
	}

	return &dto.DapodikSyncHistoryResponse{
		ID:              history.ID,
		SyncType:        string(history.SyncType),
		SyncTypeLabel:   getSyncTypeLabel(history.SyncType),
		Status:          string(history.Status),
		StatusLabel:     getStatusLabel(history.Status),
		StartedAt:       history.StartedAt,
		CompletedAt:     history.CompletedAt,
		Duration:        duration,
		TotalStudents:   history.TotalStudents,
		TotalClasses:    history.TotalClasses,
		TotalTeachers:   history.TotalTeachers,
		SyncedStudents:  history.SyncedStudents,
		SyncedClasses:   history.SyncedClasses,
		SyncedTeachers:  history.SyncedTeachers,
		ErrorMessage:    history.ErrorMessage,
		TriggeredByName: history.TriggeredByName,
		CreatedAt:       history.CreatedAt,
	}
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 4 {
		return strings.Repeat("•", len(apiKey))
	}
	return strings.Repeat("•", len(apiKey)-4) + apiKey[len(apiKey)-4:]
}

func normalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return strings.TrimSuffix(url, "/")
}

func getSyncTypeLabel(syncType entity.SyncType) string {
	switch syncType {
	case entity.SyncTypeUpdate:
		return "Update & Tambah Data Baru"
	case entity.SyncTypeFullReset:
		return "Reset & Import Ulang Semua Data"
	default:
		return string(syncType)
	}
}

func getStatusLabel(status entity.SyncStatus) string {
	switch status {
	case entity.SyncStatusPending:
		return "Menunggu"
	case entity.SyncStatusInProgress:
		return "Sedang Berjalan"
	case entity.SyncStatusSuccess:
		return "Berhasil"
	case entity.SyncStatusFailed:
		return "Gagal"
	default:
		return string(status)
	}
}

func formatDuration(d time.Duration) string {
	if d.Hours() >= 1 {
		return fmt.Sprintf("%.0f jam %.0f menit", d.Hours(), d.Minutes()-d.Hours()*60)
	}
	if d.Minutes() >= 1 {
		return fmt.Sprintf("%.0f menit %.0f detik", d.Minutes(), d.Seconds()-d.Minutes()*60)
	}
	return fmt.Sprintf("%.0f detik", d.Seconds())
}

// =============================================================================
// Preview Endpoints - Fetch data from Dapodik without syncing
// =============================================================================

// createExecutorForPreview creates a sync executor for preview operations
func (s *DapodikService) createExecutorForPreview(c *gin.Context) *SyncExecutor {
	claims := jwt.GetDataClaims(c)
	config, err := s.repo.FindConfigBySchoolCode(claims.SchoolCode)
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Konfigurasi DAPODIK belum diatur"))
	}

	// Get school NPSN - required for Dapodik API
	school, err := s.schoolRepo.FindSchoolByCode(claims.SchoolCode)
	if err != nil || school.NPSN == "" {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "NPSN sekolah tidak ditemukan. Pastikan data sekolah sudah diisi."))
	}

	return NewSyncExecutor(s.repo, SyncExecutorConfig{
		DapodikURL:  config.DapodikURL,
		APIKey:      config.APIKey,
		NPSN:        school.NPSN,
		UseProxy:    config.UseProxy,
		ProxyURL:    config.ProxyURL,
		ProxyAPIKey: config.ProxyAPIKey,
	}, nil)
}

// PreviewSekolah fetches school data from Dapodik
func (s *DapodikService) PreviewSekolah(c *gin.Context) *dto.PreviewSekolahResponse {
	executor := s.createExecutorForPreview(c)
	sekolah, err := executor.PreviewSekolah()
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Gagal mengambil data sekolah: "+err.Error()))
	}

	return &dto.PreviewSekolahResponse{
		SekolahID:        sekolah.SekolahID,
		Nama:             sekolah.Nama,
		NPSN:             sekolah.NPSN,
		NSS:              sekolah.NSS,
		BentukPendidikan: sekolah.BentukPendidikanStr,
		StatusSekolah:    sekolah.StatusSekolahStr,
		Alamat:           sekolah.AlamatJalan,
		Kelurahan:        sekolah.DesaKelurahan,
		Kecamatan:        sekolah.Kecamatan,
		KabupatenKota:    sekolah.KabupatenKota,
		Provinsi:         sekolah.Provinsi,
		Email:            sekolah.Email,
		Website:          sekolah.Website,
		Telepon:          sekolah.NomorTelepon,
	}
}

// PreviewPTK fetches teachers from Dapodik
func (s *DapodikService) PreviewPTK(c *gin.Context) *dto.PreviewResponse[[]dto.PreviewPTKResponse] {
	executor := s.createExecutorForPreview(c)
	teachers, total, err := executor.PreviewPTK()
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Gagal mengambil data PTK: "+err.Error()))
	}

	var result []dto.PreviewPTKResponse
	for _, ptk := range teachers {
		result = append(result, dto.PreviewPTKResponse{
			PTKID:              ptk.PTKID,
			Nama:               ptk.Nama,
			NUPTK:              ptk.NUPTK,
			NIP:                ptk.NIP,
			NIK:                ptk.NIK,
			JenisKelamin:       ptk.JenisKelamin,
			TempatLahir:        ptk.TempatLahir,
			TanggalLahir:       ptk.TanggalLahir,
			JenisPTK:           ptk.JenisPTKStr,
			JabatanPTK:         ptk.JabatanPTKStr,
			StatusKepegawaian:  ptk.StatusKepegawaianStr,
			PendidikanTerakhir: ptk.PendidikanTerakhir,
		})
	}

	return &dto.PreviewResponse[[]dto.PreviewPTKResponse]{
		Total: total,
		Data:  result,
	}
}

// PreviewRombonganBelajar fetches classes from Dapodik
func (s *DapodikService) PreviewRombonganBelajar(c *gin.Context) *dto.PreviewResponse[[]dto.PreviewRombelResponse] {
	executor := s.createExecutorForPreview(c)
	classes, total, err := executor.PreviewRombonganBelajar()
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Gagal mengambil data rombongan belajar: "+err.Error()))
	}

	var result []dto.PreviewRombelResponse
	for _, rombel := range classes {
		result = append(result, dto.PreviewRombelResponse{
			RombelID:          rombel.RombonganBelajarID,
			Nama:              rombel.Nama,
			TingkatPendidikan: rombel.TingkatPendidikanStr,
			Jurusan:           rombel.JurusanStr,
			WaliKelas:         rombel.PTKIDStr,
			Kurikulum:         rombel.KurikulumStr,
			JumlahSiswa:       len(rombel.AnggotaRombel),
			Semester:          rombel.SemesterID,
		})
	}

	return &dto.PreviewResponse[[]dto.PreviewRombelResponse]{
		Total: total,
		Data:  result,
	}
}

// PreviewPesertaDidik fetches students from Dapodik
func (s *DapodikService) PreviewPesertaDidik(c *gin.Context) *dto.PreviewResponse[[]dto.PreviewPesertaDidikResponse] {
	executor := s.createExecutorForPreview(c)
	students, total, err := executor.PreviewPesertaDidik()
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Gagal mengambil data peserta didik: "+err.Error()))
	}

	var result []dto.PreviewPesertaDidikResponse
	for _, pd := range students {
		namaAyah := ""
		if pd.NamaAyah != nil {
			namaAyah = *pd.NamaAyah
		}
		namaIbu := ""
		if pd.NamaIbu != nil {
			namaIbu = *pd.NamaIbu
		}
		email := ""
		if pd.Email != nil {
			email = *pd.Email
		}
		noHP := ""
		if pd.NomorTeleponSeluler != nil {
			noHP = *pd.NomorTeleponSeluler
		}

		result = append(result, dto.PreviewPesertaDidikResponse{
			PesertaDidikID: pd.PesertaDidikID,
			NISN:           pd.NISN,
			NIPD:           pd.NIPD,
			Nama:           pd.Nama,
			JenisKelamin:   pd.JenisKelamin,
			TempatLahir:    pd.TempatLahir,
			TanggalLahir:   pd.TanggalLahir,
			NamaKelas:      pd.NamaRombel,
			TingkatKelas:   pd.TingkatPendidikanID,
			NamaAyah:       namaAyah,
			NamaIbu:        namaIbu,
			Email:          email,
			NoHP:           noHP,
		})
	}

	return &dto.PreviewResponse[[]dto.PreviewPesertaDidikResponse]{
		Total: total,
		Data:  result,
	}
}

// PreviewPengguna fetches users from Dapodik
func (s *DapodikService) PreviewPengguna(c *gin.Context) *dto.PreviewResponse[[]dto.PreviewPenggunaResponse] {
	executor := s.createExecutorForPreview(c)
	users, total, err := executor.PreviewPengguna()
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Gagal mengambil data pengguna: "+err.Error()))
	}

	var result []dto.PreviewPenggunaResponse
	for _, user := range users {
		ptkID := ""
		if user.PTKID != nil {
			ptkID = *user.PTKID
		}

		result = append(result, dto.PreviewPenggunaResponse{
			PenggunaID: user.PenggunaID,
			Username:   user.Username,
			Nama:       user.Nama,
			Peran:      user.PeranIDStr,
			NoHP:       user.NoHP,
			PTKID:      ptkID,
		})
	}

	return &dto.PreviewResponse[[]dto.PreviewPenggunaResponse]{
		Total: total,
		Data:  result,
	}
}
