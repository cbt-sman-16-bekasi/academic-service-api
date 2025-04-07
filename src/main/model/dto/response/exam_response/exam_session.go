package exam_response

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"time"
)

type ExamSessionListResponse struct {
	school.ExamSession
	Exam         school.Exam `json:"exam"`
	TotalStudent int         `json:"total_student"`
}

type ExamDetailSessionResponse struct {
	*school.ExamSession
	Exam            school.Exam `json:"exam"`
	TotalStudent    int         `json:"total_student"`
	TotalAttendance int         `json:"total_attendance"`
}

type ExamSessionAttendanceResponse struct {
	Nisn    string    `json:"nisn"`
	Name    string    `json:"name"`
	Class   string    `json:"class"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Score   int       `json:"score"`
	Status  string    `json:"status"`
}

type ExamSessionTokenResponse struct {
	*school.TokenExamSession
	Status string `json:"status"`
}
