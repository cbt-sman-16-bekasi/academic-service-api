package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// DapodikErrorResponse represents error response from Dapodik API
type DapodikErrorResponse struct {
	Success    bool   `json:"success"`
	HTTPCode   int    `json:"http_code"`
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
}

// DapodikClient handles communication with Dapodik API
type DapodikClient struct {
	baseURL     string
	apiKey      string
	npsn        string // NPSN sekolah - required for all API calls
	useProxy    bool   // Toggle: true = via proxy, false = direct
	proxyURL    string // Optional: proxy server URL
	proxyAPIKey string // Optional: API key for proxy authentication
	httpClient  *http.Client
}

// DapodikClientConfig holds configuration for creating a DapodikClient
type DapodikClientConfig struct {
	BaseURL     string // Dapodik server URL
	APIKey      string // Dapodik API token
	NPSN        string // School NPSN
	UseProxy    bool   // Toggle: true = via proxy, false = direct to Dapodik
	ProxyURL    string // Optional: proxy server URL (e.g., http://192.168.1.100:8888)
	ProxyAPIKey string // Optional: API key for proxy server
}

// NewDapodikClient creates a new Dapodik API client
func NewDapodikClient(baseURL, apiKey, npsn string) *DapodikClient {
	return &DapodikClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		npsn:    npsn,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for large data
		},
	}
}

// NewDapodikClientWithConfig creates a new Dapodik API client with full configuration
func NewDapodikClientWithConfig(cfg DapodikClientConfig) *DapodikClient {
	return &DapodikClient{
		baseURL:     cfg.BaseURL,
		apiKey:      cfg.APIKey,
		npsn:        cfg.NPSN,
		useProxy:    cfg.UseProxy,
		proxyURL:    cfg.ProxyURL,
		proxyAPIKey: cfg.ProxyAPIKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// =============================================================================
// Response Structures
// =============================================================================

// BaseResponse is the common response structure from Dapodik API
type BaseResponse[T any] struct {
	Results int    `json:"results"`
	ID      string `json:"id"`
	Start   int    `json:"start"`
	Limit   int    `json:"limit"`
	Rows    T      `json:"rows"`
}

// Pengguna represents a user from Dapodik
type Pengguna struct {
	PenggunaID     string  `json:"pengguna_id"`
	SekolahID      string  `json:"sekolah_id"`
	Username       string  `json:"username"`
	Nama           string  `json:"nama"`
	PeranIDStr     string  `json:"peran_id_str"`
	Password       string  `json:"password"`
	Alamat         *string `json:"alamat"`
	NoTelepon      *string `json:"no_telepon"`
	NoHP           string  `json:"no_hp"`
	PTKID          *string `json:"ptk_id"`
	PesertaDidikID *string `json:"peserta_didik_id"`
}

// Sekolah represents school data from Dapodik
type Sekolah struct {
	SekolahID           string  `json:"sekolah_id"`
	Nama                string  `json:"nama"`
	NSS                 string  `json:"nss"`
	NPSN                string  `json:"npsn"`
	BentukPendidikanID  int     `json:"bentuk_pendidikan_id"`
	BentukPendidikanStr string  `json:"bentuk_pendidikan_id_str"`
	StatusSekolah       string  `json:"status_sekolah"`
	StatusSekolahStr    string  `json:"status_sekolah_str"`
	AlamatJalan         string  `json:"alamat_jalan"`
	RT                  string  `json:"rt"`
	RW                  string  `json:"rw"`
	KodeWilayah         string  `json:"kode_wilayah"`
	KodePos             string  `json:"kode_pos"`
	NomorTelepon        string  `json:"nomor_telepon"`
	NomorFax            *string `json:"nomor_fax"`
	Email               string  `json:"email"`
	Website             string  `json:"website"`
	IsSKS               bool    `json:"is_sks"`
	Lintang             string  `json:"lintang"`
	Bujur               string  `json:"bujur"`
	Dusun               string  `json:"dusun"`
	DesaKelurahan       string  `json:"desa_kelurahan"`
	Kecamatan           string  `json:"kecamatan"`
	KabupatenKota       string  `json:"kabupaten_kota"`
	Provinsi            string  `json:"provinsi"`
}

// RombonganBelajar represents class/group from Dapodik
type RombonganBelajar struct {
	RombonganBelajarID   string          `json:"rombongan_belajar_id"`
	Nama                 string          `json:"nama"`
	TingkatPendidikanID  string          `json:"tingkat_pendidikan_id"`
	TingkatPendidikanStr string          `json:"tingkat_pendidikan_id_str"`
	SemesterID           string          `json:"semester_id"`
	JenisRombel          string          `json:"jenis_rombel"`
	JenisRombelStr       string          `json:"jenis_rombel_str"`
	KurikulumID          int             `json:"kurikulum_id"`
	KurikulumStr         string          `json:"kurikulum_id_str"`
	IDRuang              string          `json:"id_ruang"`
	IDRuangStr           string          `json:"id_ruang_str"`
	MovingClass          string          `json:"moving_class"`
	PTKID                string          `json:"ptk_id"`
	PTKIDStr             string          `json:"ptk_id_str"` // Wali Kelas name
	JurusanID            string          `json:"jurusan_id"`
	JurusanStr           string          `json:"jurusan_id_str"`
	AnggotaRombel        []AnggotaRombel `json:"anggota_rombel"`
}

// AnggotaRombel represents class member
type AnggotaRombel struct {
	AnggotaRombelID     string `json:"anggota_rombel_id"`
	PesertaDidikID      string `json:"peserta_didik_id"`
	JenisPendaftaranID  string `json:"jenis_pendaftaran_id"`
	JenisPendaftaranStr string `json:"jenis_pendaftaran_id_str"`
}

// PTK represents teacher from Dapodik (Pendidik dan Tenaga Kependidikan)
type PTK struct {
	TahunAjaranID           string              `json:"tahun_ajaran_id"`
	PTKTerdaftarID          string              `json:"ptk_terdaftar_id"`
	PTKID                   string              `json:"ptk_id"`
	PTKInduk                string              `json:"ptk_induk"`
	TanggalSuratTugas       string              `json:"tanggal_surat_tugas"`
	Nama                    string              `json:"nama"`
	JenisKelamin            string              `json:"jenis_kelamin"`
	TempatLahir             string              `json:"tempat_lahir"`
	TanggalLahir            string              `json:"tanggal_lahir"`
	AgamaID                 int                 `json:"agama_id"`
	AgamaStr                string              `json:"agama_id_str"`
	NUPTK                   string              `json:"nuptk"`
	NIK                     string              `json:"nik"`
	JenisPTKID              string              `json:"jenis_ptk_id"`
	JenisPTKStr             string              `json:"jenis_ptk_id_str"`
	JabatanPTKID            string              `json:"jabatan_ptk_id"`
	JabatanPTKStr           string              `json:"jabatan_ptk_id_str"`
	StatusKepegawaianID     int                 `json:"status_kepegawaian_id"`
	StatusKepegawaianStr    string              `json:"status_kepegawaian_id_str"`
	NIP                     string              `json:"nip"`
	PendidikanTerakhir      string              `json:"pendidikan_terakhir"`
	BidangStudiTerakhir     string              `json:"bidang_studi_terakhir"`
	PangkatGolonganTerakhir string              `json:"pangkat_golongan_terakhir"`
	RwyPendFormal           []RiwayatPendidikan `json:"rwy_pend_formal"`
}

// RiwayatPendidikan represents education history
type RiwayatPendidikan struct {
	RiwayatPendidikanFormalID string `json:"riwayat_pendidikan_formal_id"`
	SatuanPendidikanFormal    string `json:"satuan_pendidikan_formal"`
	Fakultas                  string `json:"fakultas"`
	Kependidikan              string `json:"kependidikan"`
	TahunMasuk                string `json:"tahun_masuk"`
	TahunLulus                string `json:"tahun_lulus"`
	NIM                       string `json:"nim"`
	StatusKuliah              string `json:"status_kuliah"`
	IPK                       string `json:"ipk"`
	BidangStudiStr            string `json:"bidang_studi_id_str"`
	JenjangPendidikanStr      string `json:"jenjang_pendidikan_id_str"`
	GelarAkademikStr          string `json:"gelar_akademik_id_str"`
}

// PesertaDidik represents student from Dapodik
type PesertaDidik struct {
	RegistrasiID        string  `json:"registrasi_id"`
	JenisPendaftaranID  string  `json:"jenis_pendaftaran_id"`
	JenisPendaftaranStr string  `json:"jenis_pendaftaran_id_str"`
	NIPD                string  `json:"nipd"` // Nomor Induk Peserta Didik (local school ID)
	TanggalMasukSekolah string  `json:"tanggal_masuk_sekolah"`
	SekolahAsal         *string `json:"sekolah_asal"`
	PesertaDidikID      string  `json:"peserta_didik_id"`
	Nama                string  `json:"nama"`
	NISN                string  `json:"nisn"`
	JenisKelamin        string  `json:"jenis_kelamin"`
	NIK                 string  `json:"nik"`
	TempatLahir         string  `json:"tempat_lahir"`
	TanggalLahir        string  `json:"tanggal_lahir"`
	AgamaID             int     `json:"agama_id"`
	AgamaStr            string  `json:"agama_id_str"`
	NomorTeleponRumah   *string `json:"nomor_telepon_rumah"`
	NomorTeleponSeluler *string `json:"nomor_telepon_seluler"`
	NamaAyah            *string `json:"nama_ayah"`
	PekerjaanAyahID     int     `json:"pekerjaan_ayah_id"`
	PekerjaanAyahStr    string  `json:"pekerjaan_ayah_id_str"`
	NamaIbu             *string `json:"nama_ibu"`
	PekerjaanIbuID      int     `json:"pekerjaan_ibu_id"`
	PekerjaanIbuStr     string  `json:"pekerjaan_ibu_id_str"`
	NamaWali            *string `json:"nama_wali"`
	PekerjaanWaliID     *int    `json:"pekerjaan_wali_id"`
	PekerjaanWaliStr    string  `json:"pekerjaan_wali_id_str"`
	AnakKeberapa        string  `json:"anak_keberapa"`
	TinggiBadan         string  `json:"tinggi_badan"`
	BeratBadan          string  `json:"berat_badan"`
	Email               *string `json:"email"`
	SemesterID          string  `json:"semester_id"`
	AnggotaRombelID     string  `json:"anggota_rombel_id"`
	RombonganBelajarID  string  `json:"rombongan_belajar_id"`
	TingkatPendidikanID string  `json:"tingkat_pendidikan_id"`
	NamaRombel          string  `json:"nama_rombel"`
	KurikulumID         int     `json:"kurikulum_id"`
	KurikulumStr        string  `json:"kurikulum_id_str"`
	KebutuhanKhusus     string  `json:"kebutuhan_khusus"`
}

// =============================================================================
// API Methods
// =============================================================================

// doRequest performs HTTP request to Dapodik API (directly or via proxy)
func (c *DapodikClient) doRequest(endpoint string) ([]byte, error) {
	// Build target URL with NPSN query parameter
	targetURL := fmt.Sprintf("%s/%s?npsn=%s", c.baseURL, endpoint, c.npsn)

	var requestURL string
	var req *http.Request
	var err error
	var viaProxy bool

	// Check if proxy is enabled and configured
	if c.useProxy && c.proxyURL != "" {
		// Use proxy: send request to proxy with target URL as query parameter
		requestURL = fmt.Sprintf("%s/proxy?target=%s", c.proxyURL, targetURL)
		viaProxy = true
		log.Debug().
			Str("proxy", c.proxyURL).
			Str("target", targetURL).
			Str("npsn", c.npsn).
			Msg("Calling Dapodik API via proxy")
	} else {
		// Direct call to Dapodik
		requestURL = targetURL
		viaProxy = false
		log.Debug().Str("url", requestURL).Str("npsn", c.npsn).Msg("Calling Dapodik API directly")
	}

	req, err = http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Dapodik authorization header
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	// Set proxy API key if using proxy
	if viaProxy && c.proxyAPIKey != "" {
		req.Header.Set("X-API-Key", c.proxyAPIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid API key")
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle malformed response where HTTP headers are included in body
	// e.g., "HTTP/1.0 403 Forbidden\r\nCache-Control: no-cache\r\n...\r\n\r\n{json}"
	body = c.extractJSONFromBody(body)

	// Check if response is an error response from Dapodik
	if err := c.checkForErrorResponse(body); err != nil {
		return nil, err
	}

	return body, nil
}

// extractJSONFromBody extracts JSON from response body that may contain HTTP headers as text
// Dapodik API sometimes returns responses like:
// "HTTP/1.0 403 Forbidden\r\nCache-Control: no-cache\r\n...\r\n\r\n{...json...}"
func (c *DapodikClient) extractJSONFromBody(body []byte) []byte {
	// Check if body starts with "HTTP/" indicating malformed response
	if bytes.HasPrefix(body, []byte("HTTP/")) {
		log.Warn().Msg("Detected malformed Dapodik response with HTTP headers in body")

		// Find the JSON start by looking for first '{' or '['
		jsonStart := bytes.IndexAny(body, "{[")
		if jsonStart != -1 {
			return body[jsonStart:]
		}
	}

	return body
}

// checkForErrorResponse checks if the response body contains a Dapodik error
func (c *DapodikClient) checkForErrorResponse(body []byte) error {
	// Quick check if this looks like an error response
	if !bytes.Contains(body, []byte(`"success":false`)) &&
		!bytes.Contains(body, []byte(`"success": false`)) {
		return nil
	}

	var errResp DapodikErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// Not a valid error response format, continue normally
		return nil
	}

	if !errResp.Success {
		log.Error().
			Int("http_code", errResp.HTTPCode).
			Str("status_code", errResp.StatusCode).
			Str("message", errResp.Message).
			Msg("Dapodik API returned error")

		// Return specific error messages based on error type
		switch errResp.HTTPCode {
		case 403:
			return fmt.Errorf("dapodik: %s - Pastikan aplikasi sudah terdaftar di Web Service Dapodik", errResp.Message)
		case 401:
			return fmt.Errorf("dapodik: %s - Token/API Key tidak valid", errResp.Message)
		default:
			return fmt.Errorf("dapodik [%d]: %s", errResp.HTTPCode, errResp.Message)
		}
	}

	return nil
}

// GetPengguna fetches all users from Dapodik
func (c *DapodikClient) GetPengguna() ([]Pengguna, int, error) {
	body, err := c.doRequest("WebService/getPengguna")
	if err != nil {
		return nil, 0, err
	}

	var resp BaseResponse[[]Pengguna]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Rows, resp.Results, nil
}

// GetSekolah fetches school info from Dapodik
func (c *DapodikClient) GetSekolah() (*Sekolah, error) {
	body, err := c.doRequest("WebService/getSekolah")
	if err != nil {
		return nil, err
	}

	// Sekolah response has single object in rows, not array
	var resp struct {
		Results int     `json:"results"`
		Rows    Sekolah `json:"rows"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Rows, nil
}

// GetRombonganBelajar fetches all classes from Dapodik
func (c *DapodikClient) GetRombonganBelajar() ([]RombonganBelajar, int, error) {
	body, err := c.doRequest("WebService/getRombonganBelajar")
	if err != nil {
		return nil, 0, err
	}

	var resp BaseResponse[[]RombonganBelajar]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Rows, resp.Results, nil
}

// GetPTK fetches all teachers from Dapodik
func (c *DapodikClient) GetPTK() ([]PTK, int, error) {
	body, err := c.doRequest("WebService/getGtk")
	if err != nil {
		return nil, 0, err
	}

	var resp BaseResponse[[]PTK]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Rows, resp.Results, nil
}

// GetPesertaDidik fetches all students from Dapodik
func (c *DapodikClient) GetPesertaDidik() ([]PesertaDidik, int, error) {
	body, err := c.doRequest("WebService/getPesertaDidik")
	if err != nil {
		return nil, 0, err
	}

	var resp BaseResponse[[]PesertaDidik]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Rows, resp.Results, nil
}
