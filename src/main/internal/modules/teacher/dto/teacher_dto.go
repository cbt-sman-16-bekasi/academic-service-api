package dto

// ----- Request DTOs -----

// TeacherModifyRequest for create/update teacher
type TeacherModifyRequest struct {
	Nuptk    string `json:"nuptk" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Gender   string `json:"gender"`
	IsAccess bool   `json:"isAccess"`
}

// TeacherMappingSubjectClass for mapping teacher to class subjects
type TeacherMappingSubjectClass struct {
	TeacherId uint   `json:"teacherId" binding:"required"`
	SubjectId string `json:"subjectId" binding:"required"`
	ClassId   []uint `json:"classId" binding:"required"`
}

// ----- Response DTOs -----

// GeneralLabelKeyResponse for key-label pairs
type GeneralLabelKeyResponse struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// TeacherDetailResponse for teacher detail
type TeacherDetailResponse struct {
	Nuptk    string                  `json:"nuptk"`
	Name     string                  `json:"name"`
	Username string                  `json:"username"`
	Gender   string                  `json:"gender"`
	Role     GeneralLabelKeyResponse `json:"role"`
	IsAccess bool                    `json:"isAccess"`
}

