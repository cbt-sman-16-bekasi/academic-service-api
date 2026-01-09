package handler

import (
	"regexp"
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/gin-gonic/gin"
)

type SchoolHandler struct {
	service *service.SchoolService
}

// NewSchoolHandler creates handler with injected service
func NewSchoolHandler(svc *service.SchoolService) *SchoolHandler {
	return &SchoolHandler{service: svc}
}

// GetSchool godoc
// @Summary Get school information
// @Description Return school data information
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.DetailSchoolResponse}
// @Router /academic/school [get]
func (h *SchoolHandler) GetSchool(c *gin.Context) {
	data := h.service.GetSchoolDetail(c)
	response.OK(c, "Success get data school", data)
}

// ModifySchool godoc
// @Summary Update school information
// @Description Update school data
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ModifySchoolRequest true "School update request"
// @Success 200 {object} response.BaseResponse{data=dto.DetailSchoolResponse}
// @Router /academic/school/update [put]
func (h *SchoolHandler) ModifySchool(c *gin.Context) {
	var req dto.ModifySchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res := h.service.UpdateSchool(c, req)
	observer.Trigger(model.EventInformationSchoolChanged)
	response.OK(c, "Success update school", res)
}

// GetAllClassCode godoc
// @Summary Get all class codes
// @Description Return class code list (10, 11, 12)
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=[]dto.ClassCodeResponse}
// @Router /academic/class-code [get]
func (h *SchoolHandler) GetAllClassCode(c *gin.Context) {
	data := h.service.GetAllClassCodes(c)
	response.OK(c, "Success get all class code", data)
}

// GetAllSubject godoc
// @Summary Get all subjects
// @Description Return subject list
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=[]dto.SubjectResponse}
// @Router /academic/subjects [get]
func (h *SchoolHandler) GetAllSubject(c *gin.Context) {
	resp := h.service.GetAllSubjects()
	response.OK(c, "Success get all subject", resp)
}

// GetDashboard godoc
// @Summary Get dashboard data
// @Description Return dashboard user data
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse{data=dto.DashboardResponse}
// @Router /academic/dashboard [get]
func (h *SchoolHandler) GetDashboard(c *gin.Context) {
	dt := h.service.GetDashboard(c)
	response.OK(c, "Success get dashboard", dt)
}

// RetrieveConfigSchool godoc
// @Summary Get system configuration
// @Description Return system configuration for login page
// @Tags School
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponse{data=dto.ConfigurationResponse}
// @Router /academic/load/config [get]
func (h *SchoolHandler) RetrieveConfigSchool(c *gin.Context) {
	re := regexp.MustCompile(`^https?://`)
	origin := re.ReplaceAllString(c.Request.Header.Get("Origin"), "")

	res := h.service.LoadConfiguration(origin)
	response.OK(c, "Success retrieve config", res)
}

// GetAllClassSubject godoc
// @Summary Get all class-subject mappings
// @Description Return class-subject mapping list with pagination
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Number of items per page (default: 10)"
// @Success 200 {object} response.BaseResponse{data=database.Paginator}
// @Router /academic/class/subject/all [get]
func (h *SchoolHandler) GetAllClassSubject(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := h.service.GetAllClassSubjects(c, req)
	response.OK(c, "Success get all class subject", data)
}

// GetClassSubject godoc
// @Summary Get class-subject detail
// @Description Return class-subject mapping detail
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class Subject ID"
// @Success 200 {object} response.BaseResponse{data=dto.DetailClassSubjectResponse}
// @Router /academic/class/subject/detail/{id} [get]
func (h *SchoolHandler) GetClassSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := h.service.GetClassSubjectDetail(uint(id))
	response.OK(c, "Success get class subject", detail)
}

// CreateClassSubject godoc
// @Summary Create class-subject mapping
// @Description Create new class-subject mapping
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ModifyClassSubjectRequest true "Class subject data"
// @Success 200 {object} response.BaseResponse{data=dto.DetailClassSubjectResponse}
// @Router /academic/class/subject/create [post]
func (h *SchoolHandler) CreateClassSubject(c *gin.Context) {
	var req dto.ModifyClassSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res := h.service.CreateClassSubject(req)
	observer.Trigger(model.EventSubjectsChanged)
	response.OK(c, "Success create class subject", res)
}

// ModifyClassSubject godoc
// @Summary Update class-subject mapping
// @Description Update existing class-subject mapping
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class Subject ID"
// @Param request body dto.ModifyClassSubjectRequest true "Class subject data"
// @Success 200 {object} response.BaseResponse{data=dto.DetailClassSubjectResponse}
// @Router /academic/class/subject/update/{id} [put]
func (h *SchoolHandler) ModifyClassSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var req dto.ModifyClassSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res := h.service.UpdateClassSubject(uint(id), req)
	observer.Trigger(model.EventSubjectsChanged)
	response.OK(c, "Success update class subject", res)
}

// DeleteClassSubject godoc
// @Summary Delete class-subject mapping
// @Description Delete class-subject mapping
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class Subject ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/class/subject/delete/{id} [delete]
func (h *SchoolHandler) DeleteClassSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	h.service.DeleteClassSubject(uint(id))
	observer.Trigger(model.EventSubjectsChanged)
	response.OK(c, "Success delete class subject", gin.H{"id": id})
}
