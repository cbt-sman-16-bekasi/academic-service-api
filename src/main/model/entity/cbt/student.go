package cbt

import (
	"gorm.io/gorm"
	"time"
)

const TableNameStudentAnswers = "cbt_service.student_answers"
const TableNameStudentHistoryTaken = "cbt_service.student_history_taken"

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
	StartAt       time.Time  `json:"start_at"`
	EndAt         *time.Time `json:"end_at"`
	Score         int        `json:"score"`
	TotalCorrect  int        `json:"total_correct"`
	TotalWrong    int        `json:"total_wrong"`
	Status        string     `json:"status"`
	RemainingTime int        `json:"remaining_time"`
	IsFinished    bool       `json:"is_finished"`
	IsForced      bool       `json:"is_forced"`
}

func (s StudentHistoryTaken) TableName() string {
	return TableNameStudentHistoryTaken
}
