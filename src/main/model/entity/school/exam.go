package school

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/core"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
	"time"
)

const (
	TableNameTypeExam             = "school_service.type_exam"
	TableNameExam                 = "school_service.exam"
	TableNameExamMember           = "school_service.exam_member"
	TableNameExamQuestion         = "school_service.exam_question"
	TableNameMasterBankQuestion   = "school_service.master_bank_question"
	TableNameBankQuestion         = "school_service.bank_question"
	TableNameExamAnswerOption     = "school_service.exam_answer_option"
	TableNameExamBankAnswerOption = "school_service.bank_answer_option"
	TableNameExamSession          = "school_service.exam_session"
	TableNameExamSessionToken     = "school_service.exam_session_token"
	TableNameExamSessionMember    = "school_service.exam_session_member"
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
	TypeQuestion   string             `json:"type_question" gorm:"type:varchar(50)"`
	TotalScore     int                `json:"total_score" gorm:"type:int"`
	ExamMember     []ExamMember       `json:"exam_member" gorm:"foreignKey:ExamCode;references:Code"`
	ExamQuestion   []ExamQuestion     `json:"exam_question" gorm:"foreignKey:ExamCode;references:Code"`
	core.AuditUser
}

func (s *Exam) TableName() string {
	return TableNameExam
}

type ExamMember struct {
	gorm.Model
	ExamCode    string `json:"exam_code"`
	DetailExam  Exam   `json:"detail_exam" gorm:"foreignKey:ExamCode;references:Code"`
	Class       uint   `json:"class"`
	DetailClass Class  `json:"detail_class" gorm:"foreignKey:Class;references:ID"`
}

func (s *ExamMember) TableName() string {
	return TableNameExamMember
}

type ExamQuestion struct {
	gorm.Model
	ExamCode       string             `json:"exam_code"`
	QuestionId     string             `gorm:"unique" json:"question_id"`
	BankQuestionId string             `json:"bank_question_id"`
	Question       string             `json:"question" gorm:"type:text"`
	Answer         string             `json:"-"`
	AnswerSingle   string             `json:"answer_single"`
	TypeQuestion   string             `json:"type_question" gorm:"type:varchar(50)"`
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

type MasterBankQuestion struct {
	gorm.Model
	Code            string             `gorm:"unique" json:"code"`
	SubjectCode     string             `gorm:"type:varchar(50)" json:"subject_code"`
	DetailSubject   curriculum.Subject `gorm:"foreignKey:SubjectCode;references:Code" json:"detail_subject"`
	ClassCode       string             `gorm:"type:varchar(50)" json:"class_code"`
	DetailClassCode ClassCode          `gorm:"foreignKey:ClassCode;references:Code" json:"detail_class_code"`
	TypeQuestion    string             `json:"type_question" gorm:"type:varchar(50)"`
	core.AuditUser
}

func (s *MasterBankQuestion) TableName() string {
	return TableNameMasterBankQuestion
}

type BankQuestion struct {
	gorm.Model
	MasterBankQuestionCode   string             `gorm:"type:varchar(50)" json:"master_bank_question_code"`
	DetailMasterBankQuestion MasterBankQuestion `gorm:"foreignKey:MasterBankQuestionCode;references:Code" json:"detail_master_bank_question_code"`
	QuestionId               string             `gorm:"unique" json:"question_id"`
	TypeQuestion             string             `json:"type_question" gorm:"type:varchar(50)"`
	Question                 string             `json:"question" gorm:"type:text"`
	Answer                   string             `json:"-"`
	AnswerSingle             string             `json:"answer_single"`
	QuestionFrom             string             `json:"question_from" gorm:"comment:MANUAL or IMPORT"`
	QuestionOption           []BankAnswerOption `json:"question_option" gorm:"foreignKey:QuestionId;references:QuestionId"`
}

func (s *BankQuestion) TableName() string {
	return TableNameBankQuestion
}

type BankAnswerOption struct {
	gorm.Model
	QuestionId string `json:"question_id"`
	AnswerId   string `json:"answer_id"`
	Option     string `json:"option" gorm:"type:text"`
}

func (s *BankAnswerOption) TableName() string {
	return TableNameExamBankAnswerOption
}

type ExamSession struct {
	gorm.Model
	SessionId         string              `gorm:"unique" json:"session_id"`
	ExamCode          string              `json:"-"`
	DetailExam        Exam                `json:"detail_exam" gorm:"foreignKey:ExamCode;references:Code"`
	Name              string              `json:"name"`
	StartDate         time.Time           `json:"start_date"`
	EndDate           time.Time           `json:"end_date"`
	ExamSessionMember []ExamSessionMember `json:"exam_member" gorm:"foreignKey:SessionId;references:SessionId"`
	core.AuditUser
}

func (s *ExamSession) TableName() string {
	return TableNameExamSession
}

type ExamSessionMember struct {
	gorm.Model
	SessionId   string      `json:"session_id"`
	DetailExam  ExamSession `json:"detail_exam" gorm:"foreignKey:SessionId;references:SessionId"`
	Class       uint        `json:"class"`
	DetailClass Class       `json:"detail_class" gorm:"foreignKey:Class;references:ID"`
}

func (s *ExamSessionMember) TableName() string {
	return TableNameExamSessionMember
}

type TokenExamSession struct {
	gorm.Model
	ExamSession       string      `json:"-"`
	DetailExamSession ExamSession `json:"detail_exam_session" gorm:"foreignKey:ExamSession;references:SessionId"`
	StartActiveToken  time.Time   `json:"start_active_token"`
	EndActiveToken    time.Time   `json:"end_active_token"`
	Token             string      `json:"token"`
	core.AuditUser
}

func (s *TokenExamSession) TableName() string {
	return TableNameExamSessionToken
}
