package dto

import "time"

// ===== Request DTOs =====

// CreateAcademicYearRequest request untuk membuat tahun ajaran baru
type CreateAcademicYearRequest struct {
	YearStart   int        `json:"year_start" binding:"required"` // e.g., 2024
	YearEnd     int        `json:"year_end" binding:"required"`   // e.g., 2025
	Semester    int        `json:"semester" binding:"required"`   // 1 = Ganjil, 2 = Genap
	StartDate   *time.Time `json:"start_date"`                    // Optional: semester start date
	EndDate     *time.Time `json:"end_date"`                      // Optional: semester end date
	Description string     `json:"description"`                   // Optional: description
	IsActive    bool       `json:"is_active"`                     // Set as active semester
}

// UpdateAcademicYearRequest request untuk update tahun ajaran
type UpdateAcademicYearRequest struct {
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description string     `json:"description"`
	IsActive    *bool      `json:"is_active"`
}

// SetActiveRequest request untuk set tahun ajaran aktif
type SetActiveRequest struct {
	SemesterID string `json:"semester_id" binding:"required"` // e.g., "20251"
}

// CreateFromDapodikRequest request untuk membuat dari data Dapodik
type CreateFromDapodikRequest struct {
	SemesterID string `json:"semester_id" binding:"required"` // e.g., "20251"
}

// ===== Response DTOs =====

// AcademicYearResponse response untuk tahun ajaran
type AcademicYearResponse struct {
	ID           uint       `json:"id"`
	SemesterID   string     `json:"semester_id"`   // e.g., "20251"
	Year         string     `json:"year"`          // e.g., "2024/2025"
	YearStart    int        `json:"year_start"`    // e.g., 2024
	YearEnd      int        `json:"year_end"`      // e.g., 2025
	Semester     int        `json:"semester"`      // 1 = Ganjil, 2 = Genap
	SemesterName string     `json:"semester_name"` // "Ganjil" / "Genap"
	FullName     string     `json:"full_name"`     // e.g., "2024/2025 Ganjil"
	IsActive     bool       `json:"is_active"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Description  string     `json:"description,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AcademicYearListResponse response untuk list tahun ajaran
type AcademicYearListResponse struct {
	Items      []AcademicYearResponse `json:"items"`
	Total      int64                  `json:"total"`
	ActiveYear *AcademicYearResponse  `json:"active_year,omitempty"`
}

// StudentClassHistoryResponse response untuk riwayat kelas siswa
type StudentClassHistoryResponse struct {
	ID               uint       `json:"id"`
	StudentID        uint       `json:"student_id"`
	StudentName      string     `json:"student_name"`
	ClassID          uint       `json:"class_id"`
	ClassName        string     `json:"class_name"`
	AcademicYearID   uint       `json:"academic_year_id"`
	SemesterID       string     `json:"semester_id"`
	SemesterFullName string     `json:"semester_full_name"`
	JoinedAt         time.Time  `json:"joined_at"`
	LeftAt           *time.Time `json:"left_at,omitempty"`
	Status           string     `json:"status"`
	Notes            string     `json:"notes,omitempty"`
}

// ClassSummaryByYearResponse response untuk ringkasan kelas per tahun ajaran
type ClassSummaryByYearResponse struct {
	SemesterID    string `json:"semester_id"`
	SemesterName  string `json:"semester_name"`
	TotalClasses  int    `json:"total_classes"`
	TotalStudents int    `json:"total_students"`
}

// ===== Dropdown/Select Options =====

// AcademicYearOption for dropdown selection
type AcademicYearOption struct {
	Value    string `json:"value"` // semester_id
	Label    string `json:"label"` // e.g., "2024/2025 Ganjil"
	IsActive bool   `json:"is_active"`
}
