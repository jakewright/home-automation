package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

const (
	defaultSchemaFilename   = "schema.sql"
	defaultMockDataFilename = "mock_data.sql"
)

// GetDefaultSchema returns the schema found at
// schema/schema.sql in the service's directory.
func GetDefaultSchema(serviceName string) (string, error) {
	filename := fmt.Sprintf("./%s/schema/%s", serviceName, defaultSchemaFilename)
	return readFileIfExists(filename)
}

// GetMockSQL returns the schema found at schema/mock_data.sql
//in the service's directory.
func GetMockSQL(serviceName string) (string, error) {
	filename := fmt.Sprintf("./%s/schema/%s", serviceName, defaultMockDataFilename)
	return readFileIfExists(filename)
}

func readFileIfExists(filename string) (string, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to stat %s: %w", filename, err)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", filename, err)
	}

	return string(b), nil
}

// ApplySchema applies the schema to the database.
func ApplySchema(db *config.Database, serviceName, schema string) error {
	switch db.Engine {
	case "mysql":
		return applyMySQLSchema(db, serviceName, schema)
	default:
		return fmt.Errorf("unsupported database engine %q", db.Engine)
	}
}

func applyMySQLSchema(db *config.Database, serviceName, schema string) error {
	running, err := service.IsRunning(db.Service)
	if err != nil {
		return fmt.Errorf("failed to get status of %s: %v", db.Service, err)
	}

	if !running {
		if err := service.Run([]string{db.Service}); err != nil {
			return fmt.Errorf("failed to start %s: %v", db.Service, err)
		}

		op := output.Info("Waiting for database to startup")
		time.Sleep(time.Second * 5)
		op.Complete()
	}

	op := output.Info("Applying schema for %s", serviceName)
	if err := service.Exec(db.Service, schema, "mysql", "-u"+db.Username, "-p"+db.Password); err != nil {
		op.Failed()
		return err
	}

	op.Complete()
	return nil
}
