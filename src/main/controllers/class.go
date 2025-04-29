package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/class_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type ClassController struct {
	classService *class_service.ClassService
}

func NewClassController() *ClassController {
	return &ClassController{
		classService: class_service.NewClassService(),
	}
}

// GetAllClass Get data Class
// @Summary This endpoint about list of class
// @Description Return Class list data
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param page query int false "Page number (default: 1)" default(1)
// @Param size query int false "Number of items per page (default: 10)" default(10)
// @Param search.key query string false "Search key field (e.g., name, email, etc.)"
// @Param search.value query string false "Search value for the specified key"
// @Param filter query object false "Filter conditions (varies depending on the API)"
//
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=school.Class}} "Pagination data class"
// @Router /academic/class/all [get]
func (s *ClassController) GetAllClass(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	var data *database.Paginator = s.classService.FindAllClass(req)
	response.SuccessResponse("Success get all class", data).Json(c)
}

// GetDetailClass Get detail of Class
// @Summary This endpoint about detail of class
// @Description Return detail information class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int false "Class ID"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassResponse} "Detail of class"
// @Router /academic/class/detail/{id} [get]
func (s *ClassController) GetDetailClass(c *gin.Context) {

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := s.classService.GetDetailClass(uint(id))
	response.SuccessResponse("Success get detail class", detail).Json(c)
}

// CreateNewClass Create New Class
// @Summary This endpoint about create new class
// @Description Create New Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body class_request.ModifyClassRequest true "Body Request of Create New Class"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassResponse} "Create class response"
// @Router /academic/class/create [post]
func (s *ClassController) CreateNewClass(c *gin.Context) {
	var req class_request.ModifyClassRequest
	_ = c.BindJSON(&req)

	resp := s.classService.CreateNewClass(req)
	observer.Trigger(model.EventClassChanged)
	response.SuccessResponse("Success create new class", resp).Json(c)
}

// UpdateClass Update Class
// @Summary This endpoint about Update class
// @Description Update Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Class ID"
// @Param request body class_request.ModifyClassRequest true "Body Request of Update Class"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassResponse} "Update class response"
// @Router /academic/class/update/{id} [put]
func (s *ClassController) UpdateClass(c *gin.Context) {
	var req class_request.ModifyClassRequest
	_ = c.BindJSON(&req)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.classService.ModifyClass(uint(id), req)
	observer.Trigger(model.EventClassChanged)
	response.SuccessResponse("Success update class", resp).Json(c)
}

// DeleteClass Delete New Class
// @Summary This endpoint about Delete class
// @Description Delete Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Class ID"
//
// @Success 200 {object} response.BaseResponse "Delete class response"
// @Router /academic/class/delete/{id} [delete]
func (s *ClassController) DeleteClass(c *gin.Context) {

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.classService.DeleteById(uint(id))
	observer.Trigger(model.EventClassChanged)
	response.SuccessResponse("Success delete class", gin.H{"id": id}).Json(c)
}
