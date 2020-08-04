package database

import (
	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// Gorm is a database accessor that uses a gorm DB instance
type Gorm struct {
	db *gorm.DB
}

var _ Database = (*Gorm)(nil)

// NewGorm returns a new database using the specifed gorm instance
func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{
		db: db,
	}
}

// Find marshals all matching records into out
func (g *Gorm) Find(out interface{}, where ...interface{}) error {
	if err := g.db.Find(out, where...).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return oops.WithCode(err, oops.ErrNotFound)
		}

		return oops.Wrap(err, oops.ErrInternalService, "failed to execute find")
	}

	return nil
}

// Create adds a new record based on value
func (g *Gorm) Create(value interface{}) error {
	if err := g.db.Create(value).Error; err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to execute create")
	}
	return nil
}

// Delete performs a hard-delete of matching rows
func (g *Gorm) Delete(value interface{}, where ...interface{}) error {
	// Unscoped() disables soft delete
	if err := g.db.Unscoped().Delete(value, where...).Error; err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to execute delete")
	}
	return nil
}
