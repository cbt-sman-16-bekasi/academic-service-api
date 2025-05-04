package teacher_request

type TeacherModifyRequest struct {
	Nuptk    string `json:"nuptk"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	IsAccess bool   `json:"isAccess"`
}

type TeacherMappingSubjectClass struct {
	TeacherId uint   `json:"teacherId"`
	SubjectId string `json:"subjectId"`
	ClassId   []uint `json:"classId"`
}
