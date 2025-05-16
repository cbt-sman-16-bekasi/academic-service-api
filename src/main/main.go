package main

import (
	_ "github.com/Sistem-Informasi-Akademik/academic-system-information-service/docs"
	_ "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity"
	_ "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/routes"
	_ "github.com/yon-module/yon-framework/config"
)

// @title CBT Documentation
// @version 1.0
// @description This is documentation technical

// @contact.name CBT Team
// @contact.url https://sman16bekasi.id
// @contact.email cbt@sman16bekasi.id

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @schemes https
// @host api.sman16bekasi.id/gateway-service
// @path /

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Use "Bearer {your_token}" to authenticate. Some endpoints are public and do not require authentication.

// @tag.name Auth
// @tag.description Operation about Authentication and Profile

// @tag.name School
// @tag.description Operation about school data

// @tag.name Student
// @tag.description Operation about student

// @tag.name Teacher
// @tag.description Operation about teacher

// @tag.name CBT
// @tag.description Operation about CBT configuration

// @tag.name Exam
// @tag.description Operation about Exam configuration

// @tag.name Exam Question
// @tag.description Operation about Management Question configuration

// @tag.name Exam Session
// @tag.description Operation about Management Question Session configuration

// @tag.name Type Exam
// @tag.description Operation about Mastering type exam
func main() {
}
