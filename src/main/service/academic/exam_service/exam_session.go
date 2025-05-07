package exam_service

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/exam_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/redisstore"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/exam_repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/reporting_service"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/logger"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
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

func (e *ExamSessionService) GetAllExamSession(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)
	if claims.Role != "ADMIN" {
		filter := map[string]interface{}{}
		filter["created_by"] = jwt.GetID(claims.Username)

		request.Filter = &filter
	}
	paging := database.NewPagination[map[string]interface{}]().
		SetModal([]view.ExamSessionView{}).
		SetRequest(&request).
		FindAllPaging()

	return paging
}

func (e *ExamSessionService) GetDetailExamSession(id uint) exam_response.ExamDetailSessionResponse {
	data := e.examSessionRepository.FindById(id)
	if data.ID == 0 {
		return exam_response.ExamDetailSessionResponse{}
	}

	var summarySession view.SummaryExamSession
	e.examSessionRepository.Database.Where("session_id = ?", data.SessionId).First(&summarySession)

	return exam_response.ExamDetailSessionResponse{
		ExamSession:     data,
		Exam:            data.DetailExam,
		TotalStudent:    0,
		TotalAttendance: summarySession.TotalLogin,
		TotalSubmit:     summarySession.TotalStudentSubmit,
		TotalCheating:   summarySession.TotalCheating,
		TotalTimesOver:  summarySession.TotalTimeIsOver,
	}
}

func (e *ExamSessionService) CreateExamSession(c *gin.Context, request exam_request.ModifyExamSessionRequest) exam_request.ModifyExamSessionRequest {
	data := &school.ExamSession{
		SessionId: "SESSION-" + helper.RandomString(10),
		ExamCode:  request.ExamCode,
		Name:      request.Name,
		StartDate: request.StartAt,
		EndDate:   request.EndAt,
	}

	claims := jwt.GetDataClaims(c)
	data.CreatedBy = uint(jwt.GetID(claims.Username))

	e.examSessionRepository.Database.Create(&data)

	var dataExamSessionMember []school.ExamSessionMember
	for _, classId := range request.ClassId {
		dataExamSessionMember = append(dataExamSessionMember, school.ExamSessionMember{
			SessionId: data.SessionId,
			Class:     uint(classId),
		})
	}

	e.examSessionRepository.Database.Create(&dataExamSessionMember)
	return request
}

func (e *ExamSessionService) UpdateExamSession(c *gin.Context, id uint, request exam_request.ModifyExamSessionRequest) exam_request.ModifyExamSessionRequest {
	existing := e.examSessionRepository.FindById(id)
	if existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam session not found")))
	}

	existing.ExamCode = request.ExamCode
	existing.Name = request.Name
	existing.StartDate = request.StartAt
	existing.EndDate = request.EndAt

	claims := jwt.GetDataClaims(c)
	existing.ModifiedBy = uint(jwt.GetID(claims.Username))
	e.examSessionRepository.Database.Save(&existing)

	e.examSessionRepository.Database.Where("session_id = ?", existing.SessionId).Delete(&school.ExamSessionMember{})
	var dataExamSessionMember []school.ExamSessionMember
	for _, classId := range request.ClassId {
		dataExamSessionMember = append(dataExamSessionMember, school.ExamSessionMember{
			SessionId: existing.SessionId,
			Class:     uint(classId),
		})
	}

	e.examSessionRepository.Database.Create(&dataExamSessionMember)
	return request
}

func (e *ExamSessionService) DeleteExamSession(id uint) {
	e.examSessionRepository.Database.Where("id = ?", id).Delete(&school.ExamSession{})
}

func (e *ExamSessionService) GetAllAttendance(request exam_request.ExamSessionAttendanceRequest) []exam_response.ExamSessionAttendanceResponse {
	var studentClasses []student.StudentClass
	e.examSessionRepository.Database.Where("class_id = ?", request.ClassId).
		Preload("DetailStudent").
		Preload("DetailClass").
		Find(&studentClasses)

	var responses []exam_response.ExamSessionAttendanceResponse

	for _, class := range studentClasses {
		var studentAttendance cbt.StudentHistoryTaken
		e.examSessionRepository.Database.Where("session_id = ? AND student_id = ?", request.ExamSessionId, class.DetailStudent.ID).First(&studentAttendance)
		status := studentAttendance.Status

		if status == "STARTED" {
			status = "Aktif"
		}
		if status == "COMPLETED" {
			status = "Selesai"
		}

		if studentAttendance.IsForced {
			status = "Dikumpulkan Oleh Sistem"
			if studentAttendance.IsTimeOver {
				status = "Mengumpulkan Pada Waktu Habis"
			}
			if studentAttendance.IsCheating {
				status = "Terindikasi Kecurangan"
			}
		}
		responses = append(responses, exam_response.ExamSessionAttendanceResponse{
			Nisn:    class.DetailStudent.Nisn,
			Name:    strings.ToUpper(class.DetailStudent.Name),
			Class:   class.DetailClass.ClassName,
			StartAt: &studentAttendance.StartAt,
			EndAt:   studentAttendance.EndAt,
			Score:   studentAttendance.Score,
			Status:  status,
		})
	}

	return responses
}

func (e *ExamSessionService) GenerateToken(c *gin.Context, request exam_request.ExamSessionGenerateToken) *school.TokenExamSession {
	data := &school.TokenExamSession{
		Model:            gorm.Model{},
		ExamSession:      request.ExamSessionId,
		StartActiveToken: request.StartAt,
		EndActiveToken:   request.EndAt,
		Token:            strings.ToUpper(helper.RandomString(6)),
	}

	claims := jwt.GetDataClaims(c)
	data.CreatedBy = uint(jwt.GetID(claims.Username))
	e.examSessionRepository.Database.Create(&data)
	return data
}

func (e *ExamSessionService) GetAllToken(c *gin.Context, request exam_request.ExamSessionTokenFilterRequest) (res []exam_response.ExamSessionTokenResponse) {
	var data []school.TokenExamSession
	q := e.examSessionRepository.Database.Preload("DetailExamSession").
		Preload("DetailExamSession.DetailExam").
		Preload("DetailExamSession.DetailExam.DetailSubject").
		Preload("DetailExamSession.DetailExam.DetailTypeExam").
		Where("end_active_token >= ?", time.Now())

	claims := jwt.GetDataClaims(c)
	if claims.Role != "ADMIN" {
		q = q.Where("created_by", jwt.GetID(claims.Username))
	}

	q = q.Order("id desc").
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
		Preload("DetailExamSession.DetailExam").
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

	remainingInMinutes := examSession.EndDate.Sub(timeNow).Minutes()
	if int(remainingInMinutes) > examSession.DetailExam.Duration {
		remainingInMinutes = float64(examSession.DetailExam.Duration)
	}

	studentData := e.studentRepo.FindByNISN(claims.Username)
	var existingHistoryTaken cbt.StudentHistoryTaken
	e.examSessionRepository.Database.Where("session_id = ? AND student_id = ?", examSession.SessionId, studentData.ID).First(&existingHistoryTaken)
	if existingHistoryTaken.ID == 0 {
		existingHistoryTaken = cbt.StudentHistoryTaken{
			ExamCode:      tokenExamSession.DetailExamSession.ExamCode,
			SessionId:     examSession.SessionId,
			StudentId:     studentData.ID,
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
	existingHistoryTaken.IsTimeOver = request.IsTimeOver
	existingHistoryTaken.IsForced = request.IsForced
	existingHistoryTaken.IsFinished = true
	existingHistoryTaken.EndAt = &timeNow
	existingHistoryTaken.Status = "COMPLETED"

	// Correction result
	e.examSessionRepository.Database.Where("session_id = ? AND student_id = ?", request.ExamSessionId, student.ID).
		Delete(&cbt.StudentAnswers{})

	var studentAnswers []cbt.StudentAnswers
	totalScore := 0
	totalAllScore := 0
	totalCorrect := 0

	var questions []school.ExamQuestion
	err := redisstore.GetJSON(request.ExamCode, &questions)
	if err != nil {
		e.examSessionRepository.Database.Where("exam_code", request.ExamCode).Preload("QuestionOption").Find(&questions)
		_ = redisstore.SetJSON(request.ExamCode, &questions, time.Hour*24)
	}

	if questions == nil {
		e.examSessionRepository.Database.Where("exam_code", request.ExamCode).Preload("QuestionOption").Find(&questions)
		_ = redisstore.SetJSON(request.ExamCode, &questions, time.Hour*24)
	}
	for _, question := range questions {
		totalAllScore += question.Score
	}

	for _, submit := range request.Result {
		score := 0

		var question school.ExamQuestion
		for _, examQuestion := range questions {
			if examQuestion.QuestionId == submit.QuestionId {
				question = examQuestion
				break
			}
		}

		if submit.AnswerId == question.Answer && question.TypeQuestion == "PILIHAN_GANDA" {
			score = question.Score
			totalCorrect++
		}

		if question.TypeQuestion == "ESSAY" {
			essayHelper := helper.NewCosineSimilarity(question.AnswerSingle, submit.AnswerId, question.Score)
			score = essayHelper.EvaluateScoreEssay()
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

	existingHistoryTaken.Score = (totalScore / totalAllScore) * 100
	existingHistoryTaken.TotalCorrect = totalCorrect
	e.examSessionRepository.Database.Save(&existingHistoryTaken)
	return existingHistoryTaken
}

func (e *ExamSessionService) ExportExamSessionAttendanceToExcel(c *gin.Context, responses []exam_response.ExamSessionAttendanceResponse, request exam_request.ExamSessionAttendanceRequest) {
	var examSession school.ExamSession
	_ = e.studentRepo.Database.Where("session_id = ?", request.ExamSessionId).
		Preload("DetailExam").
		Preload("DetailExam.DetailSubject").
		Preload("DetailExam.DetailTypeExam").
		First(&examSession)

	f := excelize.NewFile()
	sheet := "Attendance"
	index, _ := f.NewSheet(sheet)

	// Header
	headers := []string{"No", "NISN", "Name", "Class", "Start At", "End At", "Score", "Status"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data
	for i, r := range responses {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), i+1)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), r.Nisn)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), r.Name)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), r.Class)
		// Start At
		if r.StartAt != nil {
			f.SetCellValue(sheet, "E"+strconv.Itoa(row), r.StartAt.Format("2006-01-02 15:04:05"))
		} else {
			f.SetCellValue(sheet, "E"+strconv.Itoa(row), "")
		}

		// End At
		if r.EndAt != nil {
			f.SetCellValue(sheet, "F"+strconv.Itoa(row), r.EndAt.Format("2006-01-02 15:04:05"))
		} else {
			f.SetCellValue(sheet, "F"+strconv.Itoa(row), "")
		}

		f.SetCellValue(sheet, "G"+strconv.Itoa(row), r.Score)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), r.Status)
	}

	f.SetActiveSheet(index)

	fileName := fmt.Sprintf(
		"Data Peserta_%s_%s_%s.xlsx",
		examSession.DetailExam.DetailTypeExam.Code,
		examSession.DetailExam.DetailSubject.Subject,
		examSession.EndDate.Format("20060102"),
	)
	// Stream Excel ke response
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("File-Name", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Expires", "0")

	if err := f.Write(c.Writer); err != nil {
		panic(exception.NewIntenalServerExceptionStruct(
			response.ServerError, "Failed generate report"),
		)
	}
}

func (e *ExamSessionService) ExamSessionMember(sessionId string) []school.ExamSessionMember {
	var examSessionMembers []school.ExamSessionMember
	e.examSessionRepository.Database.Where("session_id = ?", sessionId).Preload("DetailClass").Find(&examSessionMembers)
	return examSessionMembers
}

func (e *ExamSessionService) GenerateReportSession() {
	var schoolData school.School
	e.examSessionRepository.Database.Where("school_code=?", "db74a42e-23a7-4cd2-bbe5-49cf79f86453").First(&schoolData)

	var examSessionReport []view.ExamSessionReadyReport
	e.examSessionRepository.Database.Find(&examSessionReport)

	for _, session := range examSessionReport {
		logger.Log.Info().Msgf("ExamSession Report %v", session)
		report := reporting_service.NewReport(schoolData)

		var reportSession []reporting_service.DataExamSession
		classIds, _ := StringToUintSlice(session.KelasID)
		classNames := strings.Split(session.Kelas, ",")
		for i, classId := range classIds {
			dataScore := e.GetAllAttendance(exam_request.ExamSessionAttendanceRequest{
				ExamSessionId: session.SessionID,
				ClassId:       &classId,
			})

			var reportScore []reporting_service.DataNilai
			for _, data := range dataScore {
				reportScore = append(reportScore, reporting_service.DataNilai{
					NISN:      data.Nisn,
					Name:      data.Name,
					ClassName: data.Class,
					Gender:    "-",
					Score:     float64(data.Score),
				})
			}
			dataSession := reporting_service.DataExamSession{
				TypeExam:     session.TypeExam,
				Subject:      session.Subject,
				ClassName:    classNames[i],
				SessionName:  session.SessionName,
				SessionStart: session.StartDate,
				SessionEnd:   session.EndDate,
				ScoreData:    reportScore,
			}
			reportSession = append(reportSession, dataSession)
		}
		report.SetData(reportSession)

		err := report.Generate()
		if err != nil {
			e.examSessionRepository.Database.Where("session_id = ?", session.SessionID).Update("error_report", err.Error())
			logger.Log.Error().Msgf("Generate reportSession Error: %s", err.Error())
			continue
		}

		resUrlReport, isSuccess := report.GetResult()
		logger.Log.Info().Str("URL", fmt.Sprintf("%v", resUrlReport)).Str("Status", fmt.Sprintf("%v", isSuccess)).Msg("Generate reportSession Result")
		if !isSuccess {
			errorReport := ErrorsToString(report.GetError())
			e.examSessionRepository.Database.Where("session_id = ?", session.SessionID).Model(&school.ExamSession{}).Update("error_report", errorReport)
			logger.Log.Error().Msgf("Generate reportSession Error: %s", errorReport)
			continue
		}

		e.examSessionRepository.Database.Debug().Where("session_id = ?", session.SessionID).Model(&school.ExamSession{}).Updates(map[string]interface{}{
			"report_url":    resUrlReport,
			"status_report": "READY",
		})
	}

}

func (e *ExamSessionService) GetAllReport(request exam_request.ExamSessionReportRequest) []view.ExamSessionReportScoreView {
	var data []view.ExamSessionReportScoreView
	q := e.examSessionRepository.Database.Where("exam_code=?", request.ExamCode)
	if request.SessionId != nil {
		q = q.Where("session_id=?", *request.SessionId)
	}
	q.Find(&data)
	return data
}

func (e *ExamSessionService) GetAnswerStudent(request exam_request.ExamSessionStudentAnswer) []school.ExamEssayResult {
	var data []school.ExamEssayResult
	e.examSessionRepository.Database.Raw(`select q.question_id, sa.session_id, q.question, q.answer_single, sa.answer_id as answer_user, sa.id as answerID, sa.score from school_service.exam_question q
left join school_service.exam e on e.code = q.exam_code
LEFT JOIN cbt_service.student_answers sa ON sa.question_id = q.question_id AND sa.student_id = ? AND sa.session_id = ?
WHERE e.type_question = 'ESSAY' and q.exam_code = ?`, request.StudentId, request.SessionId, request.ExamCode).Scan(&data)
	return data
}

func StringToUintSlice(s string) ([]uint, error) {
	parts := strings.Split(s, ",")
	result := make([]uint, 0, len(parts))

	for _, part := range parts {
		p := strings.TrimSpace(part)
		num, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(num))
	}

	return result, nil
}

func ErrorsToString(errs []error) string {
	strs := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			strs = append(strs, err.Error())
		}
	}
	return strings.Join(strs, ", ")
}
