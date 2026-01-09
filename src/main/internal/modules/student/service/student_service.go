package service

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/repository"
	userService "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	sharedHelper "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type StudentService struct {
	repo        *repository.StudentRepository
	userService *userService.UserService
}

// NewStudentService creates service with injected dependencies
func NewStudentService(repo *repository.StudentRepository, userSvc *userService.UserService) *StudentService {
	return &StudentService{
		repo:        repo,
		userService: userSvc,
	}
}

// AllStudent returns paginated students
func (s *StudentService) AllStudent(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)

	// Filter by school_code for multi-tenant isolation
	filter := map[string]interface{}{}
	if request.Filter != nil {
		filter = *request.Filter
	}
	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	return database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetRequest(&request).
		SetModel([]view.VStudent{}).
		FindAllPaging()
}

// DetailStudent returns student detail
func (s *StudentService) DetailStudent(id uint) dto.DetailStudentResponse {
	detail, err := s.repo.FindByID(id)
	if err != nil || detail.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Student not found"))
	}

	return dto.DetailStudentResponse{
		Nisn:   detail.NISN,
		Name:   detail.Name,
		Gender: detail.Gender,
		Class: dto.GeneralLabelKeyResponse{
			Key:   detail.ClassID,
			Label: detail.ClassName,
		},
	}
}

// CreateStudent creates a new student
func (s *StudentService) CreateStudent(c *gin.Context, req dto.StudentModifyRequest) dto.DetailStudentResponse {
	claims := jwt.GetDataClaims(c)

	// Get STUDENT role
	role := s.userService.GetRole("STUDENT")
	if role == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Undefined role 'STUDENT'"))
	}

	// Check if NISN already exists
	existing, _ := s.repo.FindByNISN(req.Nisn)
	if existing != nil && existing.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User with NISN '%s' already exists", req.Nisn)))
	}

	// Create user first
	resultUser := s.userService.CreateNewUser(&user.User{
		Username:   req.Nisn,
		Role:       role.ID,
		Status:     1,
		SchoolCode: claims.SchoolCode,
	})

	// Create student
	newStudent := student.Student{
		UserId: resultUser.ID,
		Name:   req.Name,
		Nisn:   req.Nisn,
		Gender: req.Gender,
	}

	if err := s.repo.Create(&newStudent); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Create student class
	studentClass := student.StudentClass{
		StudentId: newStudent.ID,
		ClassId:   req.ClassId,
	}
	if err := s.repo.CreateStudentClass(&studentClass); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.DetailStudentResponse{
		Nisn:   req.Nisn,
		Name:   req.Name,
		Gender: req.Gender,
		Class: dto.GeneralLabelKeyResponse{
			Key:   role.Code,
			Label: role.Name,
		},
	}
}

// UpdateStudent updates an existing student
func (s *StudentService) UpdateStudent(id uint, req dto.StudentModifyRequest) dto.DetailStudentResponse {
	existingStudent, err := s.repo.FindByID(id)
	if err != nil || existingStudent.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Student not found"))
	}

	// Check NISN uniqueness if changed
	if existingStudent.NISN != req.Nisn {
		duplicate, _ := s.repo.FindByNISN(req.Nisn)
		if duplicate != nil && duplicate.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User with NISN '%s' already exists", req.Nisn)))
		}
	}

	// Update student
	studentUpdate := &student.Student{
		Name:   req.Name,
		Nisn:   req.Nisn,
		Gender: req.Gender,
	}
	if err := s.repo.Update(id, studentUpdate); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Update student class
	if err := s.repo.UpdateStudentClass(id, req.ClassId); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Update user username if NISN changed
	if existingStudent.NISN != req.Nisn {
		sc, _ := s.repo.FindStudentClassByStudentID(id)
		if sc != nil && sc.DetailStudent.DetailUser.ID != 0 {
			sc.DetailStudent.DetailUser.Username = req.Nisn
			s.repo.DB().Save(&sc.DetailStudent.DetailUser)
		}
	}

	return dto.DetailStudentResponse{
		Nisn:   req.Nisn,
		Name:   req.Name,
		Gender: req.Gender,
		Class:  dto.GeneralLabelKeyResponse{},
	}
}

// DeleteStudent deletes a student
func (s *StudentService) DeleteStudent(id uint) {
	existingStudent, err := s.repo.FindByID(id)
	if err != nil || existingStudent.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Student not found"))
	}

	// Delete student class
	s.repo.Delete(id)
	// Delete student
	s.repo.DeleteStudent(id)
	// Delete user
	s.repo.DeleteUser(existingStudent.UserID)
}

// DownloadTemplateUpload generates excel template for student upload
func (s *StudentService) DownloadTemplateUpload(c *gin.Context) {
	f := excelize.NewFile()
	sheet := "Upload_Student"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"NISN", "NAMA", "JENIS_KELAMIN(laki-laki/perempuan)", "CLASS_ID"}
	for i, h := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, h)
	}
	f.SetCellValue(sheet, "A2", "12234234234234")
	f.SetCellValue(sheet, "B2", "Nama Siswa")
	f.SetCellValue(sheet, "C2", "perempuan")
	f.SetCellValue(sheet, "D2", 1)

	refSheet := "Class_Reference"
	index, err := f.NewSheet(refSheet)
	if err != nil {
		response.InternalError(c, "Failed create reference sheet", err)
		return
	}

	f.SetCellValue(refSheet, "A1", "classId")
	f.SetCellValue(refSheet, "B1", "className")

	var classes []school.Class
	s.repo.DB().Find(&classes)

	for i, class := range classes {
		f.SetCellValue(refSheet, fmt.Sprintf("A%d", i+2), class.ID)
		f.SetCellValue(refSheet, fmt.Sprintf("B%d", i+2), class.ClassName)
	}

	f.SetActiveSheet(index)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=template_student.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Expires", "0")

	if err := f.Write(c.Writer); err != nil {
		response.InternalError(c, "Failed Generate file template", err)
		return
	}
}

// UploadTemplate handles student upload from excel
func (s *StudentService) UploadTemplate(c *gin.Context) {
	claims := jwt.GetDataClaims(c)

	file, err := c.FormFile("file")
	if err != nil {
		response.InternalError(c, "Failed Upload Student", err)
		return
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".xlsx" {
		response.InternalError(c, "Format file must be .xlsx", nil)
		return
	}

	src, err := file.Open()
	if err != nil {
		response.InternalError(c, "Failed open file", err)
		return
	}
	defer src.Close()

	rows, err := sharedHelper.ReadAndValidateExcel(src)
	if err != nil {
		response.InternalError(c, "Failed Upload Student", err)
		return
	}

	tx := s.repo.DB().Begin()
	for _, row := range rows {
		nisn := row[0]
		nama := row[1]
		jk := row[2]
		classIdStr := row[3]

		classIdInt, _ := strconv.Atoi(classIdStr)

		var classCheck school.Class
		s.repo.DB().Where("id = ? or class_name = ?", classIdInt, classIdStr).First(&classCheck)
		if classCheck.ID == 0 {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Class with ID '"+classIdStr+"' does not exist"))
		}

		role := s.userService.GetRole("STUDENT")
		if role == nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Undefined role 'STUDENT'"))
		}

		// Check existing
		existing, _ := s.repo.FindByNISN(nisn)
		if existing != nil && existing.ID != 0 {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User with NISN '%s' already exists", nisn)))
		}

		// Create user
		newUser := user.User{
			Username:   nisn,
			Role:       role.ID,
			Status:     1,
			SchoolCode: claims.SchoolCode,
		}
		if err := tx.Create(&newUser).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save user"))
		}

		// Create student
		newStudent := student.Student{
			UserId: newUser.ID,
			Name:   nama,
			Nisn:   nisn,
			Gender: jk,
		}
		if err := tx.Create(&newStudent).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save student"))
		}

		// Create student class
		studentClass := student.StudentClass{
			StudentId: newStudent.ID,
			ClassId:   classCheck.ID,
		}
		if err := tx.Create(&studentClass).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save student class"))
		}
	}

	tx.Commit()
	response.OK(c, "Success Upload Template", nil)
}
