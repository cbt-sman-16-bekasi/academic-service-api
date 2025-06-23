package exam_response

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"time"
)

type ExamSessionListResponse struct {
	school.ExamSession
	Exam         school.Exam `json:"exam"`
	TotalStudent int         `json:"total_student"`
	IsActive     bool        `json:"is_active"`
}

type ExamDetailSessionResponse struct {
	*school.ExamSession
	Exam            school.Exam `json:"exam"`
	TotalStudent    int64       `json:"total_student"`
	TotalAttendance int64       `json:"total_attendance"`
	TotalSubmit     int64       `json:"total_submit"`
	TotalCheating   int64       `json:"total_cheating"`
	TotalTimesOver  int64       `json:"total_times_over"`
}

type ExamSessionAttendanceResponse struct {
	Nisn                string     `json:"nisn"`
	StudentId           uint       `json:"student_id"`
	Name                string     `json:"name"`
	Class               string     `json:"class"`
	StartAt             *time.Time `json:"start_at"`
	EndAt               *time.Time `json:"end_at"`
	Score               float64    `json:"score"`
	Status              string     `json:"status"`
	NeedCorrection      bool       `json:"need_correction"`
	LastCorrectionScore string     `json:"lastCorrectionScore"`
	LastCorrectionBy    string     `json:"lastCorrectionBy"`
}

type ExamSessionTokenResponse struct {
	*school.TokenExamSession
	Status string `json:"status"`
}
