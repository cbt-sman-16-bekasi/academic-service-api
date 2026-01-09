-- ============================================================
-- MIGRATION: Add school_code column for Multi-Tenant Support
-- Description: Menambahkan kolom school_code ke semua tabel 
--              untuk mendukung multi-tenant (multi-sekolah)
-- Author: AI Assistant
-- Date: 2025-12-28
-- ============================================================

-- ============================================================
-- SCHOOL SERVICE TABLES
-- ============================================================

-- m_class
ALTER TABLE school_service.m_class 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_class_school_code 
ON school_service.m_class(school_code);

-- m_class_code
ALTER TABLE school_service.m_class_code 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_class_code_school_code 
ON school_service.m_class_code(school_code);

-- m_class_subject
ALTER TABLE school_service.m_class_subject 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_class_subject_school_code 
ON school_service.m_class_subject(school_code);

-- type_exam
ALTER TABLE school_service.type_exam 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_type_exam_school_code 
ON school_service.type_exam(school_code);

-- exam
ALTER TABLE school_service.exam 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_school_code 
ON school_service.exam(school_code);

-- exam_member
ALTER TABLE school_service.exam_member 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_member_school_code 
ON school_service.exam_member(school_code);

-- exam_question
ALTER TABLE school_service.exam_question 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_question_school_code 
ON school_service.exam_question(school_code);

-- exam_answer_option
ALTER TABLE school_service.exam_answer_option 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_answer_option_school_code 
ON school_service.exam_answer_option(school_code);

-- master_bank_question
ALTER TABLE school_service.master_bank_question 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_master_bank_question_school_code 
ON school_service.master_bank_question(school_code);

-- bank_question
ALTER TABLE school_service.bank_question 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_bank_question_school_code 
ON school_service.bank_question(school_code);

-- bank_answer_option
ALTER TABLE school_service.bank_answer_option 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_bank_answer_option_school_code 
ON school_service.bank_answer_option(school_code);

-- exam_session
ALTER TABLE school_service.exam_session 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_session_school_code 
ON school_service.exam_session(school_code);

-- exam_session_token
ALTER TABLE school_service.exam_session_token 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_session_token_school_code 
ON school_service.exam_session_token(school_code);

-- exam_session_member
ALTER TABLE school_service.exam_session_member 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_exam_session_member_school_code 
ON school_service.exam_session_member(school_code);

-- ============================================================
-- STUDENT SERVICE TABLES
-- ============================================================

-- m_student
ALTER TABLE student_service.m_student 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_student_school_code 
ON student_service.m_student(school_code);

-- m_student_class
ALTER TABLE student_service.m_student_class 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_student_class_school_code 
ON student_service.m_student_class(school_code);

-- ============================================================
-- TEACHER SERVICE TABLES
-- ============================================================

-- m_teacher
ALTER TABLE teacher_service.m_teacher 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_teacher_school_code 
ON teacher_service.m_teacher(school_code);

-- m_teacher_class_subject
ALTER TABLE teacher_service.m_teacher_class_subject 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_teacher_class_subject_school_code 
ON teacher_service.m_teacher_class_subject(school_code);

-- ============================================================
-- CURRICULUM SERVICE TABLES
-- ============================================================

-- m_subject
ALTER TABLE curriculum_service.m_subject 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_subject_school_code 
ON curriculum_service.m_subject(school_code);

-- m_curriculum_subject
ALTER TABLE curriculum_service.m_curriculum_subject 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_m_curriculum_subject_school_code 
ON curriculum_service.m_curriculum_subject(school_code);

-- ============================================================
-- CBT SERVICE TABLES
-- ============================================================

-- student_answers
ALTER TABLE cbt_service.student_answers 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_student_answers_school_code 
ON cbt_service.student_answers(school_code);

-- student_history_taken
ALTER TABLE cbt_service.student_history_taken 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_student_history_taken_school_code 
ON cbt_service.student_history_taken(school_code);

-- history_reset_session
ALTER TABLE cbt_service.history_reset_session 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_history_reset_session_school_code 
ON cbt_service.history_reset_session(school_code);

-- history_change_score
ALTER TABLE cbt_service.history_change_score 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_history_change_score_school_code 
ON cbt_service.history_change_score(school_code);

-- suspicious_activity
ALTER TABLE cbt_service.suspicious_activity 
ADD COLUMN IF NOT EXISTS school_code VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_suspicious_activity_school_code 
ON cbt_service.suspicious_activity(school_code);

-- ============================================================
-- UPDATE EXISTING DATA WITH DEFAULT SCHOOL CODE
-- (Run this after adding columns to populate existing data)
-- ============================================================

-- Ganti 'YOUR_SCHOOL_CODE' dengan school_code yang sesuai
-- Contoh untuk sekolah yang sudah ada:
-- UPDATE school_service.m_class SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;

-- ============================================================
-- SCHOOL SERVICE
-- ============================================================
UPDATE school_service.m_class SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.m_class_code SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.m_class_subject SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.type_exam SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam_member SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam_question SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam_answer_option SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.master_bank_question SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.bank_question SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.bank_answer_option SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam_session SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam_session_token SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE school_service.exam_session_member SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;

-- ============================================================
-- STUDENT SERVICE
-- ============================================================
UPDATE student_service.m_student SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE student_service.m_student_class SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;

-- ============================================================
-- TEACHER SERVICE
-- ============================================================
UPDATE teacher_service.m_teacher SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE teacher_service.m_teacher_class_subject SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;

-- ============================================================
-- CURRICULUM SERVICE
-- ============================================================
UPDATE curriculum_service.m_subject SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE curriculum_service.m_curriculum_subject SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;

-- ============================================================
-- CBT SERVICE
-- ============================================================
UPDATE cbt_service.student_answers SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE cbt_service.student_history_taken SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE cbt_service.history_reset_session SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE cbt_service.history_change_score SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;
UPDATE cbt_service.suspicious_activity SET school_code = 'db74a42e-23a7-4cd2-bbe5-49cf79f86453' WHERE school_code IS NULL;

-- ============================================================
-- END OF MIGRATION
-- ============================================================

