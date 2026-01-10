package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	TableNameAcademicYear        = "school_service.m_academic_year"
	TableNameStudentClassHistory = "school_service.t_student_class_history"
)

// SemesterType represents the type of semester
type SemesterType int

const (
	SemesterGanjil SemesterType = 1 // Odd semester
	SemesterGenap  SemesterType = 2 // Even semester
)

func (s SemesterType) String() string {
	switch s {
	case SemesterGanjil:
		return "Ganjil"
	case SemesterGenap:
		return "Genap"
	default:
		return "Unknown"
	}
}

// AcademicYear represents a semester in an academic year
// Example: 2024/2025 Semester 1 (Ganjil)
type AcademicYear struct {
	gorm.Model
	SchoolCode   string       `gorm:"type:varchar(20);not null;index" json:"school_code"`
	SemesterID   string       `gorm:"type:varchar(10);not null;uniqueIndex:idx_school_semester" json:"semester_id"` // Dapodik format: YYYYS (e.g., "20251")
	Year         string       `gorm:"type:varchar(20);not null" json:"year"`                                        // e.g., "2024/2025"
	YearStart    int          `gorm:"not null" json:"year_start"`                                                   // e.g., 2024
	YearEnd      int          `gorm:"not null" json:"year_end"`                                                     // e.g., 2025
	Semester     SemesterType `gorm:"not null" json:"semester"`                                                     // 1 = Ganjil, 2 = Genap
	SemesterName string       `gorm:"type:varchar(20);not null" json:"semester_name"`                               // "Ganjil" / "Genap"
	IsActive     bool         `gorm:"default:false" json:"is_active"`                                               // Current active semester
	StartDate    *time.Time   `gorm:"null" json:"start_date"`                                                       // Semester start date
	EndDate      *time.Time   `gorm:"null" json:"end_date"`                                                         // Semester end date
	Description  string       `gorm:"type:varchar(255);null" json:"description"`                                    // Optional description
}

func (a *AcademicYear) TableName() string {
	return TableNameAcademicYear
}

// GetFullName returns the full name like "2024/2025 Ganjil"
func (a *AcademicYear) GetFullName() string {
	return a.Year + " " + a.SemesterName
}

// StudentClassHistory tracks student class assignments per semester
// This allows tracking which class a student was in for each semester
type StudentClassHistory struct {
	gorm.Model
	SchoolCode       string     `gorm:"type:varchar(20);not null;index" json:"school_code"`
	StudentID        uint       `gorm:"not null;index" json:"student_id"`
	ClassID          uint       `gorm:"not null;index" json:"class_id"`
	AcademicYearID   uint       `gorm:"not null;index" json:"academic_year_id"`
	SemesterID       string     `gorm:"type:varchar(10);not null;index" json:"semester_id"` // Denormalized for easier queries
	DapodikRombelID  string     `gorm:"type:varchar(100);null" json:"dapodik_rombel_id"`    // rombongan_belajar_id from Dapodik
	DapodikAnggotaID string     `gorm:"type:varchar(100);null" json:"dapodik_anggota_id"`   // anggota_rombel_id from Dapodik
	JoinedAt         time.Time  `gorm:"not null" json:"joined_at"`
	LeftAt           *time.Time `gorm:"null" json:"left_at"`
	Status           string     `gorm:"type:varchar(20);default:'active'" json:"status"` // active, graduated, transferred, dropped
	Notes            string     `gorm:"type:text;null" json:"notes"`
}

func (s *StudentClassHistory) TableName() string {
	return TableNameStudentClassHistory
}

// ParseSemesterID parses a Dapodik semester_id (e.g., "20251") into components
// Returns: yearEnd (2025), semester (1), yearStart (2024), yearString ("2024/2025")
func ParseSemesterID(semesterID string) (yearEnd int, semester SemesterType, yearStart int, yearString string) {
	if len(semesterID) != 5 {
		return 0, 0, 0, ""
	}

	// Parse year (first 4 digits)
	var year int
	for i := 0; i < 4; i++ {
		year = year*10 + int(semesterID[i]-'0')
	}

	// Parse semester (last digit)
	sem := int(semesterID[4] - '0')

	yearEnd = year
	yearStart = year - 1
	semester = SemesterType(sem)
	yearString = formatYear(yearStart, yearEnd)

	return
}

// FormatSemesterID creates a Dapodik format semester_id from year and semester
// Example: FormatSemesterID(2024, 2025, 1) = "20251"
func FormatSemesterID(yearStart, yearEnd int, semester SemesterType) string {
	return formatInt(yearEnd) + formatInt(int(semester))
}

func formatYear(start, end int) string {
	return formatInt(start) + "/" + formatInt(end)
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
