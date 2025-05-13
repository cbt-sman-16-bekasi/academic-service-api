package exam_service

import (
	"encoding/json"
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/bucket"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/parsedocx"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/exam_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/exam_repository"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/logger"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

type ExamService struct {
	examRepository *exam_repository.ExamRepository
}

func NewExamService() *ExamService {
	return &ExamService{
		examRepository: exam_repository.NewExamRepository(),
	}
}

func (e *ExamService) GetAllExam(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)
	if claims.Role != "ADMIN" {
		filter := map[string]interface{}{}
		filter["created_by"] = jwt.GetID(claims.Username)

		request.Filter = &filter
	}
	paging := database.NewPagination[map[string]interface{}]().
		SetModal([]school.Exam{}).
		SetPreloads(
			"DetailSubject",
			"ExamMember",
			"ExamMember.DetailClass",
			"ExamMember.DetailClass.DetailClassCode",
			"DetailTypeExam",
			"DetailTypeExam.DetailRole").
		SetRequest(&request).
		FindAllPaging()

	jsByte, _ := json.Marshal(paging.Records)
	var records []school.Exam
	_ = json.Unmarshal(jsByte, &records)

	var newResponse []exam_response.ExamListResponse
	for _, record := range records {
		var member string
		for i, examMember := range record.ExamMember {
			if i > 0 {
				member += ", "
			}
			member += examMember.DetailClass.ClassName
		}

		var totalQuestion int64
		_ = e.examRepository.Database.Where("exam_code = ?", record.Code).Model(&school.ExamQuestion{}).Count(&totalQuestion)
		newResponse = append(newResponse, exam_response.ExamListResponse{
			ID:            record.ID,
			Code:          record.Code,
			Name:          record.Name,
			Subject:       record.DetailSubject,
			Member:        member,
			TypeExam:      record.DetailTypeExam,
			TotalQuestion: int(totalQuestion),
			Duration:      record.Duration,
			TotalScore:    record.ScoreQuestion,
		})
	}

	paging.Records = newResponse
	return paging
}

func (e *ExamService) GetDetailExam(id uint) exam_response.ExamDetailResponse {
	data := e.examRepository.FindById(id)
	if data.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam with id %d not found", id)))
	}

	totalScore := len(data.ExamQuestion) * data.ScoreQuestion

	return exam_response.ExamDetailResponse{
		Exam:          data,
		TotalQuestion: len(data.ExamQuestion),
		TotalScore:    totalScore,
	}
}

func (e *ExamService) CreateNewExam(c *gin.Context, request exam_request.ModifyExamRequest) *school.Exam {
	data := &school.Exam{
		Code:           "EXAM-" + helper.RandomString(10),
		Name:           request.Name,
		Description:    request.Description,
		SubjectCode:    request.SubjectCode,
		TypeExam:       request.TypeExamId,
		RandomQuestion: request.RandomQuestion,
		RandomAnswer:   request.RandomAnswer,
		ShowResult:     request.ShowResult,
		Duration:       request.Duration,
		TypeQuestion:   request.TypeQuestion,
		ScoreQuestion:  request.Score,
	}
	claims := jwt.GetDataClaims(c)
	data.CreatedBy = uint(jwt.GetID(claims.Username))

	e.examRepository.Database.Create(data)

	for _, member := range request.ClassId {
		examMember := &school.ExamMember{
			ExamCode: data.Code,
			Class:    member,
		}
		e.examRepository.Database.Create(examMember)
	}
	return data
}

func (e *ExamService) UpdateExam(c *gin.Context, id uint, request exam_request.ModifyExamRequest) *school.Exam {
	existing := e.examRepository.Repository.FindById(id)
	if existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam with id %d not found", id)))
	}

	existing.Name = request.Name
	existing.Description = request.Description
	existing.SubjectCode = request.SubjectCode
	existing.TypeExam = request.TypeExamId
	existing.RandomQuestion = request.RandomQuestion
	existing.RandomAnswer = request.RandomAnswer
	existing.ShowResult = request.ShowResult
	existing.Duration = request.Duration
	existing.TypeQuestion = request.TypeQuestion
	existing.ScoreQuestion = request.Score

	claims := jwt.GetDataClaims(c)
	existing.ModifiedBy = uint(jwt.GetID(claims.Username))

	e.examRepository.Database.Save(existing)
	e.examRepository.Database.Where("exam_code = ?", existing.Code).Delete(&school.ExamMember{})

	for _, member := range request.ClassId {
		examMember := &school.ExamMember{
			ExamCode: existing.Code,
			Class:    member,
		}
		e.examRepository.Database.Create(examMember)
	}

	return existing
}

func (e *ExamService) DeleteExam(id uint) {
	_ = e.examRepository.Repository.DeleteById(id)
}

func (e *ExamService) GetAllExamQuestion(examId uint) []school.ExamQuestion {
	existing := e.examRepository.FindById(examId)

	return existing.ExamQuestion
}

func (e *ExamService) GetDetailExamQuestion(id uint) exam_response.DetailExamQuestionResponse {
	existing := e.examRepository.FindByIdQuestion(id)
	if existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam with id %d not found", id)))
	}

	options := existing.QuestionOption
	return exam_response.DetailExamQuestionResponse{
		ExamCode:   existing.ExamCode,
		QuestionId: existing.QuestionId,
		Question:   existing.Question,
		OptionA:    e.getOptionByAnswerId(existing.QuestionId+"_A", options).Option,
		OptionB:    e.getOptionByAnswerId(existing.QuestionId+"_B", options).Option,
		OptionC:    e.getOptionByAnswerId(existing.QuestionId+"_C", options).Option,
		OptionD:    e.getOptionByAnswerId(existing.QuestionId+"_D", options).Option,
		OptionE:    e.getOptionByAnswerId(existing.QuestionId+"_E", options).Option,
		Answer:     strings.Split(existing.AnswerSingle, "_")[1],
		Score:      existing.Score,
	}
}

func (e *ExamService) getOptionByAnswerId(answerId string, options []school.ExamAnswerOption) (option school.ExamAnswerOption) {
	for _, answerOption := range options {
		if answerOption.AnswerId == answerId {
			option = answerOption
			break
		}
	}
	return option
}

func (e *ExamService) CreateExamQuestion(request exam_request.ModifyExamQuestionRequest) exam_request.ModifyExamQuestionRequest {
	exam := e.examRepository.FindByCode(request.ExamCode)
	tx := e.examRepository.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	questionID := "QUESTION-" + helper.RandomString(10)
	dataQuestion := &school.ExamQuestion{
		ExamCode:     request.ExamCode,
		QuestionId:   questionID,
		Question:     request.Question,
		Answer:       questionID + "_" + request.Answer,
		AnswerSingle: request.Answer,
		Score:        exam.ScoreQuestion,
		QuestionFrom: "MANUAL",
		TypeQuestion: exam.TypeQuestion,
	}

	if err := tx.Create(dataQuestion).Error; err != nil {
		tx.Rollback()
		logger.Log.Error().Msg("Failed to create exam question")
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed to create exam question"))
	}
	var options = e.setQuestionOptions(questionID, request)
	if err := tx.Create(&options).Error; err != nil {
		tx.Rollback()
		logger.Log.Error().Msg("Failed to create exam question Options")
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed to create exam question Options"))
	}

	e.insertBankQuestion(tx, exam, questionID, request)

	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed Commit save to database"))
	}
	return request
}

func (e *ExamService) AddQuestionFromBank(request exam_request.AddExamQuestionFromBank) exam_request.AddExamQuestionFromBank {
	var questionBank []school.BankQuestion
	e.examRepository.Database.Where("id IN ?", request.BankExamQuestion).Preload("QuestionOption").Find(&questionBank)

	var questions []school.ExamQuestion
	for _, question := range questionBank {
		qID := helper.RandomString(10)
		qst := school.ExamQuestion{
			Model:          gorm.Model{},
			ExamCode:       request.ExamCode,
			QuestionId:     question.QuestionId + qID,
			BankQuestionId: question.QuestionId,
			Question:       question.Question,
			Answer:         question.Answer,
			AnswerSingle:   question.AnswerSingle,
			TypeQuestion:   question.TypeQuestion,
			Score:          0,
			QuestionFrom:   "BANK",
		}

		if question.TypeQuestion == "PILIHAN_GANDA" {
			var options []school.ExamAnswerOption

			for _, option := range question.QuestionOption {
				options = append(options, school.ExamAnswerOption{
					QuestionId: question.QuestionId,
					AnswerId:   option.AnswerId,
					Option:     option.Option,
				})
			}
			e.examRepository.Database.Create(&qst)
		}
		questions = append(questions, qst)
	}
	e.examRepository.Database.Create(&questions)

	return request
}

func (e *ExamService) setQuestionOptions(questionID string, request exam_request.ModifyExamQuestionRequest) []school.ExamAnswerOption {
	var options []school.ExamAnswerOption
	options = append(options, school.ExamAnswerOption{
		QuestionId: questionID,
		AnswerId:   questionID + "_A",
		Option:     request.OptionA,
	})
	options = append(options, school.ExamAnswerOption{
		QuestionId: questionID,
		AnswerId:   questionID + "_B",
		Option:     request.OptionB,
	})
	options = append(options, school.ExamAnswerOption{
		QuestionId: questionID,
		AnswerId:   questionID + "_C",
		Option:     request.OptionC,
	})
	options = append(options, school.ExamAnswerOption{
		QuestionId: questionID,
		AnswerId:   questionID + "_D",
		Option:     request.OptionD,
	})
	options = append(options, school.ExamAnswerOption{
		QuestionId: questionID,
		AnswerId:   questionID + "_E",
		Option:     request.OptionE,
	})
	return options
}

func (e *ExamService) UpdateExamQuestion(id uint, request exam_request.ModifyExamQuestionRequest) exam_request.ModifyExamQuestionRequest {
	exam := e.examRepository.FindByCode(request.ExamCode)

	question := e.examRepository.FindByIdQuestion(id)
	if question.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam with id %d not found", id)))
	}
	questionId := question.QuestionId

	tx := e.examRepository.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	question.Question = request.Question
	question.Score = exam.ScoreQuestion
	question.Answer = questionId + "_" + request.Answer
	question.AnswerSingle = request.Answer

	if err := tx.Save(&question).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
	}

	if err := tx.Where("question_id = ?", questionId).Delete(&school.ExamAnswerOption{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
	}
	var options = e.setQuestionOptions(question.QuestionId, request)
	e.examRepository.Database.Create(&options)
	if err := tx.Create(&options).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
	}

	e.updateBankQuestion(tx, exam, questionId, request)

	// 7. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		log.Println("Gagal commit saat update exam question:", err)
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
	}

	return request
}

func (e *ExamService) updateBankQuestion(tx *gorm.DB, exam school.Exam, questionId string, request exam_request.ModifyExamQuestionRequest) {
	var classCode []string
	seen := make(map[string]bool)

	for _, member := range exam.ExamMember {
		code := member.DetailClass.ClassCode
		if !seen[code] {
			classCode = append(classCode, code)
			seen[code] = true
		}
	}

	for _, member := range classCode {
		var bankQuestion school.BankQuestion
		if err := tx.Where("question_id = ?", questionId+"_"+member).First(&bankQuestion).Error; err != nil {
			var masterBankQuestion school.MasterBankQuestion
			e.examRepository.Database.Where("subject_code = ? AND class_code = ? AND type_question = ?", exam.SubjectCode, member, exam.TypeQuestion).First(&masterBankQuestion)

			if masterBankQuestion.ID == 0 {
				masterBankQuestion = school.MasterBankQuestion{
					Code:         "BANK_MASTER_" + helper.RandomString(10),
					SubjectCode:  exam.SubjectCode,
					ClassCode:    member,
					TypeQuestion: exam.TypeQuestion,
				}

				if err := tx.Create(&masterBankQuestion).Error; err != nil {
					tx.Rollback()
					panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save bank question"))
				}
			}

			bankQuestion = school.BankQuestion{
				QuestionId:             questionId + "_" + member,
				MasterBankQuestionCode: masterBankQuestion.Code,
				TypeQuestion:           exam.TypeQuestion,
				Question:               request.Question,
				Answer:                 questionId + "_" + member + "_" + request.Answer,
				AnswerSingle:           request.Answer,
				QuestionFrom:           "MANUAL",
			}
			if err = tx.Create(&bankQuestion).Error; err != nil {
				tx.Rollback()
				panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
			}
		} else {
			bankQuestion.Question = request.Question
			bankQuestion.Answer = questionId + "_" + member + "_" + request.Answer
			bankQuestion.AnswerSingle = request.Answer
			if err := tx.Save(&bankQuestion).Error; err != nil {
				tx.Rollback()
				panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
			}
		}

		// 6. Hapus dan update BankAnswerOption
		if err := tx.Where("question_id = ?", questionId+"_"+member).Delete(&school.BankAnswerOption{}).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
		}

		bankOptions := e.setBankQuestionOptions(questionId+"_"+member, request)
		if err := tx.Create(&bankOptions).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "failed to update exam question"))
		}
	}
}

func (e *ExamService) DeleteExamQuestion(id uint) {
	question := e.examRepository.FindByIdQuestion(id)
	if question.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("exam with id %d not found", id)))
	}

	e.examRepository.Database.Where("question_id = ?", question.QuestionId).Delete(&school.ExamAnswerOption{})
	e.examRepository.Database.Delete(question)
}

func (e *ExamService) DownloadTemplateUploadQuestion(examId uint, c *gin.Context) {
	typeQuestion := c.Query("typeQuestion")
	bucketName := "dokumen"
	var objectName string

	// Tentukan nama file berdasarkan tipe soal
	if typeQuestion == "PILIHAN_GANDA" {
		objectName = "template_soal_cbt_PILIHAN_GANDA.docx"
	} else {
		objectName = "template_soal_cbt_ESSAY.docx"
	}

	minioBucket := bucket.NewMinio()
	object, err := minioBucket.RetrieveObject(bucketName, objectName)
	if err != nil {
		panic(err)
	}
	defer object.Close()

	// Set header agar bisa di-download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", "attachment; filename="+objectName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Expires", "0")

	// Stream file ke response
	if _, err := io.Copy(c.Writer, object); err != nil {
		response.ErrorResponse(500, "Failed download file", nil)
		return
	}
}

func (e *ExamService) UploadQuestion(c *gin.Context) {
	idParam := c.Param("examId")
	id, _ := strconv.Atoi(idParam)
	exam := e.examRepository.FindById(uint(id))

	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorResponse(response.ServerError, "Failed Upload Question", err).Json(c)
		return
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".docx" {
		response.ErrorResponse(response.ServerError, "Failed Upload Question, Format file must be .docx", err).Json(c)
		return
	}

	src, err := file.Open()
	if err != nil {
		response.ErrorResponse(response.ServerError, "Failed open file", err).Json(c)
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		response.ErrorResponse(response.ServerError, "Failed read file", err).Json(c)
		return
	}

	result, done := e.uploadQuestion(c, exam.TypeQuestion, err, fileBytes, file)
	if !done {
		return
	}

	for i, row := range result {
		questionID := "QUESTION-" + helper.RandomString(10)
		question := row.Soal
		answer := row.Jawaban

		examQuestion := school.ExamQuestion{
			ExamCode:     exam.Code,
			QuestionId:   questionID,
			Question:     question,
			Answer:       questionID + "_" + answer,
			Score:        exam.ScoreQuestion,
			AnswerSingle: answer,
			TypeQuestion: exam.TypeQuestion,
			QuestionFrom: "IMPORT",
		}

		if err := e.examRepository.Database.Create(&examQuestion).Error; err != nil {
			response.ErrorResponse(response.ServerError, fmt.Sprintf("Gagal simpan data di baris %d", i+2), err).Json(c)
			break
		}

		if exam.TypeQuestion == "PILIHAN_GANDA" {
			var examQuestionOption []school.ExamAnswerOption

			examQuestionOption = append(examQuestionOption, school.ExamAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_A",
				Option:     row.A,
			})
			examQuestionOption = append(examQuestionOption, school.ExamAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_B",
				Option:     row.B,
			})
			examQuestionOption = append(examQuestionOption, school.ExamAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_C",
				Option:     row.C,
			})
			examQuestionOption = append(examQuestionOption, school.ExamAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_D",
				Option:     row.D,
			})
			examQuestionOption = append(examQuestionOption, school.ExamAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_E",
				Option:     row.E,
			})

			if err := e.examRepository.Database.Create(&examQuestionOption).Error; err != nil {
				response.ErrorResponse(response.ServerError, fmt.Sprintf("Gagal simpan data di baris %d", i+2), err).Json(c)
				break
			}
		}

	}

	response.SuccessResponse("Success Upload Question", nil).Json(c)
}

func (e *ExamService) uploadQuestion(c *gin.Context, typeQuestion string, err error, fileBytes []byte, file *multipart.FileHeader) ([]parsedocx.ResultParse, bool) {
	var result []parsedocx.ResultParse
	if typeQuestion == "ESSAY" {
		result, err = parsedocx.ParseDocxEssay(fileBytes, file.Filename)
		if err != nil {
			response.ErrorResponse(response.ServerError, "Failed to parse docx", err).Json(c)
			return nil, false
		}
	} else {
		result, err = parsedocx.ParseDocxPilihanGanda(fileBytes, file.Filename)
		if err != nil {
			response.ErrorResponse(response.ServerError, "Failed to parse docx", err).Json(c)
			return nil, false
		}
	}
	return result, true
}

func (e *ExamService) UploadBankQuestion(c *gin.Context) {
	idParam := c.Param("masterBankId")
	id, _ := strconv.Atoi(idParam)

	var bank school.MasterBankQuestion
	e.examRepository.Database.Where("id", uint(id)).First(&bank)

	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorResponse(response.ServerError, "Failed Upload Question", err).Json(c)
		return
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".docx" {
		response.ErrorResponse(response.ServerError, "Failed Upload Question, Format file must be .docx", err).Json(c)
		return
	}

	src, err := file.Open()
	if err != nil {
		response.ErrorResponse(response.ServerError, "Failed open file", err).Json(c)
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		response.ErrorResponse(response.ServerError, "Failed read file", err).Json(c)
		return
	}

	result, done := e.uploadQuestion(c, bank.TypeQuestion, err, fileBytes, file)
	if !done {
		return
	}

	for i, row := range result {

		questionID := "QUESTION-" + helper.RandomString(10)
		question := row.Soal
		answer := row.Jawaban

		examQuestion := school.BankQuestion{
			MasterBankQuestionCode: bank.Code,
			QuestionId:             questionID,
			Question:               question,
			Answer:                 questionID + "_" + answer,
			AnswerSingle:           answer,
			TypeQuestion:           bank.TypeQuestion,
			QuestionFrom:           "IMPORT",
		}

		if err = e.examRepository.Database.Create(&examQuestion).Error; err != nil {
			response.ErrorResponse(response.ServerError, fmt.Sprintf("Gagal simpan data di baris %d", i+2), err).Json(c)
			break
		}

		if bank.TypeQuestion == "PILIHAN_GANDA" {
			var examQuestionOption []school.BankAnswerOption
			examQuestionOption = append(examQuestionOption, school.BankAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_A",
				Option:     row.A,
			})
			examQuestionOption = append(examQuestionOption, school.BankAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_B",
				Option:     row.B,
			})
			examQuestionOption = append(examQuestionOption, school.BankAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_C",
				Option:     row.C,
			})
			examQuestionOption = append(examQuestionOption, school.BankAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_D",
				Option:     row.D,
			})
			examQuestionOption = append(examQuestionOption, school.BankAnswerOption{
				QuestionId: questionID,
				AnswerId:   questionID + "_E",
				Option:     row.E,
			})

			if err = e.examRepository.Database.Create(&examQuestionOption).Error; err != nil {
				response.ErrorResponse(response.ServerError, fmt.Sprintf("Gagal simpan data di baris %d", i+2), err).Json(c)
				break
			}
		}
	}

	response.SuccessResponse("Success Upload Question", nil).Json(c)
}

func ReadAndValidateExcel(file multipart.File) ([][]string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenReader(strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("First sheet undefined")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("Data Empty")
	}

	//header := rows[0]
	//if len(header) != len(expectedHeaders) {
	//	return nil, fmt.Errorf("Jumlah kolom tidak sesuai. Harus %v", expectedHeaders)
	//}
	//for i, v := range expectedHeaders {
	//	if strings.TrimSpace(header[i]) != v {
	//		return nil, fmt.Errorf("Kolom ke-%d harus '%s', dapat '%s'", i+1, v, header[i])
	//	}
	//}

	return rows[1:], nil
}

func (e *ExamService) insertBankQuestion(tx *gorm.DB, exam school.Exam, questionID string, request exam_request.ModifyExamQuestionRequest) {
	var classCode []string
	seen := make(map[string]bool)

	for _, member := range exam.ExamMember {
		code := member.DetailClass.ClassCode
		if !seen[code] {
			classCode = append(classCode, code)
			seen[code] = true
		}
	}

	for _, member := range classCode {
		var masterBank school.MasterBankQuestion
		e.examRepository.Database.Where("subject_code = ? AND class_code = ?",
			exam.SubjectCode,
			member).First(&masterBank)

		if masterBank.ID == 0 {
			masterBank = school.MasterBankQuestion{
				Code:         "BANK_MASTER_" + helper.RandomString(10),
				SubjectCode:  exam.SubjectCode,
				ClassCode:    member,
				TypeQuestion: exam.TypeQuestion,
			}

			if err := tx.Create(&masterBank).Error; err != nil {
				tx.Rollback()
				panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save bank question"))
			}
		}

		bankQuestion := &school.BankQuestion{
			QuestionId:             questionID + "_" + member,
			MasterBankQuestionCode: masterBank.Code,
			TypeQuestion:           exam.TypeQuestion,
			Question:               request.Question,
			Answer:                 questionID + "_" + member + "_" + request.Answer,
			AnswerSingle:           request.Answer,
			QuestionFrom:           "MANUAL",
		}

		if err := tx.Create(bankQuestion).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save bank question"))
		}

		options := e.setBankQuestionOptions(questionID+"_"+member, request)

		if err := tx.Create(&options).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed save bank question"))
		}
	}
}

func (e *ExamService) setBankQuestionOptions(questionID string, request exam_request.ModifyExamQuestionRequest) []school.BankAnswerOption {
	return []school.BankAnswerOption{
		{QuestionId: questionID, AnswerId: questionID + "_A", Option: request.OptionA},
		{QuestionId: questionID, AnswerId: questionID + "_B", Option: request.OptionB},
		{QuestionId: questionID, AnswerId: questionID + "_C", Option: request.OptionC},
		{QuestionId: questionID, AnswerId: questionID + "_D", Option: request.OptionD},
		{QuestionId: questionID, AnswerId: questionID + "_E", Option: request.OptionE},
	}
}

func (e *ExamService) GetAllBankQuestion(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)
	if claims.Role != "ADMIN" {
		filter := map[string]interface{}{}

		var teacherClassSubject []teacher.TeacherClassSubject
		_ = e.examRepository.Database.Where("teacher_id = ?", jwt.GetID(claims.Username)).Find(&teacherClassSubject)

		var subjects = make([]interface{}, 0)
		for _, subject := range teacherClassSubject {
			subjects = append(subjects, subject.SubjectCode)
		}
		filter["subject_code"] = subjects
		request.Filter = &filter
	}
	paging := database.NewPagination[map[string]interface{}]().
		SetModal([]view.MasterBankQuestionResponse{}).
		SetRequest(&request).
		FindAllPaging()
	return paging
}

func (e *ExamService) GetDetailBankQuestion(id uint) exam_response.MasterBankQuestionResponse {
	var detail school.MasterBankQuestion
	e.examRepository.Database.Where("id = ?", id).
		Preload("DetailSubject").
		Preload("DetailClassCode").
		First(&detail)

	var total int64
	e.examRepository.Database.Where("master_bank_question_code = ?", detail.Code).
		Count(&total)
	res := exam_response.MasterBankQuestionResponse{
		MasterBankQuestion: detail,
		TotalQuestion:      int(total),
	}
	return res
}

func (e *ExamService) GetDetailBankQuestionBySubject(request exam_request.ModifyMasterBankQuestionRequest) exam_response.MasterBankQuestionResponse {
	var detail school.MasterBankQuestion
	e.examRepository.Database.Where("subject_code = ? AND class_code = ? AND type_question = ?", request.SubjectCode, request.ClassCode, request.TypeQuestion).
		First(&detail)

	res := exam_response.MasterBankQuestionResponse{
		MasterBankQuestion: detail,
	}
	return res
}

func (e *ExamService) CreateMasterBankQuestion(request exam_request.ModifyMasterBankQuestionRequest) *school.MasterBankQuestion {
	code := "BANK_" + helper.RandomString(10)
	body := &school.MasterBankQuestion{
		BankName:     request.BankName,
		SubjectCode:  request.SubjectCode,
		ClassCode:    request.ClassCode,
		TypeQuestion: request.TypeQuestion,
		Code:         code,
	}
	e.examRepository.Database.Create(&body)
	return body
}

func (e *ExamService) UpdateMasterBankQuestion(id uint, request exam_request.ModifyMasterBankQuestionRequest) *school.MasterBankQuestion {
	var detail school.MasterBankQuestion
	e.examRepository.Database.Where("id = ?", id).
		First(&detail)

	if detail.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Master Bank Question Not Found"))
	}
	detail.BankName = request.BankName
	detail.SubjectCode = request.SubjectCode
	detail.ClassCode = request.ClassCode
	detail.TypeQuestion = request.TypeQuestion

	e.examRepository.Database.Save(&detail)
	return &detail
}

func (e *ExamService) DeleteMasterBankQuestion(id uint) *school.MasterBankQuestion {
	var detail school.MasterBankQuestion
	e.examRepository.Database.Where("id = ?", id).
		First(&detail)
	if detail.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Master Bank Question Not Found"))
	}
	e.examRepository.Database.Delete(&detail)
	return &detail
}

func (e *ExamService) GetQuestionByMasterCode(masterCode string) []school.BankQuestion {
	var questions []school.BankQuestion
	e.examRepository.Database.Where("master_bank_question_code = ?", masterCode).
		Preload("QuestionOption").
		Order("created_at asc").
		Find(&questions)

	return questions
}

func (e *ExamService) GetBankQuestionById(id uint) *school.BankQuestion {
	var question school.BankQuestion
	e.examRepository.Database.Where("id = ?", id).
		Preload("QuestionOption").
		First(&question)
	return &question
}

func (e *ExamService) CreateBankQuestion(request exam_request.ModifyExamQuestionRequest) *school.BankQuestion {
	var masterBankQuestion school.MasterBankQuestion
	e.examRepository.Database.Where("code = ?", request.ExamCode).First(&masterBankQuestion)

	questionId := "QUESTION_" + helper.RandomString(10)
	body := &school.BankQuestion{
		MasterBankQuestionCode: request.ExamCode,
		QuestionId:             questionId,
		TypeQuestion:           masterBankQuestion.TypeQuestion,
		Question:               request.Question,
		Answer:                 questionId + "_" + request.Answer,
		AnswerSingle:           request.Answer,
		QuestionFrom:           "MANUAL",
	}

	e.examRepository.Database.Create(&body)
	if masterBankQuestion.TypeQuestion == "PILIHAN_GANDA" {
		options := e.setBankQuestionOptions(questionId, request)
		e.examRepository.Database.Create(&options)
	}

	return body
}

func (e *ExamService) UpdateBankQuestion(id uint, request exam_request.ModifyExamQuestionRequest) *school.BankQuestion {
	var masterBankQuestion school.MasterBankQuestion
	e.examRepository.Database.Where("code = ?", request.ExamCode).First(&masterBankQuestion)

	var detail school.BankQuestion
	e.examRepository.Database.Where("id = ?", id).First(&detail)
	if detail.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Master Bank Question Not Found"))
	}

	detail.Question = request.Question
	detail.Answer = detail.QuestionId + "_" + request.Answer
	detail.AnswerSingle = request.Answer
	e.examRepository.Database.Save(&detail)

	if masterBankQuestion.TypeQuestion == "PILIHAN_GANDA" {
		options := e.setBankQuestionOptions(detail.QuestionId, request)
		e.examRepository.Database.Where("question_id = ?", detail.QuestionId).Delete(&school.BankAnswerOption{})
		e.examRepository.Database.Create(&options)
	}

	return &detail
}

func (e *ExamService) DeleteBankQuestion(id uint) *school.BankQuestion {
	var detail school.BankQuestion
	e.examRepository.Database.Where("id = ?", id).
		First(&detail)
	if detail.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Master Bank Question Not Found"))
	}
	e.examRepository.Database.Delete(&detail)
	return &detail
}

func (e *ExamService) GetExamMember(code string) []school.ExamMember {
	var data []school.ExamMember
	e.examRepository.Database.Where("exam_code = ?", code).Preload("DetailClass").Find(&data)

	return data
}
