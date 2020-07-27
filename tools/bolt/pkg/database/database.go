package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jakewright/home-automation/tools/bolt/pkg/compose"
	"github.com/jakewright/home-automation/tools/bolt/pkg/config"
	"github.com/jakewright/home-automation/tools/bolt/pkg/service"
	"github.com/jakewright/home-automation/tools/deploy/pkg/output"
)

const (
	defaultSchemaFilename   = "schema.sql"
	defaultMockDataFilename = "mock_data.sql"
)

// Database performs operations on a database container
type Database struct {
	c   *compose.Compose
	cfg *config.Database
}

// New returns a Database
func New(c *compose.Compose, config *config.Database) *Database {
	return &Database{
		c:   c,
		cfg: config,
	}
}

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
func (d *Database) ApplySchema(serviceName, schema string) error {
	switch d.cfg.Engine {
	case "mysql":
		return d.applyMySQLSchema(serviceName, schema)
	default:
		return fmt.Errorf("unsupported database engine %q", d.cfg.Engine)
	}
}

func (d *Database) applyMySQLSchema(serviceName, schema string) error {
	running, err := d.c.IsRunning(d.cfg.Service)
	if err != nil {
		return fmt.Errorf("failed to get status of %s: %v", d.cfg.Service, err)
	}

	if !running {
		if err := service.Run(d.c, []string{d.cfg.Service}); err != nil {
			return fmt.Errorf("failed to start %s: %v", d.cfg.Service, err)
		}

		op := output.Info("Waiting for database to startup")
		time.Sleep(time.Second * 5)
		op.Success()
	}

	op := output.Info("Applying schema for %s", serviceName)
	if err := d.c.Exec(d.cfg.Service, schema, "mysql", "-u"+d.cfg.Username, "-p"+d.cfg.Password); err != nil {
		op.Failed()
		return err
	}

	op.Success()
	return nil
}
