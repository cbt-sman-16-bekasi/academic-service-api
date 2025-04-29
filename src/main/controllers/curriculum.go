package controllers

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/curriculum_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/curriculum_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type CurriculumController struct {
	service *curriculum_service.SubjectService
}

func NewCurriculumController() *CurriculumController {
	return &CurriculumController{
		service: curriculum_service.NewSubjectService(),
	}
}

func (s *CurriculumController) GetAllSubject(c *gin.Context) {

	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := s.service.GetAllSubject(request)
	response.SuccessResponse("Success get data subject", resp).Json(c)
}

func (s *CurriculumController) GetSubject(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.service.GetSubject(uint64(id))
	response.SuccessResponse("Success get data subject", resp).Json(c)
}

func (s *CurriculumController) CreateSubject(c *gin.Context) {
	var subject curriculum_request.SubjectRequest
	_ = c.BindJSON(&subject)

	resp := s.service.CreateSubject(subject)
	observer.Trigger(model.EventSubjectsChanged)
	response.SuccessResponse("Success create subject", resp).Json(c)
}

func (s *CurriculumController) DeleteSubject(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.service.DeleteSubject(uint64(id))
	observer.Trigger(model.EventSubjectsChanged)
	response.SuccessResponse(fmt.Sprint("Success delete subject", id), "").Json(c)
}

func (s *CurriculumController) UpdateSubject(c *gin.Context) {
	var subject curriculum_request.SubjectRequest
	_ = c.BindJSON(&subject)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.service.UpdateSubject(uint64(id), subject)
	observer.Trigger(model.EventSubjectsChanged)
	response.SuccessResponse("Success update subject", resp).Json(c)
}
