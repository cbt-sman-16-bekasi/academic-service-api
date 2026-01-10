package service

import (
	"fmt"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/entity"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// AcademicYearService handles business logic for academic year
type AcademicYearService struct {
	repo *repository.AcademicYearRepository
}

// NewAcademicYearService creates a new service instance
func NewAcademicYearService(repo *repository.AcademicYearRepository) *AcademicYearService {
	return &AcademicYearService{repo: repo}
}

// ===== Academic Year Operations =====

// GetAll gets all academic years for the current school
func (s *AcademicYearService) GetAll(c *gin.Context) *dto.AcademicYearListResponse {
	claims := jwt.GetDataClaims(c)

	years, err := s.repo.FindBySchoolCode(claims.SchoolCode)
	if err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal mengambil data tahun ajaran"))
	}

	items := make([]dto.AcademicYearResponse, len(years))
	var activeYear *dto.AcademicYearResponse

	for i, year := range years {
		resp := s.mapToResponse(&year)
		items[i] = *resp
		if year.IsActive {
			activeYear = resp
		}
	}

	return &dto.AcademicYearListResponse{
		Items:      items,
		Total:      int64(len(items)),
		ActiveYear: activeYear,
	}
}

// GetByID gets an academic year by ID
func (s *AcademicYearService) GetByID(c *gin.Context, id uint) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindByID(id)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Tahun ajaran tidak ditemukan"))
	}

	if year.SchoolCode != claims.SchoolCode {
		panic(exception.NewBadRequestExceptionStruct(response.Forbidden, "Tidak memiliki akses ke tahun ajaran ini"))
	}

	return s.mapToResponse(year)
}

// GetActive gets the active academic year
func (s *AcademicYearService) GetActive(c *gin.Context) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindActive(claims.SchoolCode)
	if err != nil {
		return nil // No active year
	}

	return s.mapToResponse(year)
}

// GetBySemesterID gets an academic year by semester_id
func (s *AcademicYearService) GetBySemesterID(c *gin.Context, semesterID string) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindBySemesterID(claims.SchoolCode, semesterID)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Tahun ajaran tidak ditemukan"))
	}

	return s.mapToResponse(year)
}

// Create creates a new academic year
func (s *AcademicYearService) Create(c *gin.Context, req dto.CreateAcademicYearRequest) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	// Validate semester
	if req.Semester != 1 && req.Semester != 2 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Semester harus 1 (Ganjil) atau 2 (Genap)"))
	}

	// Validate year
	if req.YearEnd != req.YearStart+1 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Tahun akhir harus satu tahun setelah tahun awal"))
	}

	// Generate semester_id (Dapodik format)
	semesterID := fmt.Sprintf("%d%d", req.YearEnd, req.Semester)

	// Check if already exists
	if s.repo.ExistsBySemesterID(claims.SchoolCode, semesterID) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Tahun ajaran sudah ada"))
	}

	semester := entity.SemesterType(req.Semester)
	year := &entity.AcademicYear{
		SchoolCode:   claims.SchoolCode,
		SemesterID:   semesterID,
		Year:         fmt.Sprintf("%d/%d", req.YearStart, req.YearEnd),
		YearStart:    req.YearStart,
		YearEnd:      req.YearEnd,
		Semester:     semester,
		SemesterName: semester.String(),
		IsActive:     req.IsActive,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Description:  req.Description,
	}

	// If setting as active, deactivate others first
	if req.IsActive {
		s.repo.SetActive(claims.SchoolCode, 0) // Deactivate all
	}

	if err := s.repo.Create(year); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal membuat tahun ajaran"))
	}

	return s.mapToResponse(year)
}

// CreateFromDapodik creates an academic year from Dapodik semester_id
func (s *AcademicYearService) CreateFromDapodik(c *gin.Context, semesterID string) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	// Parse semester_id
	yearEnd, semester, yearStart, yearString := entity.ParseSemesterID(semesterID)
	if yearEnd == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Format semester_id tidak valid. Gunakan format YYYYS (contoh: 20251)"))
	}

	// Check if already exists
	if s.repo.ExistsBySemesterID(claims.SchoolCode, semesterID) {
		// Return existing
		year, _ := s.repo.FindBySemesterID(claims.SchoolCode, semesterID)
		return s.mapToResponse(year)
	}

	year := &entity.AcademicYear{
		SchoolCode:   claims.SchoolCode,
		SemesterID:   semesterID,
		Year:         yearString,
		YearStart:    yearStart,
		YearEnd:      yearEnd,
		Semester:     semester,
		SemesterName: semester.String(),
		IsActive:     false,
	}

	if err := s.repo.Create(year); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal membuat tahun ajaran"))
	}

	return s.mapToResponse(year)
}

// Update updates an academic year
func (s *AcademicYearService) Update(c *gin.Context, id uint, req dto.UpdateAcademicYearRequest) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindByID(id)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Tahun ajaran tidak ditemukan"))
	}

	if year.SchoolCode != claims.SchoolCode {
		panic(exception.NewBadRequestExceptionStruct(response.Forbidden, "Tidak memiliki akses ke tahun ajaran ini"))
	}

	// Update fields
	if req.StartDate != nil {
		year.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		year.EndDate = req.EndDate
	}
	if req.Description != "" {
		year.Description = req.Description
	}
	if req.IsActive != nil {
		if *req.IsActive {
			// Set this as active and deactivate others
			s.repo.SetActive(claims.SchoolCode, id)
		} else {
			year.IsActive = false
		}
	}

	if err := s.repo.Update(year); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal update tahun ajaran"))
	}

	return s.mapToResponse(year)
}

// SetActive sets an academic year as active
func (s *AcademicYearService) SetActive(c *gin.Context, id uint) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindByID(id)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Tahun ajaran tidak ditemukan"))
	}

	if year.SchoolCode != claims.SchoolCode {
		panic(exception.NewBadRequestExceptionStruct(response.Forbidden, "Tidak memiliki akses ke tahun ajaran ini"))
	}

	if err := s.repo.SetActive(claims.SchoolCode, id); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal set tahun ajaran aktif"))
	}

	year.IsActive = true
	return s.mapToResponse(year)
}

// SetActiveBySemesterID sets an academic year as active by semester_id
func (s *AcademicYearService) SetActiveBySemesterID(c *gin.Context, semesterID string) *dto.AcademicYearResponse {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindBySemesterID(claims.SchoolCode, semesterID)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Tahun ajaran tidak ditemukan"))
	}

	if err := s.repo.SetActive(claims.SchoolCode, year.ID); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal set tahun ajaran aktif"))
	}

	year.IsActive = true
	return s.mapToResponse(year)
}

// Delete deletes an academic year
func (s *AcademicYearService) Delete(c *gin.Context, id uint) {
	claims := jwt.GetDataClaims(c)

	year, err := s.repo.FindByID(id)
	if err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Tahun ajaran tidak ditemukan"))
	}

	if year.SchoolCode != claims.SchoolCode {
		panic(exception.NewBadRequestExceptionStruct(response.Forbidden, "Tidak memiliki akses ke tahun ajaran ini"))
	}

	if year.IsActive {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Tidak bisa menghapus tahun ajaran yang sedang aktif"))
	}

	if err := s.repo.Delete(id); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal menghapus tahun ajaran"))
	}
}

// GetOptions gets dropdown options for academic years
func (s *AcademicYearService) GetOptions(c *gin.Context) []dto.AcademicYearOption {
	claims := jwt.GetDataClaims(c)

	years, err := s.repo.FindBySchoolCode(claims.SchoolCode)
	if err != nil {
		return []dto.AcademicYearOption{}
	}

	options := make([]dto.AcademicYearOption, len(years))
	for i, year := range years {
		options[i] = dto.AcademicYearOption{
			Value:    year.SemesterID,
			Label:    year.GetFullName(),
			IsActive: year.IsActive,
		}
	}

	return options
}

// ===== Student Class History Operations =====

// GetStudentHistory gets class history for a student
func (s *AcademicYearService) GetStudentHistory(c *gin.Context, studentID uint) []dto.StudentClassHistoryResponse {
	claims := jwt.GetDataClaims(c)

	histories, err := s.repo.FindHistoryByStudent(claims.SchoolCode, studentID)
	if err != nil {
		return []dto.StudentClassHistoryResponse{}
	}

	responses := make([]dto.StudentClassHistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = s.mapHistoryToResponse(&h)
	}

	return responses
}

// AssignStudentToClass assigns a student to a class for the current semester
func (s *AcademicYearService) AssignStudentToClass(c *gin.Context, studentID, classID uint) *dto.StudentClassHistoryResponse {
	claims := jwt.GetDataClaims(c)

	// Get active semester
	activeSemester, err := s.repo.FindActive(claims.SchoolCode)
	if err != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Tidak ada tahun ajaran aktif. Silakan set tahun ajaran aktif terlebih dahulu."))
	}

	history := &entity.StudentClassHistory{
		SchoolCode:     claims.SchoolCode,
		StudentID:      studentID,
		ClassID:        classID,
		AcademicYearID: activeSemester.ID,
		SemesterID:     activeSemester.SemesterID,
		JoinedAt:       time.Now(),
		Status:         "active",
	}

	result, err := s.repo.FindOrCreateHistory(history)
	if err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, "Gagal assign siswa ke kelas"))
	}

	return &dto.StudentClassHistoryResponse{
		ID:               result.ID,
		StudentID:        result.StudentID,
		ClassID:          result.ClassID,
		AcademicYearID:   result.AcademicYearID,
		SemesterID:       result.SemesterID,
		SemesterFullName: activeSemester.GetFullName(),
		JoinedAt:         result.JoinedAt,
		Status:           result.Status,
	}
}

// ===== Helper Methods =====

func (s *AcademicYearService) mapToResponse(year *entity.AcademicYear) *dto.AcademicYearResponse {
	return &dto.AcademicYearResponse{
		ID:           year.ID,
		SemesterID:   year.SemesterID,
		Year:         year.Year,
		YearStart:    year.YearStart,
		YearEnd:      year.YearEnd,
		Semester:     int(year.Semester),
		SemesterName: year.SemesterName,
		FullName:     year.GetFullName(),
		IsActive:     year.IsActive,
		StartDate:    year.StartDate,
		EndDate:      year.EndDate,
		Description:  year.Description,
		CreatedAt:    year.CreatedAt,
		UpdatedAt:    year.UpdatedAt,
	}
}

func (s *AcademicYearService) mapHistoryToResponse(h *entity.StudentClassHistory) dto.StudentClassHistoryResponse {
	return dto.StudentClassHistoryResponse{
		ID:             h.ID,
		StudentID:      h.StudentID,
		ClassID:        h.ClassID,
		AcademicYearID: h.AcademicYearID,
		SemesterID:     h.SemesterID,
		JoinedAt:       h.JoinedAt,
		LeftAt:         h.LeftAt,
		Status:         h.Status,
		Notes:          h.Notes,
	}
}
