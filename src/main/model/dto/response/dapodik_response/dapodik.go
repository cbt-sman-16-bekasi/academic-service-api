package dapodik_response

import "time"

// DapodikConfigResponse response untuk konfigurasi DAPODIK
type DapodikConfigResponse struct {
	ID         uint       `json:"id"`
	DapodikURL string     `json:"dapodik_url"`
	APIKey     string     `json:"api_key"` // Masked for security
	IsActive   bool       `json:"is_active"`
	LastSyncAt *time.Time `json:"last_sync_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// DapodikSyncHistoryResponse response untuk riwayat sinkronisasi
type DapodikSyncHistoryResponse struct {
	ID              uint       `json:"id"`
	SyncType        string     `json:"sync_type"`
	SyncTypeLabel   string     `json:"sync_type_label"`
	Status          string     `json:"status"`
	StatusLabel     string     `json:"status_label"`
	StartedAt       *time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	Duration        string     `json:"duration"`
	TotalStudents   int        `json:"total_students"`
	TotalClasses    int        `json:"total_classes"`
	TotalTeachers   int        `json:"total_teachers"`
	SyncedStudents  int        `json:"synced_students"`
	SyncedClasses   int        `json:"synced_classes"`
	SyncedTeachers  int        `json:"synced_teachers"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	TriggeredByName string     `json:"triggered_by_name"`
	CreatedAt       time.Time  `json:"created_at"`
}

// TestConnectionResponse response untuk test koneksi
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Latency string `json:"latency"`
}

// SyncSummaryResponse response untuk summary sinkronisasi
type SyncSummaryResponse struct {
	TotalSync      int        `json:"total_sync"`
	LastSyncAt     *time.Time `json:"last_sync_at"`
	LastSyncStatus string     `json:"last_sync_status"`
	TotalSuccess   int        `json:"total_success"`
	TotalFailed    int        `json:"total_failed"`
}
