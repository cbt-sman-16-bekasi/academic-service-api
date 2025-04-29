package exam_request

type ModifyExamRequest struct {
	Name           string `json:"name"`
	SubjectCode    string `json:"subject_code"`
	ClassId        []uint `json:"class_id"`
	TypeExamId     string `json:"type_exam_id"`
	Description    string `json:"description"`
	TypeQuestion   string `json:"type_question"`
	RandomQuestion bool   `json:"random_question"`
	RandomAnswer   bool   `json:"random_answer"`
	ShowResult     bool   `json:"show_result"`
	Duration       int    `json:"duration"`
	Score          int    `json:"score"`
}

type ModifyExamQuestionRequest struct {
	ExamCode string `json:"exam_code"`
	Question string `json:"question"`
	OptionA  string `json:"option_a"`
	OptionB  string `json:"option_b"`
	OptionC  string `json:"option_c"`
	OptionD  string `json:"option_d"`
	OptionE  string `json:"option_e"`
	Answer   string `json:"answer"`
	Score    int    `json:"score"`
}

type ModifyMasterBankQuestionRequest struct {
	BankName     string `json:"bank_name"`
	SubjectCode  string `json:"subject_code"`
	ClassCode    string `json:"class_code"`
	TypeQuestion string `json:"type_question"`
}
