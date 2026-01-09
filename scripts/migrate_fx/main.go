package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const modulesDir = "src/main/internal/modules"

// ModuleConfig contains module-specific configuration
type ModuleConfig struct {
	Name          string // lowercase: auth, class, etc.
	CapName       string // capitalized: Auth, Class, etc.
	RepoStruct    string // e.g., AuthRepository
	ServiceStruct string // e.g., AuthService
	HandlerStruct string // e.g., AuthHandler
	HasRepository bool
	HasService    bool
	HasHandler    bool
	HasRoutes     bool
}

// Module definitions with correct struct names
var modules = []ModuleConfig{
	{Name: "auth", CapName: "Auth", RepoStruct: "AuthRepository", ServiceStruct: "AuthService", HandlerStruct: "AuthHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "class", CapName: "Class", RepoStruct: "ClassRepository", ServiceStruct: "ClassService", HandlerStruct: "ClassHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "curriculum", CapName: "Curriculum", RepoStruct: "SubjectRepository", ServiceStruct: "SubjectService", HandlerStruct: "SubjectHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "dapodik", CapName: "Dapodik", RepoStruct: "DapodikRepository", ServiceStruct: "DapodikService", HandlerStruct: "DapodikHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "exam", CapName: "Exam", RepoStruct: "ExamRepository", ServiceStruct: "ExamService", HandlerStruct: "ExamHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "school", CapName: "School", RepoStruct: "SchoolRepository", ServiceStruct: "SchoolService", HandlerStruct: "SchoolHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "spp", CapName: "Spp", RepoStruct: "PaymentRepository", ServiceStruct: "PaymentService", HandlerStruct: "PaymentHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
	{Name: "storage", CapName: "Storage", RepoStruct: "", ServiceStruct: "", HandlerStruct: "StorageHandler", HasRepository: false, HasService: false, HasHandler: true, HasRoutes: true},
	{Name: "teacher", CapName: "Teacher", RepoStruct: "TeacherRepository", ServiceStruct: "TeacherService", HandlerStruct: "TeacherHandler", HasRepository: true, HasService: true, HasHandler: true, HasRoutes: true},
}

const moduleTemplate = `package {{.Name}}

import (
{{- if .HasHandler}}
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/{{.Name}}/handler"
{{- end}}
{{- if .HasRepository}}
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/{{.Name}}/repository"
{{- end}}
{{- if .HasService}}
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/{{.Name}}/service"
{{- end}}
{{- if .HasRoutes}}
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
{{- end}}
	"go.uber.org/fx"
{{- if .HasRepository}}
	"gorm.io/gorm"
{{- end}}
)

{{if .HasRepository}}
// New{{.RepoStruct}} creates repository with injected DB
func New{{.RepoStruct}}(db *gorm.DB) *repository.{{.RepoStruct}} {
	return repository.New{{.RepoStruct}}WithDB(db)
}
{{end}}

{{if .HasService}}
// New{{.ServiceStruct}} creates service with injected dependencies
func New{{.ServiceStruct}}(repo *repository.{{.RepoStruct}}) *service.{{.ServiceStruct}} {
	return service.New{{.ServiceStruct}}WithDeps(repo)
}
{{end}}

{{if .HasHandler}}
// New{{.HandlerStruct}} creates handler with injected service
{{- if .HasService}}
func New{{.HandlerStruct}}(svc *service.{{.ServiceStruct}}) *handler.{{.HandlerStruct}} {
	return handler.New{{.HandlerStruct}}WithService(svc)
}
{{- else}}
func New{{.HandlerStruct}}() *handler.{{.HandlerStruct}} {
	return handler.New{{.HandlerStruct}}()
}
{{- end}}
{{end}}

{{if .HasRoutes}}
// Register{{.CapName}}Routes registers routes with DI
func Register{{.CapName}}Routes(rg *types.RouterGroup, h *handler.{{.HandlerStruct}}) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}
{{end}}

// Module provides all {{.Name}} module dependencies
var Module = fx.Module("{{.Name}}",
	fx.Provide(
{{- if .HasRepository}}
		New{{.RepoStruct}},
{{- end}}
{{- if .HasService}}
		New{{.ServiceStruct}},
{{- end}}
{{- if .HasHandler}}
		New{{.HandlerStruct}},
{{- end}}
	),
{{- if .HasRoutes}}
	fx.Invoke(Register{{.CapName}}Routes),
{{- end}}
)
`

func main() {
	fmt.Println("FX Migration Tool")
	fmt.Println("=================")
	fmt.Println()

	tmpl, err := template.New("module").Parse(moduleTemplate)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	for _, mod := range modules {
		moduleFile := filepath.Join(modulesDir, mod.Name, "module.go")

		// Check if file already exists
		if _, err := os.Stat(moduleFile); err == nil {
			fmt.Printf("✓ %s/module.go already exists, skipping\n", mod.Name)
			continue
		}

		// Create module.go
		f, err := os.Create(moduleFile)
		if err != nil {
			fmt.Printf("✗ Error creating %s: %v\n", moduleFile, err)
			continue
		}

		err = tmpl.Execute(f, mod)
		f.Close()
		if err != nil {
			fmt.Printf("✗ Error writing %s: %v\n", moduleFile, err)
			continue
		}

		fmt.Printf("✓ Created %s/module.go\n", mod.Name)
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Add DI constructors to repository files:")
	fmt.Println("   func NewXxxRepositoryWithDB(db *gorm.DB) *XxxRepository { ... }")
	fmt.Println()
	fmt.Println("2. Add DI constructors to service files:")
	fmt.Println("   func NewXxxServiceWithDeps(repo *XxxRepository) *XxxService { ... }")
	fmt.Println()
	fmt.Println("3. Add DI constructors to handler files:")
	fmt.Println("   func NewXxxHandlerWithService(svc *XxxService) *XxxHandler { ... }")
	fmt.Println()
	fmt.Println("4. Add RegisterRoutesWithDI to routes.go files")
	fmt.Println()
	fmt.Println("5. Update app.go to include all modules")
	fmt.Println()
	fmt.Println("Run: go build ./... to verify")
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
