package school

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
	"time"
)

const (
	TableNameTypeExam         = "school_service.type_exam"
	TableNameExam             = "school_service.exam"
	TableNameExamMember       = "school_service.exam_member"
	TableNameExamQuestion     = "school_service.exam_question"
	TableNameExamAnswerOption = "school_service.exam_answer_option"
	TableNameExamSession      = "school_service.exam_session"
	TableNameExamSessionToken = "school_service.exam_session_token"
)

type TypeExam struct {
	gorm.Model
	Code       string    `gorm:"unique" json:"code"`
	Name       string    `gorm:"unique" json:"name"`
	Color      string    `json:"color"`
	Role       string    `json:"-"`
	DetailRole user.Role `json:"detail_role" gorm:"foreignKey:Role;references:Code"`
}

func (s *TypeExam) TableName() string {
	return TableNameTypeExam
}

type Exam struct {
	gorm.Model
	Code           string             `gorm:"unique" json:"code"`
	Name           string             `json:"name"`
	Description    string             `json:"description" gorm:"type:text"`
	SubjectCode    string             `json:"-"`
	DetailSubject  curriculum.Subject `gorm:"foreignKey:SubjectCode;references:Code" json:"subject_code"`
	TypeExam       string             `json:"-"`
	DetailTypeExam TypeExam           `json:"detail_type_exam" gorm:"foreignKey:TypeExam;references:Code"`
	RandomQuestion bool               `json:"random_question"`
	RandomAnswer   bool               `json:"random_answer"`
	ShowResult     bool               `json:"show_result"`
	Duration       int                `json:"duration"`
	ExamMember     []ExamMember       `json:"exam_member" gorm:"foreignKey:ExamCode;references:Code"`
}

func (s *Exam) TableName() string {
	return TableNameExam
}

type ExamMember struct {
	gorm.Model
	ExamCode           string       `json:"-"`
	DetailExam         Exam         `json:"detail_exam" gorm:"foreignKey:ExamCode;references:Code"`
	ClassSubject       uint         `json:"-"`
	DetailClassSubject ClassSubject `json:"detail_class_subject" gorm:"foreignKey:ClassSubject;references:ID"`
}

func (s *ExamMember) TableName() string {
	return TableNameExamMember
}

type ExamQuestion struct {
	gorm.Model
	ExamCode       string             `json:"exam_code"`
	QuestionId     string             `gorm:"unique" json:"question_id"`
	Question       string             `json:"question" gorm:"type:text"`
	Answer         string             `json:"answer"`
	Score          int                `json:"score"`
	QuestionFrom   string             `json:"question_from" gorm:"comment:MANUAL or BANK_QUESTION"`
	QuestionOption []ExamAnswerOption `json:"question_option" gorm:"foreignKey:QuestionId;references:QuestionId"`
}

func (s *ExamQuestion) TableName() string {
	return TableNameExamQuestion
}

type ExamAnswerOption struct {
	gorm.Model
	QuestionId string `json:"question_id"`
	AnswerId   string `json:"answer_id"`
	Option     string `json:"option" gorm:"type:text"`
}

func (s *ExamAnswerOption) TableName() string {
	return TableNameExamAnswerOption
}

type ExamSession struct {
	gorm.Model
	ExamCode   string    `json:"-"`
	DetailExam Exam      `json:"detail_exam" gorm:"foreignKey:ExamCode;references:Code"`
	Name       string    `json:"name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

func (s *ExamSession) TableName() string {
	return TableNameExamSession
}

type TokenExamSession struct {
	gorm.Model
	ExamSession       uint      `json:"-"`
	DetailExamSession string    `json:"detail_exam_session"`
	StartActiveToken  time.Time `json:"start_active_token"`
	EndActiveToken    time.Time `json:"end_active_token"`
	Token             string    `json:"token"`
}

func (s *TokenExamSession) TableName() string {
	return TableNameExamSessionToken
}
