package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// AcademicYearHandler handles HTTP requests for academic year
type AcademicYearHandler struct {
	service *service.AcademicYearService
}

// NewAcademicYearHandler creates a new handler instance
func NewAcademicYearHandler(service *service.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{service: service}
}

// GetAll godoc
// @Summary Get all academic years
// @Description Get all academic years for the current school
// @Tags Academic Year
// @Accept json
// @Produce json
// @Success 200 {object} dto.AcademicYearListResponse
// @Router /academic/academic-year [get]
func (h *AcademicYearHandler) GetAll(c *gin.Context) {
	result := h.service.GetAll(c)
	response.OK(c, "Berhasil mengambil data tahun ajaran", result)
}

// GetActive godoc
// @Summary Get active academic year
// @Description Get the currently active academic year
// @Tags Academic Year
// @Accept json
// @Produce json
// @Success 200 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/active [get]
func (h *AcademicYearHandler) GetActive(c *gin.Context) {
	result := h.service.GetActive(c)
	response.OK(c, "Berhasil mengambil tahun ajaran aktif", result)
}

// GetByID godoc
// @Summary Get academic year by ID
// @Description Get an academic year by its ID
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param id path int true "Academic Year ID"
// @Success 200 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/{id} [get]
func (h *AcademicYearHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	result := h.service.GetByID(c, uint(id))
	response.OK(c, "Berhasil mengambil tahun ajaran", result)
}

// GetBySemesterID godoc
// @Summary Get academic year by semester ID
// @Description Get an academic year by Dapodik semester_id
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param semesterId path string true "Semester ID (e.g., 20251)"
// @Success 200 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/semester/{semesterId} [get]
func (h *AcademicYearHandler) GetBySemesterID(c *gin.Context) {
	semesterID := c.Param("semesterId")
	result := h.service.GetBySemesterID(c, semesterID)
	response.OK(c, "Berhasil mengambil tahun ajaran", result)
}

// Create godoc
// @Summary Create academic year
// @Description Create a new academic year
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param request body dto.CreateAcademicYearRequest true "Create request"
// @Success 201 {object} dto.AcademicYearResponse
// @Router /academic/academic-year [post]
func (h *AcademicYearHandler) Create(c *gin.Context) {
	var req dto.CreateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.Create(c, req)
	response.CreatedResponse(c, "Berhasil membuat tahun ajaran", result)
}

// CreateFromDapodik godoc
// @Summary Create academic year from Dapodik
// @Description Create an academic year from Dapodik semester_id format
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param request body dto.CreateFromDapodikRequest true "Create request"
// @Success 201 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/from-dapodik [post]
func (h *AcademicYearHandler) CreateFromDapodik(c *gin.Context) {
	var req dto.CreateFromDapodikRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.CreateFromDapodik(c, req.SemesterID)
	response.CreatedResponse(c, "Berhasil membuat tahun ajaran dari Dapodik", result)
}

// Update godoc
// @Summary Update academic year
// @Description Update an existing academic year
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param id path int true "Academic Year ID"
// @Param request body dto.UpdateAcademicYearRequest true "Update request"
// @Success 200 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/{id} [put]
func (h *AcademicYearHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.UpdateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.Update(c, uint(id), req)
	response.OK(c, "Berhasil update tahun ajaran", result)
}

// SetActive godoc
// @Summary Set active academic year
// @Description Set an academic year as the active one
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param id path int true "Academic Year ID"
// @Success 200 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/{id}/set-active [post]
func (h *AcademicYearHandler) SetActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	result := h.service.SetActive(c, uint(id))
	response.OK(c, "Tahun ajaran berhasil diaktifkan", result)
}

// SetActiveBySemesterID godoc
// @Summary Set active by semester ID
// @Description Set an academic year as active by semester_id
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param request body dto.SetActiveRequest true "Set active request"
// @Success 200 {object} dto.AcademicYearResponse
// @Router /academic/academic-year/set-active [post]
func (h *AcademicYearHandler) SetActiveBySemesterID(c *gin.Context) {
	var req dto.SetActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.SetActiveBySemesterID(c, req.SemesterID)
	response.OK(c, "Tahun ajaran berhasil diaktifkan", result)
}

// Delete godoc
// @Summary Delete academic year
// @Description Delete an academic year
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param id path int true "Academic Year ID"
// @Success 200
// @Router /academic/academic-year/{id} [delete]
func (h *AcademicYearHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.service.Delete(c, uint(id))
	response.OK(c, "Tahun ajaran berhasil dihapus", nil)
}

// GetOptions godoc
// @Summary Get dropdown options
// @Description Get academic year options for dropdown/select
// @Tags Academic Year
// @Accept json
// @Produce json
// @Success 200 {array} dto.AcademicYearOption
// @Router /academic/academic-year/options [get]
func (h *AcademicYearHandler) GetOptions(c *gin.Context) {
	result := h.service.GetOptions(c)
	response.OK(c, "Berhasil mengambil opsi tahun ajaran", result)
}

// GetStudentHistory godoc
// @Summary Get student class history
// @Description Get class history for a student across semesters
// @Tags Academic Year
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {array} dto.StudentClassHistoryResponse
// @Router /academic/academic-year/student/{studentId}/history [get]
func (h *AcademicYearHandler) GetStudentHistory(c *gin.Context) {
	studentID, _ := strconv.ParseUint(c.Param("studentId"), 10, 32)
	result := h.service.GetStudentHistory(c, uint(studentID))
	response.OK(c, "Berhasil mengambil riwayat kelas siswa", result)
}
