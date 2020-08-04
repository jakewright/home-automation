package bootstrap

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/healthz"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"

	// Register MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type mySQLConfig struct {
	Host         string `envconfig:"MYSQL_HOST"`
	Username     string `envconfig:"MYSQL_USERNAME"`
	Password     string `envconfig:"MYSQL_PASSWORD"`
	DatabaseName string `envconfig:"default=home_automation"`
	Charset      string `envconfig:"default=utf8mb4"`
}

func (c *mySQLConfig) load() error {
	config.Load(&c)

	switch {
	case c.Host == "":
		return oops.InternalService("MYSQL_HOST not set")
	case c.Username == "":
		return oops.InternalService("MYSQL_USERNAME not set")
	case c.Password == "":
		return oops.InternalService("MYSQL_PASSWORD not set")
	}

	return nil
}

// getMySQL returns a cached instance of a gorm.DB connected
// to a MySQL server. If this is being called for the first
// time, a new connection to MySQL is made. Connection options
// are read from config.
func (s *Service) getMySQL() (*gorm.DB, error) {
	if s.mysqlCon == nil {
		conf := mySQLConfig{}
		if err := conf.load(); err != nil {
			return nil, err
		}

		prefix := tableNamePrefix(s.name)

		// Set a default table prefix
		gorm.DefaultTableNameHandler = func(_ *gorm.DB, defaultTableName string) string {
			return prefix + "_" + defaultTableName
		}

		addr := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.DatabaseName,
			conf.Charset,
		)

		db, err := gorm.Open("mysql", addr)
		if err != nil {
			return nil, err
		}

		s.runner.addDeferred(func() error {
			err := db.Close()
			if err != nil {
				slog.Errorf("Failed to close MySQL connection: %v", err)
			} else {
				slog.Debugf("Closed MySQL connection")
			}
			return err
		})

		// Always load associations
		db.InstantSet("gorm:auto_preload", true)

		healthz.RegisterCheck("mysql", func(ctx context.Context) error {
			return db.DB().PingContext(ctx)
		})

		s.mysqlCon = db
	}

	return s.mysqlCon, nil
}

func tableNamePrefix(serviceName string) string {
	// Replace non alphanumeric characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	prefix := re.ReplaceAllString(serviceName, "_")
	return strings.ToLower(prefix)
}
