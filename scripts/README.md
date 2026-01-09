# Migration Scripts

## FX Dependency Injection Migration

This directory contains scripts to help migrate modules to use Uber FX for dependency injection.

### Available Scripts

#### `migrate_fx/main.go`

A Go-based tool that generates `module.go` files for each module with the FX pattern.

**Usage:**
```bash
go run scripts/migrate_fx/main.go
```

**What it does:**
1. Creates `module.go` files for each module
2. Sets up FX providers for Repository → Service → Handler chain
3. Registers routes with DI

#### `migrate_to_fx.sh`

A bash script that attempts to add DI constructors to existing files.

**Usage:**
```bash
./scripts/migrate_to_fx.sh
```

### DI Pattern

Each module follows this pattern:

```
module/
├── repository/
│   └── xxx_repository.go     # NewXxxRepository() + NewXxxRepositoryWithDB(db)
├── service/
│   └── xxx_service.go        # NewXxxService() + NewXxxServiceWithDeps(repo)
├── handler/
│   └── xxx_handler.go        # NewXxxHandler() + NewXxxHandlerWithService(svc)
├── routes.go                 # RegisterRoutes() + RegisterRoutesWithDI(router, handler)
└── module.go                 # FX Module definition
```

### Constructor Patterns

**Repository:**
```go
// Legacy (non-DI)
func NewXxxRepository() *XxxRepository {
    return &XxxRepository{db: database.GetDB()}
}

// DI-friendly
func NewXxxRepositoryWithDB(db *gorm.DB) *XxxRepository {
    return &XxxRepository{db: db}
}
```

**Service:**
```go
// Legacy (non-DI)
func NewXxxService() *XxxService {
    return &XxxService{repo: repository.NewXxxRepository()}
}

// DI-friendly
func NewXxxServiceWithDeps(repo *repository.XxxRepository) *XxxService {
    return &XxxService{repo: repo}
}
```

**Handler:**
```go
// Legacy (non-DI)
func NewXxxHandler() *XxxHandler {
    return &XxxHandler{service: service.NewXxxService()}
}

// DI-friendly
func NewXxxHandlerWithService(svc *service.XxxService) *XxxHandler {
    return &XxxHandler{service: svc}
}
```

**Routes:**
```go
// Legacy (non-DI)
func RegisterRoutes(router *gin.RouterGroup) {
    h := handler.NewXxxHandler()
    registerRoutes(router, h)
}

// DI-friendly
func RegisterRoutesWithDI(router *gin.RouterGroup, h *handler.XxxHandler) {
    registerRoutes(router, h)
}

// Common logic
func registerRoutes(router *gin.RouterGroup, h *handler.XxxHandler) {
    // Route definitions
}
```

### Module Definition

```go
package xxx

import (
    "go.uber.org/fx"
    "gorm.io/gorm"
)

func NewXxxRepository(db *gorm.DB) *repository.XxxRepository {
    return repository.NewXxxRepositoryWithDB(db)
}

func NewXxxService(repo *repository.XxxRepository) *service.XxxService {
    return service.NewXxxServiceWithDeps(repo)
}

func NewXxxHandler(svc *service.XxxService) *handler.XxxHandler {
    return handler.NewXxxHandlerWithService(svc)
}

func RegisterXxxRoutes(rg *types.RouterGroup, h *handler.XxxHandler) {
    RegisterRoutesWithDI(rg.RouterGroup, h)
}

var Module = fx.Module("xxx",
    fx.Provide(
        NewXxxRepository,
        NewXxxService,
        NewXxxHandler,
    ),
    fx.Invoke(RegisterXxxRoutes),
)
```

### Backward Compatibility

Both legacy constructors (without DI) and DI-friendly constructors are maintained, allowing:
1. Gradual migration
2. Easy testing with mocks
3. Flexible initialization based on context

