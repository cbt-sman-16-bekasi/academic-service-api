package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

// GetAllTypeExam Get all data type exam
// @Summary This endpoint about list of type exam
// @Description Return Pagination Type Exam list data
// @Tags Type Exam
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
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=school.TypeExam}} "Pagination data Type Exam"
// @Router /academic/exam/type-exam/all [get]
func (e *ExamController) GetAllTypeExam(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := e.typeExamService.GetAll(c, request)
	response.SuccessResponse("Success get all type Exam", resp).Json(c)
}

// GetDetailTypeExam Get detail of type exam
// @Summary This endpoint about detail type exam
// @Description Return detail of type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Type Exam ID"
//
// @Success 200 {object} response.BaseResponse{data=school.TypeExam} "Detail response type exam"
// @Router /academic/exam/type-exam/detail/{id} [get]
func (e *ExamController) GetDetailTypeExam(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	res := e.typeExamService.GetDetail(uint(id))
	response.SuccessResponse("Success get type Exam", res).Json(c)
}

// CreateTypeExam Create new type exam
// @Summary This endpoint about Create type exam
// @Description Return Create of type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ModifyTypeExamRequest true "Request body"
//
// @Success 200 {object} response.BaseResponse{data=school.TypeExam} "Detail response type exam"
// @Router /academic/exam/type-exam/create [post]
func (e *ExamController) CreateTypeExam(c *gin.Context) {
	var request exam_request.ModifyTypeExamRequest
	_ = c.BindJSON(&request)

	res := e.typeExamService.CreateTypeExam(request)
	response.SuccessResponse("Success create new type Exam", res).Json(c)
}

// ModifyTypeExam Modify type exam
// @Summary This endpoint about Modify type exam
// @Description Return Modify of type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ModifyTypeExamRequest true "Request body"
// @Param id path int true "Type Exam ID"
//
// @Success 200 {object} response.BaseResponse{data=school.TypeExam} "Detail response type exam"
// @Router /academic/exam/type-exam/update/{id} [put]
func (e *ExamController) ModifyTypeExam(c *gin.Context) {
	var request exam_request.ModifyTypeExamRequest
	_ = c.BindJSON(&request)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	res := e.typeExamService.ModifyTypeExam(uint(id), request)
	response.SuccessResponse("Success update type Exam", res).Json(c)
}

// DeleteTypeExam Create new type exam
// @Summary This endpoint about Create type exam
// @Description Return Create of type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ModifyTypeExamRequest true "Request body"
// @Param id path int true "Type Exam ID"
//
// @Success 200 {object} response.BaseResponse "Delete response type exam"
// @Router /academic/exam/type-exam/delete/{id} [delete]
func (e *ExamController) DeleteTypeExam(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	e.typeExamService.DeleteTypeExam(uint(id))
	response.SuccessResponse("Success delete type Exam", nil).Json(c)
}
