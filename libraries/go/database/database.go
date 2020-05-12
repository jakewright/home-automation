package database

import (
	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/oops"
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
	if err := mustGetDefaultDB().Find(out, where...).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return oops.WithCode(err, oops.ErrNotFound)
		}

		return oops.Wrap(err, oops.ErrInternalService, "failed to execute find")
	}

	return nil
}

// Create inserts value into the database
func Create(value interface{}) error {
	if err := mustGetDefaultDB().Create(value).Error; err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to execute create")
	}
	return nil
}

// Delete deletes a value from the database
func Delete(value interface{}, where ...interface{}) error {
	// Unscoped() disables soft delete
	if err := mustGetDefaultDB().Unscoped().Delete(value, where...).Error; err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to execute delete")
	}
	return nil
}
