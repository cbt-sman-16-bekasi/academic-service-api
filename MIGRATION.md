# Migration Guide: Legacy to Modular Architecture

## Overview

This document outlines the strategy for gradually migrating the existing codebase from the legacy flat structure to the new modular architecture.

## Current State

### Legacy Structure (To Be Migrated)
```
src/main/
├── controllers/          # HTTP handlers
├── service/academic/     # Business logic
├── repository/           # Data access
├── model/               # Entities & DTOs
└── routes/              # Route definitions
```

### Target Structure
```
src/main/internal/modules/{module_name}/
├── entity/              # Domain entities
├── dto/                 # Request/Response DTOs
├── repository/          # Data access layer
├── service/             # Business logic
├── handler/             # HTTP handlers
└── routes.go            # Module route registration
```

## Migration Priority

### Phase 1: New Features (DONE ✅)
New features like SPP, PPDB, E-Learning should be developed directly in the modular structure.

### Phase 2: Independent Modules (Next)
Modules with minimal cross-dependencies:
1. **Dapodik** - External integration, minimal dependencies
2. **Curriculum** - Master data, used by others but doesn't depend on much

### Phase 3: Core Modules (Later)
Modules with more dependencies:
1. **Class** - Depends on School, used by Student, Teacher
2. **Teacher** - Depends on User, Class
3. **Student** - Depends on User, Class

### Phase 4: Complex Modules (Final)
1. **Exam/CBT** - Most complex, depends on Student, Class, Teacher
2. **Auth** - Core authentication

## Migration Steps Per Module

### Step 1: Create Module Structure
```bash
mkdir -p src/main/internal/modules/{module}/
mkdir -p src/main/internal/modules/{module}/entity
mkdir -p src/main/internal/modules/{module}/dto
mkdir -p src/main/internal/modules/{module}/repository
mkdir -p src/main/internal/modules/{module}/service
mkdir -p src/main/internal/modules/{module}/handler
```

### Step 2: Move Entities
```go
// FROM: src/main/model/entity/{module}/{entity}.go
// TO:   src/main/internal/modules/{module}/entity/{entity}.go

// Update package name
package entity  // was: package {module}

// Update imports in other files
```

### Step 3: Move DTOs
```go
// FROM: src/main/model/dto/request/{module}_request/*.go
// FROM: src/main/model/dto/response/{module}_response/*.go
// TO:   src/main/internal/modules/{module}/dto/*.go
```

### Step 4: Move Repository
```go
// FROM: src/main/repository/{module}_repository/*.go
// TO:   src/main/internal/modules/{module}/repository/*.go
```

### Step 5: Move Service
```go
// FROM: src/main/service/academic/{module}_service/*.go
// TO:   src/main/internal/modules/{module}/service/*.go
```

### Step 6: Create Handler (Rename Controller)
```go
// FROM: src/main/controllers/{module}.go
// TO:   src/main/internal/modules/{module}/handler/{module}_handler.go

// Rename pattern:
// - {Module}Controller → {Module}Handler
// - All methods stay the same
```

### Step 7: Create Module Routes
```go
// Create: src/main/internal/modules/{module}/routes.go

package {module}

func RegisterRoutes(router *gin.RouterGroup) {
    handler := handler.New{Module}Handler()
    
    group := router.Group("/{module}")
    group.Use(jwt.AuthMiddleware())
    group.Use(jwt.SchoolScopeMiddleware())
    {
        group.GET("/all", handler.GetAll)
        group.GET("/:id", handler.GetDetail)
        // ... etc
    }
}
```

### Step 8: Register in Main Routes
```go
// In src/main/routes/routes.go

import "{module}"

func init() {
    // Remove legacy route
    // server.AddRoutes({module}Routes)  // DELETE THIS
    
    // Add modular route
    server.AddRoute({module}.RegisterRoutes)  // ADD THIS
}
```

### Step 9: Update Migration Registration
```go
// In src/main/model/entity/migration.go
// Update entity imports to use new module path
```

### Step 10: Cleanup
- Remove old files from legacy locations
- Run `go build ./...` to verify
- Run tests

## Cross-Module Dependencies

When Module A depends on Module B's entity:

### Option 1: Import Directly (Simple)
```go
import (
    schoolEntity "github.com/.../internal/modules/school/entity"
)

type Student struct {
    ClassID uint
    Class   schoolEntity.Class `gorm:"foreignKey:ClassID"`
}
```

### Option 2: Shared Kernel (For Core Entities)
Move truly shared entities to:
```
src/main/internal/shared/domain/
├── user.go      # User, Role
└── school.go    # School, SchoolCode
```

## Testing Migration

After migrating each module:

1. **Build Check**
   ```bash
   go build ./...
   ```

2. **Run Application**
   ```bash
   go run src/main/cmd/api/main.go
   ```

3. **Verify Routes**
   Check that routes are registered in startup logs

4. **API Testing**
   Test all endpoints of the migrated module

## Rollback Strategy

If issues arise:
1. Keep legacy code until fully verified
2. Use feature flags if needed
3. Gradual traffic shift if in production

## Progress Tracking

| Module | Status | Notes |
|--------|--------|-------|
| SPP | ✅ New | Developed in modular structure |
| PPDB | 📌 Pending | New module |
| E-Learning | 📌 Pending | New module |
| Dapodik | ⏳ To Migrate | Phase 2 |
| Curriculum | ⏳ To Migrate | Phase 2 |
| Class | ⏳ To Migrate | Phase 3 |
| Teacher | ⏳ To Migrate | Phase 3 |
| Student | ⏳ To Migrate | Phase 3 |
| Exam/CBT | ⏳ To Migrate | Phase 4 |
| Auth | ⏳ To Migrate | Phase 4 |

