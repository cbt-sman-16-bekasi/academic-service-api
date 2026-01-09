package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type DapodikHandler struct {
	service *service.DapodikService
}

// NewDapodikHandler creates handler with injected service
func NewDapodikHandler(svc *service.DapodikService) *DapodikHandler {
	return &DapodikHandler{service: svc}
}

// GetDapodikConfig mengambil konfigurasi DAPODIK
// @Summary Get DAPODIK configuration
// @Description Mengambil konfigurasi koneksi ke API DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.DapodikConfigResponse} "DAPODIK config response"
// @Router /academic/dapodik/config [get]
func (h *DapodikHandler) GetDapodikConfig(c *gin.Context) {
	config := h.service.GetConfig(c)
	response.OK(c, "Berhasil mengambil konfigurasi DAPODIK", config)
}

// SaveDapodikConfig menyimpan konfigurasi DAPODIK
// @Summary Save DAPODIK configuration
// @Description Menyimpan atau update konfigurasi koneksi ke API DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SaveDapodikConfigRequest true "Body Request untuk menyimpan config"
// @Success 200 {object} response.BaseResponse{data=dto.DapodikConfigResponse} "DAPODIK config response"
// @Router /academic/dapodik/config [post]
func (h *DapodikHandler) SaveDapodikConfig(c *gin.Context) {
	var request dto.SaveDapodikConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	config := h.service.SaveConfig(c, request)
	response.OK(c, "Berhasil menyimpan konfigurasi DAPODIK", config)
}

// TestDapodikConnection menguji koneksi ke API DAPODIK
// @Summary Test DAPODIK connection
// @Description Menguji koneksi ke server DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TestConnectionRequest true "Body Request untuk test koneksi"
// @Success 200 {object} response.BaseResponse{data=dto.TestConnectionResponse} "Test connection response"
// @Router /academic/dapodik/test-connection [post]
func (h *DapodikHandler) TestDapodikConnection(c *gin.Context) {
	var request dto.TestConnectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.TestConnection(request)
	response.OK(c, "Test koneksi selesai", result)
}

// TriggerSync memulai sinkronisasi dengan DAPODIK
// @Summary Trigger DAPODIK sync
// @Description Memulai proses sinkronisasi data dengan DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TriggerSyncRequest true "Body Request untuk trigger sync"
// @Success 200 {object} response.BaseResponse{data=dto.DapodikSyncHistoryResponse} "Sync triggered response"
// @Router /academic/dapodik/sync [post]
func (h *DapodikHandler) TriggerSync(c *gin.Context) {
	var request dto.TriggerSyncRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.TriggerSync(c, request)
	response.OK(c, "Sinkronisasi telah dijadwalkan", result)
}

// GetSyncHistory mengambil riwayat sinkronisasi
// @Summary Get sync history
// @Description Mengambil riwayat sinkronisasi dengan pagination
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)" default(1)
// @Param size query int false "Number of items per page (default: 10)" default(10)
// @Success 200 {object} response.BaseResponse{data=database.Paginator} "Sync history with pagination"
// @Router /academic/dapodik/history [get]
func (h *DapodikHandler) GetSyncHistory(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := h.service.GetSyncHistory(c, req)
	response.OK(c, "Berhasil mengambil riwayat sinkronisasi", data)
}

// GetSyncHistoryDetail mengambil detail riwayat sinkronisasi
// @Summary Get sync history detail
// @Description Mengambil detail satu riwayat sinkronisasi
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "History ID"
// @Success 200 {object} response.BaseResponse{data=dto.DapodikSyncHistoryResponse} "Sync history detail"
// @Router /academic/dapodik/history/{id} [get]
func (h *DapodikHandler) GetSyncHistoryDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	data := h.service.GetSyncHistoryDetail(c, uint(id))
	response.OK(c, "Berhasil mengambil detail riwayat sinkronisasi", data)
}

// GetSyncSummary mengambil ringkasan sinkronisasi
// @Summary Get sync summary
// @Description Mengambil ringkasan statistik sinkronisasi
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.SyncSummaryResponse} "Sync summary response"
// @Router /academic/dapodik/summary [get]
func (h *DapodikHandler) GetSyncSummary(c *gin.Context) {
	data := h.service.GetSyncSummary(c)
	response.OK(c, "Berhasil mengambil ringkasan sinkronisasi", data)
}

// =============================================================================
// Preview Endpoints - View data from Dapodik before syncing
// =============================================================================

// PreviewSekolah mengambil data sekolah dari Dapodik
// @Summary Preview sekolah data from Dapodik
// @Description Mengambil data sekolah dari API Dapodik tanpa menyimpan
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.PreviewSekolahResponse} "Sekolah data"
// @Router /academic/dapodik/preview/sekolah [get]
func (h *DapodikHandler) PreviewSekolah(c *gin.Context) {
	data := h.service.PreviewSekolah(c)
	response.OK(c, "Berhasil mengambil data sekolah dari Dapodik", data)
}

// PreviewPTK mengambil data PTK/guru dari Dapodik
// @Summary Preview PTK data from Dapodik
// @Description Mengambil data PTK (guru) dari API Dapodik tanpa menyimpan
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.PreviewResponse[[]dto.PreviewPTKResponse]} "PTK data"
// @Router /academic/dapodik/preview/ptk [get]
func (h *DapodikHandler) PreviewPTK(c *gin.Context) {
	data := h.service.PreviewPTK(c)
	response.OK(c, "Berhasil mengambil data PTK dari Dapodik", data)
}

// PreviewRombonganBelajar mengambil data rombongan belajar dari Dapodik
// @Summary Preview rombongan belajar data from Dapodik
// @Description Mengambil data rombongan belajar (kelas) dari API Dapodik tanpa menyimpan
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.PreviewResponse[[]dto.PreviewRombelResponse]} "Rombel data"
// @Router /academic/dapodik/preview/rombel [get]
func (h *DapodikHandler) PreviewRombonganBelajar(c *gin.Context) {
	data := h.service.PreviewRombonganBelajar(c)
	response.OK(c, "Berhasil mengambil data rombongan belajar dari Dapodik", data)
}

// PreviewPesertaDidik mengambil data peserta didik dari Dapodik
// @Summary Preview peserta didik data from Dapodik
// @Description Mengambil data peserta didik (siswa) dari API Dapodik tanpa menyimpan
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.PreviewResponse[[]dto.PreviewPesertaDidikResponse]} "Peserta didik data"
// @Router /academic/dapodik/preview/peserta-didik [get]
func (h *DapodikHandler) PreviewPesertaDidik(c *gin.Context) {
	data := h.service.PreviewPesertaDidik(c)
	response.OK(c, "Berhasil mengambil data peserta didik dari Dapodik", data)
}

// PreviewPengguna mengambil data pengguna dari Dapodik
// @Summary Preview pengguna data from Dapodik
// @Description Mengambil data pengguna dari API Dapodik tanpa menyimpan
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.PreviewResponse[[]dto.PreviewPenggunaResponse]} "Pengguna data"
// @Router /academic/dapodik/preview/pengguna [get]
func (h *DapodikHandler) PreviewPengguna(c *gin.Context) {
	data := h.service.PreviewPengguna(c)
	response.OK(c, "Berhasil mengambil data pengguna dari Dapodik", data)
}
