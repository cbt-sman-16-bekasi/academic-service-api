package exam_request

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"time"
)

type ModifyExamSessionRequest struct {
	Name     string    `json:"name"`
	ExamCode string    `json:"exam_code"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	ClassId  []int     `json:"class_id"`
}

type ExamSessionGenerateToken struct {
	ExamCode      string    `json:"exam_code"`
	ExamSessionId string    `json:"exam_session_id"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
}

type ExamSessionAttendanceRequest struct {
	ExamSessionId string `form:"exam_session_id"`
	ClassId       *uint  `form:"class_id"`
}

type ExamSessionTokenFilterRequest struct {
	ExamSessionId string `form:"exam_session_id"`
}

type ExamSessionStartDoWork struct {
	ExamCode      string `json:"exam_code"`
	ExamSessionId string `json:"exam_session_id"`
	Token         string `json:"token"`
}

type ExamSessionSubmit struct {
	ExamCode      string             `json:"exam_code"`
	ExamSessionId string             `json:"exam_session_id"`
	IsForced      bool               `json:"is_forced"`
	IsTimeOver    bool               `json:"is_time_over"`
	IsCheat       bool               `json:"is_cheat"`
	Result        []ExamResultSubmit `json:"result"`
}

type ExamResultSubmit struct {
	QuestionId string `json:"question_id"`
	AnswerId   string `json:"answer_id"`
}

type ExamSessionReportRequest struct {
	ExamCode  string  `json:"exam_code" form:"exam_code" binding:"required"`
	SessionId *string `json:"session_id" form:"session_id"`
}

type ExamSessionStudentAnswer struct {
	ExamCode     string                    `json:"exam_code" form:"exam_code" binding:"required"`
	SessionId    string                    `json:"session_id" form:"session_id"`
	StudentId    string                    `json:"student_id" form:"student_id"`
	AnswerResult *[]school.ExamEssayResult `json:"answer_result"`
}
