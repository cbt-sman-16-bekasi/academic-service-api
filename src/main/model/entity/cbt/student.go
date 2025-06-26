package cbt

import (
	"gorm.io/gorm"
	"time"
)

const TableNameStudentAnswers = "cbt_service.student_answers"
const TableNameStudentHistoryTaken = "cbt_service.student_history_taken"
const TableNameHistoryResetSession = "cbt_service.history_reset_session"
const TableNameHistoryChangeScore = "cbt_service.history_change_score"

type StudentAnswers struct {
	gorm.Model
	ExamCode string `json:"exam_code"`
	//DetailExam    school.Exam        `json:"detail_exam" gorm:"foreignKey:ExamCode;references:ExamCode"`
	SessionId string `json:"session_id"`
	//DetailSession school.ExamSession `json:"detail_session" gorm:"foreignKey:SessionId;references:SessionId"`
	StudentId uint `json:"student_id"`
	//DetailStudent student.Student    `json:"detail_student" gorm:"foreignKey:StudentId;references:ID"`
	QuestionId string `json:"question_id"`
	AnswerId   string `json:"answer_id"`
	Score      int    `json:"score"`
}

func (s StudentAnswers) TableName() string {
	return TableNameStudentAnswers
}

type StudentHistoryTaken struct {
	gorm.Model
	ExamCode string `json:"exam_code"`
	//DetailExam    school.Exam        `json:"detail_exam" gorm:"foreignKey:ExamCode;references:ExamCode"`
	SessionId string `json:"session_id"`
	//DetailSession school.ExamSession `json:"detail_session" gorm:"foreignKey:SessionId;references:SessionId"`
	StudentId uint `json:"student_id"`
	//DetailStudent student.Student    `json:"detail_student" gorm:"foreignKey:StudentId;references:ID"`
	StartAt             *time.Time `json:"start_at"`
	EndAt               *time.Time `json:"end_at"`
	Score               float64    `json:"score"`
	TotalCorrect        int        `json:"total_correct"`
	TotalWrong          int        `json:"total_wrong"`
	Status              string     `json:"status"`
	RemainingTime       int        `json:"remaining_time"`
	IsFinished          bool       `json:"is_finished"`
	IsForced            bool       `json:"is_forced"`
	IsTimeOver          bool       `json:"is_time_over"`
	IsCheating          bool       `json:"is_cheating"`
	NeedCorrection      bool       `json:"need_correction" gorm:"DEFAULT:true"`
	LastScore           float64    `json:"lastScore" gorm:"column:last_score"`
	LastCorrectionScore *time.Time `json:"lastCorrectionScore" gorm:"column:last_correction_score"`
	LastCorrectionBy    string     `json:"lastCorrectionBy" gorm:"column:last_correction_by"`
	LastResetSession    *time.Time `json:"lastResetSession" gorm:"column:last_reset_session"`
	LastResetBy         string     `json:"lastResetBy" gorm:"column:last_reset_by"`
	LastResetReason     string     `json:"lastResetReason" gorm:"column:last_reset_reason"`
}

func (s StudentHistoryTaken) TableName() string {
	return TableNameStudentHistoryTaken
}

type HistoryResetSession struct {
	gorm.Model
	ExamCode string `json:"exam_code"`
	//DetailExam    school.Exam        `json:"detail_exam" gorm:"foreignKey:ExamCode;references:ExamCode"`
	SessionId string `json:"session_id"`
	//DetailSession school.ExamSession `json:"detail_session" gorm:"foreignKey:SessionId;references:SessionId"`
	StudentId uint `json:"student_id"`
	//DetailStudent student.Student    `json:"detail_student" gorm:"foreignKey:StudentId;references:ID"`
	Reason  string `json:"reason"`
	ResetBy string `json:"reset_by"`
}

func (s HistoryResetSession) TableName() string {
	return TableNameHistoryResetSession
}

type HistoryChangeScoreSession struct {
	gorm.Model
	ExamCode string `json:"exam_code"`
	//DetailExam    school.Exam        `json:"detail_exam" gorm:"foreignKey:ExamCode;references:ExamCode"`
	SessionId string `json:"session_id"`
	//DetailSession school.ExamSession `json:"detail_session" gorm:"foreignKey:SessionId;references:SessionId"`
	StudentId uint `json:"student_id"`
	//DetailStudent student.Student    `json:"detail_student" gorm:"foreignKey:StudentId;references:ID"`
	Reason    string  `json:"reason"`
	LastScore float64 `json:"last_score"`
	NewScore  float64 `json:"new_score"`
	ChangeBy  string  `json:"change_by"`
}

func (s HistoryChangeScoreSession) TableName() string {
	return TableNameHistoryChangeScore
}
