package observer

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/redisstore"
)

func RegisterEvent() {
	eventChanged := map[string]string{
		model.EventClassChanged:             redisstore.CacheClass,
		model.EventSubjectsChanged:          redisstore.CacheSubjects,
		model.EventInformationSchoolChanged: redisstore.CacheInformationSchool,
		model.EventTeacherChanged:           redisstore.CacheTeachers,
		model.EventStudentChanged:           redisstore.CacheStudents,
		model.EventExamChanged:              redisstore.CacheExam,
		model.EventExamSessionChanged:       redisstore.CacheExamSession,
		model.EventBankQuestionChanged:      redisstore.CacheBankQuestion,
		model.EventTypeExamChanged:          redisstore.CacheTypeExam,
		model.EventExamSessionReportChanged: redisstore.CacheExamSessionReport,
	}

	for eventKey, prefix := range eventChanged {
		Register(eventKey, func(p string) func() {
			return func() {
				_ = redisstore.DeleteByPrefix(prefix)
			}
		}(prefix))
	}
}
