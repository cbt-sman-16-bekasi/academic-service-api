package redisstore

import "time"

const (
	CacheInformationSchool = "cache::information_school::"
	CacheClass             = "cache::class::"
	CacheSubjects          = "cache::subjects::"
	CacheTeachers          = "cache::teachers::"
	CacheStudents          = "cache::students::"
	CacheTypeExam          = "cache::type_exam::"
	CacheBankQuestion      = "cache::bank_question::"
	CacheExam              = "cache::exam::"
	CacheExamSession       = "cache::exam_session::"
	CacheExamSessionReport = "cache::exam_session_report::"
)

var TtlDuration = time.Duration(30 * time.Second)
var TtlOneDay = time.Duration(24 * time.Hour)
