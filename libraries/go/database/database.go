package database

import (
	"github.com/jinzhu/gorm"
)

// DefaultDB is a global instance of a gorm DB
var DefaultDB *gorm.DB

func mustGetDefaultDB() *gorm.DB {
	if DefaultDB == nil {
		panic("Database used before default DB set")
	}

	return DefaultDB
}

// Find finds records that match given conditions
func Find(out interface{}, where ...interface{}) *gorm.DB {
	return mustGetDefaultDB().Find(out, where...)
}

// Create inserts value into the database
func Create(value interface{}) *gorm.DB {
	return mustGetDefaultDB().Create(value)
}
