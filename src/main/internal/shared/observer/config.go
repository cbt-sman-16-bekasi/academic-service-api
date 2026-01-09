package observer

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
)

// RegisterEvent registers cache invalidation events
func RegisterEvent() {
	eventChanged := map[string]string{
		model.EventClassChanged:             cache.CacheClass,
		model.EventSubjectsChanged:          cache.CacheSubjects,
		model.EventInformationSchoolChanged: cache.CacheInformationSchool,
		model.EventTeacherChanged:           cache.CacheTeachers,
		model.EventStudentChanged:           cache.CacheStudents,
		model.EventExamChanged:              cache.CacheExam,
		model.EventExamSessionChanged:       cache.CacheExamSession,
		model.EventBankQuestionChanged:      cache.CacheBankQuestion,
		model.EventTypeExamChanged:          cache.CacheTypeExam,
		model.EventExamSessionReportChanged: cache.CacheExamSessionReport,
	}

	for eventKey, prefix := range eventChanged {
		Register(eventKey, func(p string) func() {
			return func() {
				_ = cache.DeleteByPrefix(p)
			}
		}(prefix))
	}
}
