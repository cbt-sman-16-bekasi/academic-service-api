package dapodik_request

// SaveDapodikConfigRequest request untuk menyimpan konfigurasi DAPODIK
type SaveDapodikConfigRequest struct {
	DapodikURL string `json:"dapodik_url" binding:"required"`
	APIKey     string `json:"api_key" binding:"required"`
}

// TriggerSyncRequest request untuk memulai sinkronisasi
type TriggerSyncRequest struct {
	SyncType string `json:"sync_type" binding:"required"` // UPDATE atau FULL_RESET
}

// TestConnectionRequest request untuk test koneksi ke DAPODIK
type TestConnectionRequest struct {
	DapodikURL string `json:"dapodik_url" binding:"required"`
	APIKey     string `json:"api_key" binding:"required"`
}
