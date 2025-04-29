package student_service

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/auth_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/student_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/auth_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/student_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/exam_service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/user_service"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	response2 "github.com/yon-module/yon-framework/server/response"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type StudentService struct {
	studentRepo    *school_repository.StudentRepository
	userRepository *school_repository.UserRepository
}

func NewStudentService() *StudentService {
	return &StudentService{
		studentRepo:    school_repository.NewStudentRepository(),
		userRepository: school_repository.NewUserRepository(),
	}
}

func (s *StudentService) DetailStudent(id uint) student_response.DetailStudentResponse {
	detail := s.studentRepo.FindById(id)

	return student_response.DetailStudentResponse{
		Nisn:   detail.NISN,
		Name:   detail.Name,
		Gender: detail.Gender,
		Class: response.GeneralLabelKeyResponse{
			Key:   detail.ClassID,
			Label: detail.ClassName,
		},
	}
}

func (s *StudentService) AllStudent(request pagination.Request[map[string]interface{}]) *database.Paginator {
	paging := database.NewPagination[map[string]interface{}]().
		SetRequest(&request).
		SetModal([]view.VStudent{}).FindAllPaging()
	return paging
}

func (s *StudentService) CreateStudent(request student_request.StudentModifyRequest) student_response.DetailStudentResponse {
	role := s.userRepository.ReadRole("STUDENT")
	if role == nil {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Undefined role 'STUDENT'"))
	}

	// Check existing user
	var existingUser student.Student
	s.studentRepo.Database.Where("nisn = ?", request.Nisn).First(&existingUser)
	if existingUser.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, fmt.Sprintf("User with NISN '%s' already exists", request.Nisn)))
	}

	// Create new user
	userService := user_service.NewUserService()
	resultUser := userService.CreateNewUser(&user.User{
		Username:   request.Nisn,
		Role:       role.ID,
		Status:     1,
		SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
	})

	// Register new student
	stud := student.Student{
		UserId: resultUser.ID,
		Name:   request.Name,
		Nisn:   request.Nisn,
		Gender: request.Gender,
	}

	s.userRepository.Database.Save(&stud)

	// Register new Student Class
	studentClass := student.StudentClass{
		StudentId: stud.ID,
		ClassId:   request.ClassId,
	}
	s.studentRepo.Database.Save(&studentClass)

	return student_response.DetailStudentResponse{
		Nisn:   request.Nisn,
		Name:   request.Name,
		Gender: request.Gender,
		Class: response.GeneralLabelKeyResponse{
			Key:   role.Code,
			Label: role.Name,
		},
	}
}

func (s *StudentService) UpdateStudent(id uint, request student_request.StudentModifyRequest) student_response.DetailStudentResponse {
	existingStudent := s.studentRepo.FindById(id)
	if existingStudent.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Student does not exist"))
	}

	if existingStudent.NISN != request.Nisn {
		var existingUser student.Student
		s.studentRepo.Database.Where("nisn = ?", request.Nisn).First(&existingUser)
		if existingUser.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, fmt.Sprintf("User with NISN '%s' already exists", request.Nisn)))
		}
		existingStudent.NISN = request.Nisn
	}
	existingStudent.Name = request.Name
	existingStudent.Gender = request.Gender

	s.studentRepo.Database.Model(&student.Student{}).Where("id = ?", id).Updates(&existingStudent)

	var existStudentClass student.StudentClass
	s.studentRepo.Database.Where("student_id = ?", existingStudent.ID).
		Preload("DetailStudent.DetailUser").
		Preload("DetailStudent").
		First(&existStudentClass)

	existStudentClass.ClassId = request.ClassId
	s.studentRepo.Database.Model(&student.StudentClass{}).Where("student_id = ?", id).Update("class_id", existStudentClass.ClassId)

	userData := existStudentClass.DetailStudent.DetailUser
	userData.Username = request.Nisn
	s.studentRepo.Database.Save(&userData)

	return student_response.DetailStudentResponse{
		Nisn:   request.Nisn,
		Name:   request.Name,
		Gender: request.Gender,
		Class:  response.GeneralLabelKeyResponse{},
	}
}

func (s *StudentService) DeleteById(id uint) {
	existStudentClass := s.studentRepo.FindById(id)
	if existStudentClass.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Student does not exist"))
	}
	s.studentRepo.Delete(id)
	s.studentRepo.Database.Where("id = ?", existStudentClass.ID).Delete(&student.Student{})
	s.studentRepo.Database.Where("id = ?", existStudentClass.UserID).Delete(&user.User{})
}

func (s *StudentService) LoginByNISN(request auth_request.CBTAuthRequest) auth_response.AuthResponseCBT {
	std := s.studentRepo.FindByNISN(request.Username)
	if std.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Student with NISN does not exist. Please contact your system administrator and try again."))
	}

	studentClass := s.studentRepo.GetStudentClass(std.ID)

	var examActive []view.ExamSessionActiveToday
	s.studentRepo.Database.Where("class = ?", studentClass.ClassId).Find(&examActive)

	if len(examActive) == 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "You don't have a exam with that class."))
	}

	var sessionActives []string
	for _, today := range examActive {
		sessionActives = append(sessionActives, today.SessionID)
	}

	examSession := s.studentRepo.ExamRepo.GetExamSessionActiveNow(sessionActives, std.ID)
	if examSession.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "You don't have a exam session with that class."))
	}

	exam := s.studentRepo.ExamRepo.FindByCode(examSession.ExamCode)
	examQuestionRandom := randomizeExam(exam.ExamQuestion, exam.RandomQuestion, exam.RandomAnswer)
	exam.ExamQuestion = examQuestionRandom

	exp := time.Now().Add(time.Hour * 24).Unix()
	token, err := jwt.GenerateJWT(jwt.Claims{
		Username:   request.Username,
		Role:       std.DetailUser.RoleUser.Code,
		Permission: []string{"create", "update", "delete", "read", "list"},
		SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
	})

	if err != nil {
		panic(exception.NewIntenalServerExceptionStruct(response2.ServerError, err.Error()))
	}

	var existingHistoryTaken cbt.StudentHistoryTaken
	s.studentRepo.Database.Where("session_id = ? AND student_id = ?", examSession.SessionId, std.ID).First(&existingHistoryTaken)
	return auth_response.AuthResponseCBT{
		Token:       token,
		Exp:         exp,
		Exam:        &exam,
		User:        studentClass,
		ExamSession: examSession,
		ExamTaken:   &existingHistoryTaken,
	}
}

func shuffle[T any](arr []T) []T {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func randomizeExam(questions []school.ExamQuestion, questionRandom, answerRandom bool) []school.ExamQuestion {
	if questionRandom {
		questions = shuffle(questions)
	}

	if answerRandom {
		for i := range questions {
			questions[i].QuestionOption = shuffle(questions[i].QuestionOption)
		}
	}

	return questions
}

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
		response2.ErrorResponse(response2.ServerError, "Failed create reference sheet", err)
		return
	}

	// Header untuk sheet referensi
	f.SetCellValue(refSheet, "A1", "classId")
	f.SetCellValue(refSheet, "B1", "className")

	var classess []school.Class
	s.studentRepo.Database.Find(&classess)

	for i, class := range classess {
		f.SetCellValue(refSheet, fmt.Sprintf("A%d", i+2), class.ID)
		f.SetCellValue(refSheet, fmt.Sprintf("B%d", i+2), class.ClassName)
	}

	f.SetActiveSheet(index)

	// Set header agar bisa di-download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=template_student.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Expires", "0")

	// Stream file ke response
	if err := f.Write(c.Writer); err != nil {
		response2.ErrorResponse(response2.ServerError, "Failed Generate file template", err)
		return
	}
}

func (s *StudentService) UploadTemplate(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response2.ErrorResponse(response2.ServerError, "Failed Upload Question", err)
		return
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".xlsx" {
		response2.ErrorResponse(response2.ServerError, "Failed Upload Question, Format file must be .xlsx", err)
		return
	}

	src, err := file.Open()
	if err != nil {
		response2.ErrorResponse(response2.ServerError, "Failed open file", err)
		return
	}
	defer src.Close()

	rows, err := exam_service.ReadAndValidateExcel(src)
	if err != nil {
		response2.ErrorResponse(response2.ServerError, "Failed Upload Question", err)
		return
	}

	tx := s.studentRepo.Database.Begin()
	for _, row := range rows {
		nisn := row[0]
		nama := row[1]
		jk := row[2]
		classId := row[3]

		classIdInt, _ := strconv.Atoi(classId)

		var classCheck school.Class
		s.userRepository.Database.Debug().Where("id = ? or class_name = ?", classIdInt, classId).First(&classCheck)
		if classCheck.ID == 0 {
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Class with ID '"+classId+"' does not exist. Please check again with reference on sheet 'Class_Reference'"))
		}

		request := student_request.StudentModifyRequest{
			Nisn:    nisn,
			Name:    nama,
			Gender:  jk,
			ClassId: classCheck.ID,
		}

		role := s.userRepository.ReadRole("STUDENT")
		if role == nil {
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Undefined role 'STUDENT'"))
		}

		// Check existing user
		var existingUser student.Student
		s.studentRepo.Database.Where("nisn = ?", request.Nisn).First(&existingUser)
		if existingUser.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, fmt.Sprintf("User with NISN '%s' already exists", request.Nisn)))
		}
		// Create new user
		resultUser := user.User{
			Username:   request.Nisn,
			Role:       role.ID,
			Status:     1,
			SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
		}
		isExist := s.userRepository.ReadUser(resultUser.Username)
		if isExist != nil {
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "User already exists"))
		}

		if err := tx.Create(&resultUser).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Failed save data user"))
		}

		// Register new student
		stud := student.Student{
			UserId: resultUser.ID,
			Name:   request.Name,
			Nisn:   request.Nisn,
			Gender: request.Gender,
		}

		if err := tx.Create(&stud).Error; err != nil {
			tx.Rollback()
			panic("Failed upload student " + err.Error())
		}

		// Register new Student Class
		studentClass := student.StudentClass{
			StudentId: stud.ID,
			ClassId:   request.ClassId,
		}

		if err := tx.Create(&studentClass).Error; err != nil {
			tx.Rollback()
			panic("Failed upload student " + err.Error())
		}
	}

	tx.Commit()
	response2.SuccessResponse("Success Upload Template", nil).Json(c)
}
