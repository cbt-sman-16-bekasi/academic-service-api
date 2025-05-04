package reporting_service

import "time"

type DataNilai struct {
	NISN      string
	Name      string
	ClassName string
	Gender    string
	Score     float64
}

type DataExamSession struct {
	TypeExam     string
	Subject      string
	ClassName    string
	SessionName  string
	SessionStart time.Time
	SessionEnd   time.Time
	ScoreData    []DataNilai
}
