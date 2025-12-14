package dapodik_service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/dapodik_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/dapodik_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
)

type DapodikService struct {
	db *gorm.DB
}

func NewDapodikService() *DapodikService {
	return &DapodikService{
		db: database.GetDB(),
	}
}

// GetConfig mengambil konfigurasi DAPODIK berdasarkan school code
func (s *DapodikService) GetConfig(c *gin.Context) *dapodik_response.DapodikConfigResponse {
	claims := jwt.GetDataClaims(c)

	var config school.DapodikConfig
	result := s.db.Where("school_code = ?", claims.SchoolCode).First(&config)

	if result.Error != nil {
		return nil
	}

	// Mask API key untuk keamanan (tampilkan hanya 4 karakter terakhir)
	maskedAPIKey := maskAPIKey(config.APIKey)

	return &dapodik_response.DapodikConfigResponse{
		ID:         config.ID,
		DapodikURL: config.DapodikURL,
		APIKey:     maskedAPIKey,
		IsActive:   config.IsActive,
		LastSyncAt: config.LastSyncAt,
		CreatedAt:  config.CreatedAt,
		UpdatedAt:  config.UpdatedAt,
	}
}

// SaveConfig menyimpan atau update konfigurasi DAPODIK
func (s *DapodikService) SaveConfig(c *gin.Context, req dapodik_request.SaveDapodikConfigRequest) *dapodik_response.DapodikConfigResponse {
	claims := jwt.GetDataClaims(c)

	var config school.DapodikConfig
	result := s.db.Where("school_code = ?", claims.SchoolCode).First(&config)

	if result.Error != nil {
		// Create new config
		config = school.DapodikConfig{
			SchoolCode: claims.SchoolCode,
			DapodikURL: normalizeURL(req.DapodikURL),
			APIKey:     req.APIKey,
			IsActive:   true,
		}
		s.db.Create(&config)
	} else {
		// Update existing config
		config.DapodikURL = normalizeURL(req.DapodikURL)
		config.APIKey = req.APIKey
		s.db.Save(&config)
	}

	return &dapodik_response.DapodikConfigResponse{
		ID:         config.ID,
		DapodikURL: config.DapodikURL,
		APIKey:     maskAPIKey(config.APIKey),
		IsActive:   config.IsActive,
		LastSyncAt: config.LastSyncAt,
		CreatedAt:  config.CreatedAt,
		UpdatedAt:  config.UpdatedAt,
	}
}

// TestConnection melakukan test koneksi ke API DAPODIK
func (s *DapodikService) TestConnection(req dapodik_request.TestConnectionRequest) *dapodik_response.TestConnectionResponse {
	startTime := time.Now()

	// Buat HTTP request ke endpoint DAPODIK
	url := normalizeURL(req.DapodikURL)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &dapodik_response.TestConnectionResponse{
			Success: false,
			Message: "Gagal membuat request: " + err.Error(),
			Latency: "0ms",
		}
	}

	// Set API Key header
	httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)

	resp, err := client.Do(httpReq)
	latency := time.Since(startTime)

	if err != nil {
		return &dapodik_response.TestConnectionResponse{
			Success: false,
			Message: "Koneksi gagal: " + err.Error(),
			Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return &dapodik_response.TestConnectionResponse{
			Success: false,
			Message: "API Key tidak valid atau tidak memiliki akses",
			Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
		}
	}

	if resp.StatusCode >= 400 {
		return &dapodik_response.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Server DAPODIK mengembalikan error: %d", resp.StatusCode),
			Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
		}
	}

	return &dapodik_response.TestConnectionResponse{
		Success: true,
		Message: "Koneksi berhasil! Server DAPODIK dapat dijangkau.",
		Latency: fmt.Sprintf("%dms", latency.Milliseconds()),
	}
}

// TriggerSync memulai proses sinkronisasi dengan DAPODIK
func (s *DapodikService) TriggerSync(c *gin.Context, req dapodik_request.TriggerSyncRequest) *dapodik_response.DapodikSyncHistoryResponse {
	claims := jwt.GetDataClaims(c)

	// Validasi sync type
	syncType := school.SyncType(req.SyncType)
	if syncType != school.SyncTypeUpdate && syncType != school.SyncTypeFullReset {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Tipe sinkronisasi tidak valid. Gunakan UPDATE atau FULL_RESET"))
	}

	// Cek apakah ada config
	var config school.DapodikConfig
	result := s.db.Where("school_code = ?", claims.SchoolCode).First(&config)
	if result.Error != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Konfigurasi DAPODIK belum diatur"))
	}

	// Cek apakah ada sinkronisasi yang sedang berjalan
	var runningSync school.DapodikSyncHistory
	s.db.Where("school_code = ? AND status IN ?", claims.SchoolCode, []string{string(school.SyncStatusPending), string(school.SyncStatusInProgress)}).First(&runningSync)
	if runningSync.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Masih ada sinkronisasi yang sedang berjalan"))
	}

	// Buat history baru
	now := time.Now()
	history := school.DapodikSyncHistory{
		SchoolCode:        claims.SchoolCode,
		SyncType:          syncType,
		Status:            school.SyncStatusPending,
		StartedAt:         &now,
		TriggeredByUserId: claims.Id,
		TriggeredByName:   claims.Username,
	}
	s.db.Create(&history)

	// TODO: Trigger async sync process here
	// Untuk saat ini, kita hanya return response bahwa sync telah dijadwalkan
	// Implementasi actual sync bisa menggunakan goroutine atau message queue

	return s.mapHistoryToResponse(&history)
}

// GetSyncHistory mengambil riwayat sinkronisasi dengan pagination
func (s *DapodikService) GetSyncHistory(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)

	// Set filter untuk school_code
	filter := map[string]interface{}{}
	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	page := database.NewPagination[map[string]interface{}]().
		SetRequest(&request).
		SetModal([]school.DapodikSyncHistory{}).
		FindAllPaging()

	return page
}

// GetSyncHistoryDetail mengambil detail satu riwayat sinkronisasi
func (s *DapodikService) GetSyncHistoryDetail(c *gin.Context, id uint) *dapodik_response.DapodikSyncHistoryResponse {
	claims := jwt.GetDataClaims(c)

	var history school.DapodikSyncHistory
	result := s.db.Where("id = ? AND school_code = ?", id, claims.SchoolCode).First(&history)

	if result.Error != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Riwayat sinkronisasi tidak ditemukan"))
	}

	return s.mapHistoryToResponse(&history)
}

// GetSyncSummary mengambil ringkasan sinkronisasi
func (s *DapodikService) GetSyncSummary(c *gin.Context) *dapodik_response.SyncSummaryResponse {
	claims := jwt.GetDataClaims(c)

	var totalSync int64
	s.db.Model(&school.DapodikSyncHistory{}).Where("school_code = ?", claims.SchoolCode).Count(&totalSync)

	var totalSuccess int64
	s.db.Model(&school.DapodikSyncHistory{}).Where("school_code = ? AND status = ?", claims.SchoolCode, school.SyncStatusSuccess).Count(&totalSuccess)

	var totalFailed int64
	s.db.Model(&school.DapodikSyncHistory{}).Where("school_code = ? AND status = ?", claims.SchoolCode, school.SyncStatusFailed).Count(&totalFailed)

	var lastSync school.DapodikSyncHistory
	s.db.Where("school_code = ?", claims.SchoolCode).Order("created_at DESC").First(&lastSync)

	var lastSyncAt *time.Time
	var lastSyncStatus string
	if lastSync.ID != 0 {
		lastSyncAt = lastSync.StartedAt
		lastSyncStatus = string(lastSync.Status)
	}

	return &dapodik_response.SyncSummaryResponse{
		TotalSync:      int(totalSync),
		LastSyncAt:     lastSyncAt,
		LastSyncStatus: lastSyncStatus,
		TotalSuccess:   int(totalSuccess),
		TotalFailed:    int(totalFailed),
	}
}

// Helper functions

func (s *DapodikService) mapHistoryToResponse(history *school.DapodikSyncHistory) *dapodik_response.DapodikSyncHistoryResponse {
	duration := ""
	if history.StartedAt != nil && history.CompletedAt != nil {
		dur := history.CompletedAt.Sub(*history.StartedAt)
		duration = formatDuration(dur)
	}

	return &dapodik_response.DapodikSyncHistoryResponse{
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

func getSyncTypeLabel(syncType school.SyncType) string {
	switch syncType {
	case school.SyncTypeUpdate:
		return "Update & Tambah Data Baru"
	case school.SyncTypeFullReset:
		return "Reset & Import Ulang Semua Data"
	default:
		return string(syncType)
	}
}

func getStatusLabel(status school.SyncStatus) string {
	switch status {
	case school.SyncStatusPending:
		return "Menunggu"
	case school.SyncStatusInProgress:
		return "Sedang Berjalan"
	case school.SyncStatusSuccess:
		return "Berhasil"
	case school.SyncStatusFailed:
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
