package exam_service

import (
	"encoding/json"
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/exam_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/exam_repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type ExamSessionService struct {
	examSessionRepository *exam_repository.ExamSessionRepository
	studentRepo           *school_repository.StudentRepository
}

func NewExamSessionService() *ExamSessionService {
	return &ExamSessionService{
		examSessionRepository: exam_repository.NewExamSessionRepository(),
		studentRepo:           school_repository.NewStudentRepository(),
	}
}

func (e *ExamSessionService) GetAllExamSession(request pagination.Request[map[string]interface{}]) *database.Paginator {
	paging := database.NewPagination[map[string]interface{}]().
		SetModal([]school.ExamSession{}).
		SetPreloads(
			"DetailExam",
			"DetailExam.DetailSubject",
			"DetailExam.DetailTypeExam",
			"DetailExam.ExamMember",
			"DetailExam.ExamMember.DetailClass",
		).
		SetRequest(&request).
		FindAllPaging()

	jsByte, _ := json.Marshal(paging.Records)
	var records []school.ExamSession
	_ = json.Unmarshal(jsByte, &records)

	var newResponse []exam_response.ExamSessionListResponse
	for _, record := range records {
		var totalStudent int64
		classIds := make([]uint, 0)

		for _, m := range record.DetailExam.ExamMember {
			classIds = append(classIds, m.Class)
		}

		_ = e.studentRepo.Database.Where("class_id IN ?", classIds).Model(&student.StudentClass{}).Count(&totalStudent)
		newResponse = append(newResponse, exam_response.ExamSessionListResponse{
			ExamSession:  record,
			Exam:         record.DetailExam,
			TotalStudent: int(totalStudent),
		})
	}

	paging.Records = newResponse

	return paging
}

func (e *ExamSessionService) GetDetailExamSession(id uint) exam_response.ExamDetailSessionResponse {
	data := e.examSessionRepository.FindById(id)
	if data.ID == 0 {
		return exam_response.ExamDetailSessionResponse{}
	}

	return exam_response.ExamDetailSessionResponse{
		ExamSession:     data,
		Exam:            data.DetailExam,
		TotalStudent:    0,
		TotalAttendance: 0,
	}
}

func (e *ExamSessionService) CreateExamSession(request exam_request.ModifyExamSessionRequest) exam_request.ModifyExamSessionRequest {
	data := &school.ExamSession{
		SessionId: "SESSION-" + helper.RandomString(10),
		ExamCode:  request.ExamCode,
		Name:      request.Name,
		StartDate: request.StartAt,
		EndDate:   request.EndAt,
	}

	e.examSessionRepository.Database.Create(&data)
	return request
}

func (e *ExamSessionService) UpdateExamSession(id uint, request exam_request.ModifyExamSessionRequest) exam_request.ModifyExamSessionRequest {
	existing := e.examSessionRepository.FindById(id)
	if existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam session not found")))
	}

	existing.ExamCode = request.ExamCode
	existing.Name = request.Name
	existing.StartDate = request.StartAt
	existing.EndDate = request.EndAt
	e.examSessionRepository.Database.Save(&existing)
	return request
}

func (e *ExamSessionService) DeleteExamSession(id uint) {
	e.examSessionRepository.Database.Where("id = ?", id).Delete(&school.ExamSession{})
}

func (e *ExamSessionService) GetAllAttendance(request exam_request.ExamSessionAttendanceRequest) []exam_response.ExamSessionAttendanceResponse {
	return make([]exam_response.ExamSessionAttendanceResponse, 0)
}

func (e *ExamSessionService) GenerateToken(request exam_request.ExamSessionGenerateToken) *school.TokenExamSession {
	data := &school.TokenExamSession{
		Model:            gorm.Model{},
		ExamSession:      request.ExamSessionId,
		StartActiveToken: request.StartAt,
		EndActiveToken:   request.EndAt,
		Token:            strings.ToUpper(helper.RandomString(6)),
	}
	e.examSessionRepository.Database.Create(&data)
	return data
}

func (e *ExamSessionService) GetAllToken(request exam_request.ExamSessionTokenFilterRequest) (res []exam_response.ExamSessionTokenResponse) {
	var data []school.TokenExamSession
	e.examSessionRepository.Database.Preload("DetailExamSession").
		Preload("DetailExamSession.DetailExam").
		Preload("DetailExamSession.DetailExam.DetailSubject").
		Preload("DetailExamSession.DetailExam.DetailTypeExam").
		Order("id desc").
		Find(&data)

	for _, tokenExamSession := range data {
		status := "Active"
		tolerance := 5 * time.Second

		if time.Now().After(tokenExamSession.EndActiveToken.Add(tolerance)) {
			status = "Expired"
		}
		res = append(res, exam_response.ExamSessionTokenResponse{
			TokenExamSession: &tokenExamSession,
			Status:           status,
		})
	}
	return
}

func (e *ExamSessionService) ValidateTokenDo(claims jwt.Claims, request exam_request.ExamSessionStartDoWork) cbt.StudentHistoryTaken {
	var tokenExamSession school.TokenExamSession
	e.examSessionRepository.Database.Where("token = ? AND exam_session = ?", request.Token, request.ExamSessionId).
		Preload("DetailExamSession").
		Preload(clause.Associations).
		First(&tokenExamSession)
	if tokenExamSession.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "token exam session not found"))
	}

	timeNow := time.Now()
	e.validateAgeToken(timeNow, tokenExamSession)
	var examSession = tokenExamSession.DetailExamSession
	if examSession.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "exam session not found"))
	}

	if timeNow.Before(tokenExamSession.StartActiveToken) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "start active token timeout"))
	}

	if timeNow.After(tokenExamSession.EndActiveToken) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "end active token timeout"))
	}

	remainingInMinutes := tokenExamSession.EndActiveToken.Sub(timeNow).Minutes()

	student := e.studentRepo.FindByNISN(claims.Username)
	var existingHistoryTaken cbt.StudentHistoryTaken
	e.examSessionRepository.Database.Where("session_id = ? AND student_id = ?", examSession.SessionId, student.ID).First(&existingHistoryTaken)
	if existingHistoryTaken.ID == 0 {
		existingHistoryTaken = cbt.StudentHistoryTaken{
			ExamCode:      tokenExamSession.DetailExamSession.ExamCode,
			SessionId:     examSession.SessionId,
			StudentId:     student.ID,
			StartAt:       timeNow,
			EndAt:         nil,
			Status:        "STARTED",
			RemainingTime: int(remainingInMinutes),
			IsFinished:    false,
			IsForced:      false,
		}
	} else {
		existingHistoryTaken.RemainingTime = int(remainingInMinutes)
	}

	if remainingInMinutes <= 0 {
		existingHistoryTaken.IsForced = true
		existingHistoryTaken.IsFinished = true
		existingHistoryTaken.EndAt = &timeNow
		existingHistoryTaken.Status = "COMPLETED"
	}

	e.examSessionRepository.Database.Save(&existingHistoryTaken)
	return existingHistoryTaken
}

func (e *ExamSessionService) validateAgeToken(timeNow time.Time, tokenExamSession school.TokenExamSession) {
	if timeNow.After(tokenExamSession.EndActiveToken) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Token already expired"))
	}

	if timeNow.Before(tokenExamSession.StartActiveToken) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Token is not active"))
	}
}

func (e *ExamSessionService) SubmitExamSession(claims jwt.Claims, request exam_request.ExamSessionSubmit) cbt.StudentHistoryTaken {
	if request.Result == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "request result is nil"))
	}
	student := e.studentRepo.FindByNISN(claims.Username)
	var existingHistoryTaken cbt.StudentHistoryTaken
	e.examSessionRepository.Database.Where("session_id = ? AND student_id = ?", request.ExamSessionId, student.ID).First(&existingHistoryTaken)
	if existingHistoryTaken.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "exam session not found"))
	}

	timeNow := time.Now()
	existingHistoryTaken.IsForced = request.IsForced
	existingHistoryTaken.IsFinished = true
	existingHistoryTaken.EndAt = &timeNow
	existingHistoryTaken.Status = "COMPLETED"

	// Correction result
	e.examSessionRepository.Database.Where("session_id = ? AND student_id = ?", request.ExamSessionId, student.ID).
		Delete(&cbt.StudentAnswers{})

	var studentAnswers []cbt.StudentAnswers
	totalScore := 0
	totalCorrect := 0
	for _, submit := range request.Result {
		score := 0
		var question school.ExamQuestion
		e.examSessionRepository.Database.Where("question_id = ?", submit.QuestionId).First(&question)
		if submit.AnswerId == question.Answer {
			score = question.Score
			totalCorrect++
		}
		studentAnswers = append(studentAnswers, cbt.StudentAnswers{
			ExamCode:   existingHistoryTaken.ExamCode,
			SessionId:  existingHistoryTaken.SessionId,
			StudentId:  student.ID,
			QuestionId: submit.QuestionId,
			AnswerId:   submit.AnswerId,
			Score:      score,
		})
		totalScore += score
	}

	e.examSessionRepository.Database.Save(&studentAnswers)

	existingHistoryTaken.Score = totalScore
	existingHistoryTaken.TotalCorrect = totalCorrect
	e.examSessionRepository.Database.Save(&existingHistoryTaken)
	return existingHistoryTaken
}
