package database

// Database is an interface for accessing the system's persistent data store
type Database interface {
	Find(out interface{}, where ...interface{}) error
	Create(value interface{}) error
	Delete(value interface{}, where ...interface{}) error
}
