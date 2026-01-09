package database

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// =============================================================================
// Column Validation (SQL Injection Protection)
// =============================================================================

// validColumnPattern matches valid SQL column names (alphanumeric + underscore)
var validColumnPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// isValidColumn checks if a column name is safe for SQL queries
func isValidColumn(column string) bool {
	return validColumnPattern.MatchString(column)
}

// sanitizeColumn validates and returns safe column name, or empty string if invalid
func sanitizeColumn(column string) string {
	// Handle table.column format
	parts := strings.Split(column, ".")
	for _, part := range parts {
		if !isValidColumn(part) {
			log.Warn().Str("column", column).Msg("Invalid column name rejected")
			return ""
		}
	}
	return column
}

// =============================================================================
// Param & Paginator Structs
// =============================================================================

// Param holds pagination parameters
type Param struct {
	DB       *gorm.DB
	Page     int
	Limit    int
	OrderBy  []string
	Preloads []string
	ShowSQL  bool
}

// Paginator holds pagination result
type Paginator struct {
	TotalRecord int64       `json:"totalRecord"`
	TotalPage   int         `json:"totalPage"`
	Records     interface{} `json:"records"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prevPage"`
	NextPage    int         `json:"nextPage"`
}

// =============================================================================
// Core Paging Function
// =============================================================================

// Paging executes paginated query
func Paging(p *Param, result interface{}) *Paginator {
	db := p.DB

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100 // Max limit to prevent abuse
	}

	// Apply order by
	for _, o := range p.OrderBy {
		db = db.Order(o)
	}

	// Apply preloads
	for _, preload := range p.Preloads {
		db = db.Preload(preload)
	}

	// Parallel count for performance
	done := make(chan bool, 1)
	var count int64

	go countRecords(db, result, done, &count)

	// Calculate offset
	offset := (p.Page - 1) * p.Limit

	// Execute query
	if p.ShowSQL {
		db.Debug().Limit(p.Limit).Offset(offset).Find(result)
	} else {
		db.Limit(p.Limit).Offset(offset).Find(result)
	}

	// Wait for count to complete
	<-done

	// Build paginator
	totalPage := int(math.Ceil(float64(count) / float64(p.Limit)))
	if totalPage == 0 {
		totalPage = 1
	}

	prevPage := p.Page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	nextPage := p.Page + 1
	if nextPage > totalPage {
		nextPage = totalPage
	}

	return &Paginator{
		TotalRecord: count,
		TotalPage:   totalPage,
		Records:     result,
		Offset:      offset,
		Limit:       p.Limit,
		Page:        p.Page,
		PrevPage:    prevPage,
		NextPage:    nextPage,
	}
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}

// =============================================================================
// Fluent Pagination Builder
// =============================================================================

// Pagination builder for fluent API
type Pagination[T any] struct {
	db             *gorm.DB
	model          interface{}
	preloads       []string
	request        *pagination.Request[T]
	allowedColumns []string // Whitelist for filter/sort columns
	defaultSort    string   // Default sort field
}

// NewPagination creates a new pagination builder with injected DB
func NewPagination[T any](db *gorm.DB) *Pagination[T] {
	return &Pagination[T]{
		db:          db,
		defaultSort: "id",
	}
}

// SetModel sets the model for pagination
func (p *Pagination[T]) SetModel(model interface{}) *Pagination[T] {
	p.model = model
	return p
}

// SetPreloads sets preload relations
func (p *Pagination[T]) SetPreloads(preloads ...string) *Pagination[T] {
	p.preloads = preloads
	return p
}

// SetRequest sets pagination request
func (p *Pagination[T]) SetRequest(request *pagination.Request[T]) *Pagination[T] {
	p.request = request
	return p
}

// SetAllowedColumns sets whitelist of allowed columns for filtering/sorting
// If not set, all valid column names are allowed (basic SQL injection protection still applies)
func (p *Pagination[T]) SetAllowedColumns(columns ...string) *Pagination[T] {
	p.allowedColumns = columns
	return p
}

// SetDefaultSort sets the default sort field when not specified in request
func (p *Pagination[T]) SetDefaultSort(field string) *Pagination[T] {
	p.defaultSort = field
	return p
}

// isColumnAllowed checks if column is in whitelist (if whitelist is set)
func (p *Pagination[T]) isColumnAllowed(column string) bool {
	// Basic sanitization first
	if sanitizeColumn(column) == "" {
		return false
	}

	// If no whitelist, allow all sanitized columns
	if len(p.allowedColumns) == 0 {
		return true
	}

	// Check whitelist
	for _, allowed := range p.allowedColumns {
		if column == allowed {
			return true
		}
	}

	log.Warn().Str("column", column).Msg("Column not in whitelist, rejected")
	return false
}

// FindAllPaging executes paginated query with filters
func (p *Pagination[T]) FindAllPaging() *Paginator {
	query := p.db

	pageRequest := p.request
	if pageRequest == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Please provide Paging Request"))
	}

	// Apply search filter (if any)
	if pageRequest.Search.Key != "" && pageRequest.Search.Value != "" {
		searchKey := sanitizeColumn(pageRequest.Search.Key)
		if searchKey != "" && p.isColumnAllowed(searchKey) {
			query = query.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", searchKey), "%"+pageRequest.Search.Value+"%")
		}
	}

	// Apply filter (if any)
	if pageRequest.Filter != nil {
		if filterMap, ok := any(*pageRequest.Filter).(map[string]interface{}); ok {
			for key, value := range filterMap {
				if value == nil {
					continue
				}

				// Validate column name
				safeKey := sanitizeColumn(key)
				if safeKey == "" || !p.isColumnAllowed(safeKey) {
					continue
				}

				switch v := value.(type) {
				case string:
					if v == "" {
						continue
					}
					query = query.Where(fmt.Sprintf("%s = ?", safeKey), v)

				case float64:
					query = query.Where(fmt.Sprintf("%s = ?", safeKey), v)

				case int, int64, uint, uint64:
					query = query.Where(fmt.Sprintf("%s = ?", safeKey), v)

				case bool:
					query = query.Where(fmt.Sprintf("%s = ?", safeKey), v)

				case []interface{}:
					if len(v) > 0 {
						query = query.Where(fmt.Sprintf("%s IN ?", safeKey), v)
					}

				default:
					log.Warn().Str("key", key).Interface("value", value).Type("type", value).Msg("Unhandled filter type")
				}
			}
		}
	}

	// Build order by from request or default
	sort := pageRequest.DefaultSort(p.defaultSort)
	orderBy := []string{}
	if p.isColumnAllowed(sort.Field) {
		orderBy = append(orderBy, fmt.Sprintf("%s %s", sort.Field, sort.Dir))
	} else {
		// Fallback to default
		orderBy = append(orderBy, fmt.Sprintf("%s desc", p.defaultSort))
	}

	// Determine if we should show SQL (based on environment)
	showSQL := config.GetApp().Env == "development"

	pg := Paging(&Param{
		DB:       query,
		Page:     pageRequest.Page,
		Limit:    pageRequest.Size,
		OrderBy:  orderBy,
		Preloads: p.preloads,
		ShowSQL:  showSQL,
	}, &p.model)

	return pg
}

// =============================================================================
// Deprecated Functions (for backward compatibility)
// =============================================================================

// SetModal is deprecated, use SetModel instead
// Deprecated: This function has a typo, use SetModel instead
func (p *Pagination[T]) SetModal(model interface{}) *Pagination[T] {
	log.Warn().Msg("SetModal is deprecated, use SetModel instead")
	return p.SetModel(model)
}
