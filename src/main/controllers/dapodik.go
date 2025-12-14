package controllers

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/dapodik_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/dapodik_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
)

type DapodikController struct {
	srv *dapodik_service.DapodikService
}

func NewDapodikController() *DapodikController {
	return &DapodikController{
		srv: dapodik_service.NewDapodikService(),
	}
}

// GetDapodikConfig mengambil konfigurasi DAPODIK
// @Summary Get DAPODIK configuration
// @Description Mengambil konfigurasi koneksi ke API DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} response.BaseResponse{data=dapodik_response.DapodikConfigResponse} "DAPODIK config response"
// @Router /academic/dapodik/config [get]
func (d *DapodikController) GetDapodikConfig(c *gin.Context) {
	config := d.srv.GetConfig(c)
	response.SuccessResponse("Berhasil mengambil konfigurasi DAPODIK", config).Json(c)
}

// SaveDapodikConfig menyimpan konfigurasi DAPODIK
// @Summary Save DAPODIK configuration
// @Description Menyimpan atau update konfigurasi koneksi ke API DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body dapodik_request.SaveDapodikConfigRequest true "Body Request untuk menyimpan config"
//
// @Success 200 {object} response.BaseResponse{data=dapodik_response.DapodikConfigResponse} "DAPODIK config response"
// @Router /academic/dapodik/config [post]
func (d *DapodikController) SaveDapodikConfig(c *gin.Context) {
	var request dapodik_request.SaveDapodikConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ErrorResponse(400, "Request tidak valid: "+err.Error(), nil).Json(c)
		return
	}

	config := d.srv.SaveConfig(c, request)
	response.SuccessResponse("Berhasil menyimpan konfigurasi DAPODIK", config).Json(c)
}

// TestDapodikConnection menguji koneksi ke API DAPODIK
// @Summary Test DAPODIK connection
// @Description Menguji koneksi ke server DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body dapodik_request.TestConnectionRequest true "Body Request untuk test koneksi"
//
// @Success 200 {object} response.BaseResponse{data=dapodik_response.TestConnectionResponse} "Test connection response"
// @Router /academic/dapodik/test-connection [post]
func (d *DapodikController) TestDapodikConnection(c *gin.Context) {
	var request dapodik_request.TestConnectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ErrorResponse(400, "Request tidak valid: "+err.Error(), nil).Json(c)
		return
	}

	result := d.srv.TestConnection(request)
	response.SuccessResponse("Test koneksi selesai", result).Json(c)
}

// TriggerSync memulai sinkronisasi dengan DAPODIK
// @Summary Trigger DAPODIK sync
// @Description Memulai proses sinkronisasi data dengan DAPODIK
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body dapodik_request.TriggerSyncRequest true "Body Request untuk trigger sync"
//
// @Success 200 {object} response.BaseResponse{data=dapodik_response.DapodikSyncHistoryResponse} "Sync triggered response"
// @Router /academic/dapodik/sync [post]
func (d *DapodikController) TriggerSync(c *gin.Context) {
	var request dapodik_request.TriggerSyncRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ErrorResponse(400, "Request tidak valid: "+err.Error(), nil).Json(c)
		return
	}

	result := d.srv.TriggerSync(c, request)
	response.SuccessResponse("Sinkronisasi telah dijadwalkan", result).Json(c)
}

// GetSyncHistory mengambil riwayat sinkronisasi
// @Summary Get sync history
// @Description Mengambil riwayat sinkronisasi dengan pagination
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param page query int false "Page number (default: 1)" default(1)
// @Param size query int false "Number of items per page (default: 10)" default(10)
//
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=school.DapodikSyncHistory}} "Sync history with pagination"
// @Router /academic/dapodik/history [get]
func (d *DapodikController) GetSyncHistory(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := d.srv.GetSyncHistory(c, req)
	response.SuccessResponse("Berhasil mengambil riwayat sinkronisasi", data).Json(c)
}

// GetSyncHistoryDetail mengambil detail riwayat sinkronisasi
// @Summary Get sync history detail
// @Description Mengambil detail satu riwayat sinkronisasi
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "History ID"
//
// @Success 200 {object} response.BaseResponse{data=dapodik_response.DapodikSyncHistoryResponse} "Sync history detail"
// @Router /academic/dapodik/history/{id} [get]
func (d *DapodikController) GetSyncHistoryDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	data := d.srv.GetSyncHistoryDetail(c, uint(id))
	response.SuccessResponse("Berhasil mengambil detail riwayat sinkronisasi", data).Json(c)
}

// GetSyncSummary mengambil ringkasan sinkronisasi
// @Summary Get sync summary
// @Description Mengambil ringkasan statistik sinkronisasi
// @Tags Dapodik
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} response.BaseResponse{data=dapodik_response.SyncSummaryResponse} "Sync summary response"
// @Router /academic/dapodik/summary [get]
func (d *DapodikController) GetSyncSummary(c *gin.Context) {
	data := d.srv.GetSyncSummary(c)
	response.SuccessResponse("Berhasil mengambil ringkasan sinkronisasi", data).Json(c)
}
