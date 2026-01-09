package database

import "gorm.io/gorm"

// SchoolScope returns a gorm scope function that filters by school_code
// Use this in repository methods to ensure multi-tenant data isolation
//
// Example usage:
//
//	repo.Database.Scopes(SchoolScope(schoolCode)).Find(&data)
func SchoolScope(schoolCode string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if schoolCode == "" {
			return db
		}
		return db.Where("school_code = ?", schoolCode)
	}
}
