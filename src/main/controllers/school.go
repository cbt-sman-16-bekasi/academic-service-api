package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/school_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/server/response"
)

type SchoolController struct {
	srv *school_service.SchoolService
}

func NewSchoolController() *SchoolController {
	return &SchoolController{
		srv: school_service.NewSchoolService(),
	}
}

func (s *SchoolController) GetSchool(c *gin.Context) {
	data := s.srv.RetrieveDetailSchool(c)
	response.SuccessResponse("Success get data school_repository", data).Json(c)
}
