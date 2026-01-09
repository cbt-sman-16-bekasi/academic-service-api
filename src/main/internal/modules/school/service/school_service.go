package service

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/school/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/gin-gonic/gin"
)

type SchoolService struct {
	repo *repository.SchoolRepository
}

// NewSchoolService creates service with injected repository
func NewSchoolService(repo *repository.SchoolRepository) *SchoolService {
	return &SchoolService{repo: repo}
}

// GetSchoolDetail retrieves school detail
func (s *SchoolService) GetSchoolDetail(c *gin.Context) dto.DetailSchoolResponse {
	claims := jwt.GetDataClaims(c)
	schoolData, err := s.repo.FindSchoolByCode(claims.SchoolCode)
	if err != nil || schoolData.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "School not found"))
	}

	return dto.DetailSchoolResponse{
		SchoolName:        schoolData.SchoolName,
		Id:                schoolData.SchoolCode,
		Address:           schoolData.Address,
		Email:             schoolData.Email,
		Phone:             schoolData.Phone,
		Logo:              schoolData.Logo,
		Npsn:              schoolData.NPSN,
		Nss:               schoolData.NSS,
		Banner:            schoolData.Banner,
		PrincipalName:     schoolData.PrincipalName,
		PrincipalNIP:      schoolData.PrincipalNIP,
		VicePrincipalName: schoolData.VicePrincipalName,
		VicePrincipalNIP:  schoolData.VicePrincipalNIP,
		LevelOfEducation:  "SMA",
	}
}

// UpdateSchool updates school info
func (s *SchoolService) UpdateSchool(c *gin.Context, req dto.ModifySchoolRequest) dto.DetailSchoolResponse {
	claims := jwt.GetDataClaims(c)
	schoolData, err := s.repo.FindSchoolByCode(claims.SchoolCode)
	if err != nil || schoolData.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "School not found"))
	}

	schoolData.SchoolName = req.SchoolName
	schoolData.Logo = req.Logo
	schoolData.NSS = req.Nss
	schoolData.NPSN = req.Npsn
	schoolData.Phone = req.Phone
	schoolData.Email = req.Email
	schoolData.Address = req.Address
	schoolData.Banner = req.Banner
	schoolData.PrincipalName = req.PrincipalName
	schoolData.VicePrincipalName = req.VicePrincipalName
	schoolData.PrincipalNIP = req.PrincipalNIP
	schoolData.VicePrincipalNIP = req.VicePrincipalNIP

	if err := s.repo.UpdateSchool(schoolData); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.DetailSchoolResponse{
		SchoolName:        schoolData.SchoolName,
		Id:                schoolData.SchoolCode,
		Address:           schoolData.Address,
		Email:             schoolData.Email,
		Phone:             schoolData.Phone,
		Logo:              schoolData.Logo,
		Npsn:              schoolData.NPSN,
		Nss:               schoolData.NSS,
		Banner:            schoolData.Banner,
		PrincipalName:     schoolData.PrincipalName,
		PrincipalNIP:      schoolData.PrincipalNIP,
		VicePrincipalName: schoolData.VicePrincipalName,
		VicePrincipalNIP:  schoolData.VicePrincipalNIP,
		LevelOfEducation:  "SMA",
	}
}

// GetAllClassCodes returns all class codes
func (s *SchoolService) GetAllClassCodes(ctx *gin.Context) []dto.ClassCodeResponse {
	claims := jwt.GetDataClaims(ctx)
	classCodes := s.repo.AllClassCodeSchool(claims.SchoolCode)
	var result []dto.ClassCodeResponse
	for _, cc := range classCodes {
		result = append(result, dto.ClassCodeResponse{
			Code: cc.Code,
			Name: cc.Name,
		})
	}
	return result
}

// GetAllSubjects returns all subjects
func (s *SchoolService) GetAllSubjects() []dto.SubjectResponse {
	subjects := s.repo.GetAllSubjects()
	var result []dto.SubjectResponse
	for _, subj := range subjects {
		result = append(result, dto.SubjectResponse{
			ID:          subj.ID,
			Code:        subj.Code,
			Subject:     subj.Subject,
			SubjectType: subj.SubjectType,
			Description: subj.Description,
		})
	}
	return result
}

// GetAllClassSubjects returns paginated class-subject mappings
func (s *SchoolService) GetAllClassSubjects(ctx *gin.Context, req pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(ctx)

	filter := map[string]interface{}{}

	if req.Filter != nil {
		filter = *req.Filter
	}

	filter["school_code"] = claims.SchoolCode
	req.Filter = &filter
	
	page := database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetRequest(&req).
		SetPreloads("DetailSubject", "DetailClassCode", "DetailClassCode.ClassMember").
		SetModel([]school.ClassSubject{}).
		FindAllPaging()
	return page
}

// GetClassSubjectDetail returns class-subject detail
func (s *SchoolService) GetClassSubjectDetail(id uint) dto.DetailClassSubjectResponse {
	cs, err := s.repo.FindClassSubjectByID(id)
	if err != nil || cs.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Class subject not found"))
	}

	return dto.DetailClassSubjectResponse{
		ID: cs.ID,
		ClassCode: dto.GeneralLabelKeyResponse{
			Key:   cs.ClassCode,
			Label: cs.DetailClassCode.Name,
		},
		Subject: dto.GeneralLabelKeyResponse{
			Key:   cs.SubjectCode,
			Label: cs.DetailSubject.Subject,
		},
	}
}

// CreateClassSubject creates new class-subject mapping
func (s *SchoolService) CreateClassSubject(req dto.ModifyClassSubjectRequest) dto.DetailClassSubjectResponse {
	// Verify subject exists
	subject, err := s.repo.FindSubjectByCode(req.SubjectCode)
	if err != nil || subject.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Subject code not exist"))
	}

	// Check for duplicate
	existing, _ := s.repo.FindClassSubjectByClassAndSubject(req.ClassCode, req.SubjectCode)
	if existing != nil && existing.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "This class-subject mapping already exists"))
	}

	newData := &school.ClassSubject{
		SubjectCode: req.SubjectCode,
		ClassCode:   req.ClassCode,
	}

	if err := s.repo.CreateClassSubject(newData); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.DetailClassSubjectResponse{
		ID: newData.ID,
		ClassCode: dto.GeneralLabelKeyResponse{
			Key: req.ClassCode,
		},
		Subject: dto.GeneralLabelKeyResponse{
			Key:   req.SubjectCode,
			Label: subject.Subject,
		},
	}
}

// UpdateClassSubject updates class-subject mapping
func (s *SchoolService) UpdateClassSubject(id uint, req dto.ModifyClassSubjectRequest) dto.DetailClassSubjectResponse {
	existing, err := s.repo.FindClassSubjectByID(id)
	if err != nil || existing.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Class subject not found"))
	}

	// Check for duplicate
	duplicate, _ := s.repo.FindClassSubjectByClassAndSubject(req.ClassCode, req.SubjectCode)
	if duplicate != nil && duplicate.ID != 0 && duplicate.ID != id {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "This class-subject mapping already exists"))
	}

	existing.ClassCode = req.ClassCode
	existing.SubjectCode = req.SubjectCode

	if err := s.repo.UpdateClassSubject(existing); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.DetailClassSubjectResponse{
		ID: id,
		ClassCode: dto.GeneralLabelKeyResponse{
			Key: req.ClassCode,
		},
		Subject: dto.GeneralLabelKeyResponse{
			Key: req.SubjectCode,
		},
	}
}

// DeleteClassSubject deletes class-subject mapping
func (s *SchoolService) DeleteClassSubject(id uint) {
	existing, err := s.repo.FindClassSubjectByID(id)
	if err != nil || existing.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Class subject not found"))
	}

	if err := s.repo.DeleteClassSubject(id); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}

// GetDashboard returns dashboard data based on user role
func (s *SchoolService) GetDashboard(c *gin.Context) dto.DashboardResponse {
	claims := jwt.GetDataClaims(c)

	if claims.Role != "ADMIN" {
		dashboard, _ := s.repo.GetDashboardTeacher(claims.Id)
		if dashboard == nil {
			return dto.DashboardResponse{}
		}
		return dto.DashboardResponse{
			TotalClass:       dashboard.TotalClasses,
			TotalSubject:     dashboard.TotalClassSubjects,
			TotalStudent:     dashboard.TotalStudents,
			TotalExam:        dashboard.TotalExams,
			TotalSessionExam: dashboard.TotalExamSessions,
			TotalReportExam:  dashboard.TotalReport,
			TotalAccess:      dashboard.TotalUsers,
		}
	}

	dashboard, _ := s.repo.GetDashboardAdmin()
	if dashboard == nil {
		return dto.DashboardResponse{}
	}
	return dto.DashboardResponse{
		TotalClass:       dashboard.TotalClasses,
		TotalSubject:     dashboard.TotalClassSubjects,
		TotalStudent:     dashboard.TotalStudents,
		TotalExam:        dashboard.TotalExams,
		TotalSessionExam: dashboard.TotalExamSessions,
		TotalReportExam:  dashboard.TotalReport,
		TotalAccess:      dashboard.TotalUsers,
	}
}

// LoadConfiguration loads system configuration for a given origin
func (s *SchoolService) LoadConfiguration(origin string) *dto.ConfigurationResponse {
	config, err := s.repo.FindSystemConfigByOrigin(origin)
	if err != nil || config.ID == 0 {
		return nil
	}

	schoolData, _ := s.repo.FindSchoolByCode(config.SchoolCode)

	var schoolInfo *dto.SchoolBriefInfo
	if schoolData != nil && schoolData.ID != 0 {
		schoolInfo = &dto.SchoolBriefInfo{
			SchoolCode: schoolData.SchoolCode,
			SchoolName: schoolData.SchoolName,
			Logo:       schoolData.Logo,
			Banner:     schoolData.Banner,
			Address:    schoolData.Address,
		}
	}

	return &dto.ConfigurationResponse{
		ClientID:    config.ClientId,
		ThemePreset: config.ThemePreset,
		ThemeCustom: config.ThemeCustom,
		SchoolInfo:  schoolInfo,
	}
}
