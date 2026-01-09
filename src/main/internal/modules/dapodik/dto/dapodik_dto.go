package dto

import "time"

// ----- Request DTOs -----

// SaveDapodikConfigRequest request untuk menyimpan konfigurasi DAPODIK
type SaveDapodikConfigRequest struct {
	DapodikURL  string `json:"dapodik_url" binding:"required"`
	APIKey      string `json:"api_key" binding:"required"`
	UseProxy    bool   `json:"use_proxy"`     // Toggle: true = via proxy, false = langsung ke Dapodik
	ProxyURL    string `json:"proxy_url"`     // Optional: URL proxy server (e.g., http://192.168.1.100:8888)
	ProxyAPIKey string `json:"proxy_api_key"` // Optional: API key untuk akses proxy
}

// TriggerSyncRequest request untuk memulai sinkronisasi
type TriggerSyncRequest struct {
	SyncType string `json:"sync_mode" binding:"required"` // UPDATE atau FULL_RESET
}

// TestConnectionRequest request untuk test koneksi ke DAPODIK
type TestConnectionRequest struct {
	DapodikURL  string `json:"dapodik_url" binding:"required"`
	APIKey      string `json:"api_key" binding:"required"`
	UseProxy    bool   `json:"use_proxy"`     // Toggle: true = via proxy, false = langsung
	ProxyURL    string `json:"proxy_url"`     // Optional: URL proxy server
	ProxyAPIKey string `json:"proxy_api_key"` // Optional: API key untuk akses proxy
}

// ----- Response DTOs -----

// DapodikConfigResponse response untuk konfigurasi DAPODIK
type DapodikConfigResponse struct {
	ID          uint       `json:"id"`
	DapodikURL  string     `json:"dapodik_url"`
	APIKey      string     `json:"api_key"`       // Masked for security
	UseProxy    bool       `json:"use_proxy"`     // Toggle: true = via proxy, false = langsung
	ProxyURL    string     `json:"proxy_url"`     // Proxy server URL (if configured)
	ProxyAPIKey string     `json:"proxy_api_key"` // Masked for security
	IsActive    bool       `json:"is_active"`
	LastSyncAt  *time.Time `json:"last_sync_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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

// =============================================================================
// Preview Response DTOs (for viewing Dapodik data before sync)
// =============================================================================

// PreviewResponse generic response for preview endpoints
type PreviewResponse[T any] struct {
	Total int `json:"total"`
	Data  T   `json:"data"`
}

// PreviewSekolahResponse response for school preview
type PreviewSekolahResponse struct {
	SekolahID        string `json:"sekolah_id"`
	Nama             string `json:"nama"`
	NPSN             string `json:"npsn"`
	NSS              string `json:"nss"`
	BentukPendidikan string `json:"bentuk_pendidikan"`
	StatusSekolah    string `json:"status_sekolah"`
	Alamat           string `json:"alamat"`
	Kelurahan        string `json:"kelurahan"`
	Kecamatan        string `json:"kecamatan"`
	KabupatenKota    string `json:"kabupaten_kota"`
	Provinsi         string `json:"provinsi"`
	Email            string `json:"email"`
	Website          string `json:"website"`
	Telepon          string `json:"telepon"`
}

// PreviewPTKResponse response for teacher preview
type PreviewPTKResponse struct {
	PTKID              string `json:"ptk_id"`
	Nama               string `json:"nama"`
	NUPTK              string `json:"nuptk"`
	NIP                string `json:"nip"`
	NIK                string `json:"nik"`
	JenisKelamin       string `json:"jenis_kelamin"`
	TempatLahir        string `json:"tempat_lahir"`
	TanggalLahir       string `json:"tanggal_lahir"`
	JenisPTK           string `json:"jenis_ptk"`
	JabatanPTK         string `json:"jabatan_ptk"`
	StatusKepegawaian  string `json:"status_kepegawaian"`
	PendidikanTerakhir string `json:"pendidikan_terakhir"`
}

// PreviewRombelResponse response for class preview
type PreviewRombelResponse struct {
	RombelID          string `json:"rombel_id"`
	Nama              string `json:"nama"`
	TingkatPendidikan string `json:"tingkat_pendidikan"`
	Jurusan           string `json:"jurusan"`
	WaliKelas         string `json:"wali_kelas"`
	Kurikulum         string `json:"kurikulum"`
	JumlahSiswa       int    `json:"jumlah_siswa"`
	Semester          string `json:"semester"`
}

// PreviewPesertaDidikResponse response for student preview
type PreviewPesertaDidikResponse struct {
	PesertaDidikID string `json:"peserta_didik_id"`
	NISN           string `json:"nisn"`
	NIPD           string `json:"nipd"`
	Nama           string `json:"nama"`
	JenisKelamin   string `json:"jenis_kelamin"`
	TempatLahir    string `json:"tempat_lahir"`
	TanggalLahir   string `json:"tanggal_lahir"`
	NamaKelas      string `json:"nama_kelas"`
	TingkatKelas   string `json:"tingkat_kelas"`
	NamaAyah       string `json:"nama_ayah"`
	NamaIbu        string `json:"nama_ibu"`
	Email          string `json:"email"`
	NoHP           string `json:"no_hp"`
}

// PreviewPenggunaResponse response for user preview
type PreviewPenggunaResponse struct {
	PenggunaID string `json:"pengguna_id"`
	Username   string `json:"username"`
	Nama       string `json:"nama"`
	Peran      string `json:"peran"`
	NoHP       string `json:"no_hp"`
	PTKID      string `json:"ptk_id,omitempty"`
}
