package database

import "gorm.io/gorm"

// GpaRepository is a generic repository for common CRUD operations
type GpaRepository[T interface{}] struct {
	model T
	db    *gorm.DB
}

// NewGpaRepository creates a new generic repository with injected DB
func NewGpaRepository[T any](model T, db *gorm.DB) *GpaRepository[T] {
	return &GpaRepository[T]{
		model: model,
		db:    db,
	}
}

// FindById finds entity by ID
func (g *GpaRepository[T]) FindById(id uint) *T {
	var result T
	g.db.Where("id = ?", id).First(&result)
	return &result
}

// FindByIds finds entities by multiple IDs
func (g *GpaRepository[T]) FindByIds(ids []uint) []T {
	var results []T
	g.db.Where("id IN (?)", ids).Find(&results)
	return results
}

// DeleteById deletes entity by ID
func (g *GpaRepository[T]) DeleteById(id uint) error {
	return g.db.Where("id = ?", id).Delete(&g.model).Error
}
