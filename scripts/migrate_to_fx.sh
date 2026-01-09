#!/bin/bash

# FX Migration Script
# This script adds DI-friendly constructors to modules

set -e

MODULES_DIR="src/main/internal/modules"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to add WithDB constructor to repository
add_repository_constructor() {
    local module=$1
    local repo_file="$MODULES_DIR/$module/repository/${module}_repository.go"
    
    if [ ! -f "$repo_file" ]; then
        # Try alternative naming
        for f in "$MODULES_DIR/$module/repository/"*.go; do
            if [ -f "$f" ]; then
                repo_file="$f"
                break
            fi
        done
    fi
    
    if [ ! -f "$repo_file" ]; then
        log_warn "Repository file not found for $module"
        return 1
    fi
    
    # Check if already has WithDB constructor
    if grep -q "WithDB" "$repo_file"; then
        log_info "Repository $module already has DI constructor"
        return 0
    fi
    
    log_info "Adding DI constructor to repository: $repo_file"
    
    # Get the repository struct name
    local struct_name=$(grep -oP 'type \K\w+Repository' "$repo_file" | head -1)
    
    if [ -z "$struct_name" ]; then
        log_warn "Could not find repository struct in $repo_file"
        return 1
    fi
    
    # Add the constructor after NewXxxRepository function
    sed -i "/^func New${struct_name}() \*${struct_name} {/,/^}/a\\
\\
// New${struct_name}WithDB creates repository with injected DB (for DI)\\
func New${struct_name}WithDB(db *gorm.DB) *${struct_name} {\\
	return \&${struct_name}{db: db}\\
}" "$repo_file"
    
    return 0
}

# Function to add WithDeps constructor to service
add_service_constructor() {
    local module=$1
    local svc_file="$MODULES_DIR/$module/service/${module}_service.go"
    
    if [ ! -f "$svc_file" ]; then
        # Try alternative naming
        for f in "$MODULES_DIR/$module/service/"*.go; do
            if [ -f "$f" ] && [[ ! "$f" == *"_test.go" ]]; then
                svc_file="$f"
                break
            fi
        done
    fi
    
    if [ ! -f "$svc_file" ]; then
        log_warn "Service file not found for $module"
        return 1
    fi
    
    # Check if already has WithDeps constructor
    if grep -q "WithDeps" "$svc_file"; then
        log_info "Service $module already has DI constructor"
        return 0
    fi
    
    log_info "Adding DI constructor to service: $svc_file"
    
    # Get the service struct name
    local struct_name=$(grep -oP 'type \K\w+Service' "$svc_file" | head -1)
    
    if [ -z "$struct_name" ]; then
        log_warn "Could not find service struct in $svc_file"
        return 1
    fi
    
    # Get the repo field name
    local repo_field=$(grep -oP '\s+\w+\s+\*repository\.\w+Repository' "$svc_file" | head -1 | awk '{print $1}')
    local repo_type=$(grep -oP '\*repository\.\w+Repository' "$svc_file" | head -1)
    
    if [ -z "$repo_field" ] || [ -z "$repo_type" ]; then
        log_warn "Could not determine repo field for $svc_file"
        return 1
    fi
    
    # Add constructor
    sed -i "/^func New${struct_name}() \*${struct_name} {/,/^}/a\\
\\
// New${struct_name}WithDeps creates service with injected dependencies (for DI)\\
func New${struct_name}WithDeps(${repo_field} ${repo_type}) *${struct_name} {\\
	return \&${struct_name}{${repo_field}: ${repo_field}}\\
}" "$svc_file"
    
    return 0
}

# Function to add WithService constructor to handler
add_handler_constructor() {
    local module=$1
    local handler_file="$MODULES_DIR/$module/handler/${module}_handler.go"
    
    if [ ! -f "$handler_file" ]; then
        # Try alternative naming
        for f in "$MODULES_DIR/$module/handler/"*.go; do
            if [ -f "$f" ] && [[ ! "$f" == *"_test.go" ]]; then
                handler_file="$f"
                break
            fi
        done
    fi
    
    if [ ! -f "$handler_file" ]; then
        log_warn "Handler file not found for $module"
        return 1
    fi
    
    # Check if already has WithService constructor
    if grep -q "WithService" "$handler_file"; then
        log_info "Handler $module already has DI constructor"
        return 0
    fi
    
    log_info "Adding DI constructor to handler: $handler_file"
    
    # Get the handler struct name
    local struct_name=$(grep -oP 'type \K\w+Handler' "$handler_file" | head -1)
    
    if [ -z "$struct_name" ]; then
        log_warn "Could not find handler struct in $handler_file"
        return 1
    fi
    
    # Get the service field
    local svc_field=$(grep -oP '\s+\w+\s+\*service\.\w+Service' "$handler_file" | head -1 | awk '{print $1}')
    local svc_type=$(grep -oP '\*service\.\w+Service' "$handler_file" | head -1)
    
    if [ -z "$svc_field" ] || [ -z "$svc_type" ]; then
        log_warn "Could not determine service field for $handler_file"
        return 1
    fi
    
    # Add constructor
    sed -i "/^func New${struct_name}() \*${struct_name} {/,/^}/a\\
\\
// New${struct_name}WithService creates handler with injected service (for DI)\\
func New${struct_name}WithService(svc ${svc_type}) *${struct_name} {\\
	return \&${struct_name}{${svc_field}: svc}\\
}" "$handler_file"
    
    return 0
}

# Function to update routes to support DI
update_routes() {
    local module=$1
    local routes_file="$MODULES_DIR/$module/routes.go"
    
    if [ ! -f "$routes_file" ]; then
        log_warn "Routes file not found for $module"
        return 1
    fi
    
    # Check if already has WithDI function
    if grep -q "RegisterRoutesWithDI" "$routes_file"; then
        log_info "Routes $module already has DI support"
        return 0
    fi
    
    log_info "Adding DI support to routes: $routes_file"
    
    # This is complex - skip for now, manual update needed
    log_warn "Routes need manual update for $module"
    return 0
}

# Function to create module.go file
create_module_file() {
    local module=$1
    local module_file="$MODULES_DIR/$module/module.go"
    
    if [ -f "$module_file" ]; then
        log_info "Module file already exists for $module"
        return 0
    fi
    
    log_info "Creating module.go for $module"
    
    # Capitalize first letter
    local Module=$(echo "$module" | sed 's/./\U&/')
    
    cat > "$module_file" << EOF
package $module

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/$module/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/$module/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/$module/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// New${Module}Repository creates a new repository with injected DB
func New${Module}Repository(db *gorm.DB) *repository.${Module}Repository {
	return repository.New${Module}RepositoryWithDB(db)
}

// New${Module}Service creates a new service with injected dependencies
func New${Module}Service(repo *repository.${Module}Repository) *service.${Module}Service {
	return service.New${Module}ServiceWithDeps(repo)
}

// New${Module}Handler creates a new handler with injected service
func New${Module}Handler(svc *service.${Module}Service) *handler.${Module}Handler {
	return handler.New${Module}HandlerWithService(svc)
}

// Register${Module}Routes registers routes with DI
func Register${Module}Routes(rg *types.RouterGroup, h *handler.${Module}Handler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}

// Module provides all $module module dependencies
var Module = fx.Module("$module",
	fx.Provide(
		New${Module}Repository,
		New${Module}Service,
		New${Module}Handler,
	),
	fx.Invoke(Register${Module}Routes),
)
EOF
    
    return 0
}

# Main migration function
migrate_module() {
    local module=$1
    
    echo ""
    log_info "========================================="
    log_info "Migrating module: $module"
    log_info "========================================="
    
    add_repository_constructor "$module"
    add_service_constructor "$module"
    add_handler_constructor "$module"
    update_routes "$module"
    create_module_file "$module"
}

# List of modules to migrate (excluding already done ones)
MODULES_TO_MIGRATE=(
    "auth"
    "class"
    "curriculum"
    "dapodik"
    "exam"
    "school"
    "spp"
    "storage"
    "teacher"
)

echo ""
log_info "FX Migration Script"
log_info "==================="
echo ""

for module in "${MODULES_TO_MIGRATE[@]}"; do
    migrate_module "$module"
done

echo ""
log_info "Migration complete!"
log_info "Note: Some files may need manual adjustment."
log_info "Run 'go build ./...' to verify."

