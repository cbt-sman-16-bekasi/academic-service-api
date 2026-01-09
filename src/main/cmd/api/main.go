package main

import (
	_ "github.com/Sistem-Informasi-Akademik/academic-system-information-service/docs"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/app"
)

// @title Academic System API
// @version 1.0
// @description Academic System Information Service - API Documentation

// @contact.name Development Team
// @contact.url https://github.com/Sistem-Informasi-Akademik
// @contact.email dev@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @schemes https http
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Use "Bearer {your_token}" to authenticate

// @tag.name Auth
// @tag.description Authentication and Profile operations

// @tag.name School
// @tag.description School management operations

// @tag.name Student
// @tag.description Student management operations

// @tag.name Teacher
// @tag.description Teacher management operations

// @tag.name CBT
// @tag.description Computer-Based Test operations

// @tag.name Exam
// @tag.description Exam configuration operations

// @tag.name Exam Question
// @tag.description Question management operations

// @tag.name Exam Session
// @tag.description Exam session management operations

// @tag.name Type Exam
// @tag.description Exam type master data operations
func main() {
	app.Run()
}
