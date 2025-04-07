package exam_response

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
)

type ExamListResponse struct {
	ID            uint               `json:"id"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Subject       curriculum.Subject `json:"subject"`
	Member        string             `json:"member"`
	TypeExam      school.TypeExam    `json:"type_exam"`
	TotalQuestion int                `json:"total_question"`
	Duration      int                `json:"duration"`
	TotalScore    int                `json:"total_score"`
}

type ExamDetailResponse struct {
	*school.Exam
	TotalQuestion int `json:"total_question"`
	TotalScore    int `json:"total_score"`
}

type DetailExamQuestionResponse struct {
	ExamCode   string `json:"exam_code"`
	QuestionId string `json:"question_id"`
	Question   string `json:"question"`
	OptionA    string `json:"option_a"`
	OptionB    string `json:"option_b"`
	OptionC    string `json:"option_c"`
	OptionD    string `json:"option_d"`
	OptionE    string `json:"option_e"`
	Answer     string `json:"answer"`
	Score      int    `json:"score"`
}

type MasterBankQuestionResponse struct {
	school.MasterBankQuestion
	TotalQuestion int `json:"total_question"`
}
