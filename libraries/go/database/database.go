package database

import (
	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

// DefaultDB is a global instance of a gorm DB
var DefaultDB *gorm.DB

func mustGetDefaultDB() *gorm.DB {
	if DefaultDB == nil {
		slog.Panicf("Database used before default DB set")
	}

	return DefaultDB
}

// Find finds records that match given conditions
func Find(out interface{}, where ...interface{}) error {
	err := mustGetDefaultDB().Find(out, where...).Error
	return errors.WithMessage(err, "failed to execute find")
}

// Create inserts value into the database
func Create(value interface{}) error {
	err := mustGetDefaultDB().Create(value).Error
	return errors.WithMessage(err, "failed to execute create")
}

// Delete deletes a value from the database
func Delete(value interface{}, where ...interface{}) error {
	err := mustGetDefaultDB().Delete(value, where...).Error
	return errors.WithMessage(err, "failed to execute delete")
}
