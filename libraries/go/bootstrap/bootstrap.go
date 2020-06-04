package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/dsync"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/rpc"
	"github.com/jakewright/home-automation/libraries/go/slog"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"

	// Register MySQL driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Process is a long-running task that provides service functionality
type Process interface {
	// GetName returns a friendly name for the process for use in logs
	GetName() string

	// Start kicks off the task and only returns when the task has finished
	Start() error

	// Stop will try to gracefully end the task and should be safe to run regardless of whether the process is currently running
	Stop(context.Context) error
}

// Service represents a collection of processes
type Service struct {
	processes []Process
	deferred  []func() error
}

// Opts defines basic initialisation options for a service
type Opts struct {
	// Config is a pointer to a struct which, if not nil, will
	// be populated with config from environment variables.
	Config interface{}

	// ServiceName is the name of the service e.g. service.foo
	ServiceName string

	// Firehose indicates whether a connection to Redis should be made
	Firehose bool

	// Database indicates whether a connection to MySQL should be made
	Database bool
}

// Init performs standard service startup tasks and returns a Service
func Init(opts *Opts) (*Service, error) {
	service := &Service{}

	// Load config if requested
	if opts.Config != nil {
		config.Load(opts.Config)
	}

	// Create default API client
	apiClient, err := rpc.NewHTTPClient("data")
	if err != nil {
		return nil, err
	}
	rpc.SetDefaultClient(apiClient)

	// Connect to Redis
	if opts.Firehose {
		if err := initFirehose(service); err != nil {
			return nil, err
		}
	}

	// Connect to MySQL
	if opts.Database {
		if err := initDatabase(opts, service); err != nil {
			return nil, err
		}
	}

	// Set up locking
	dsync.DefaultLocksmith = dsync.NewLocalLocksmith()

	return service, nil
}

func initFirehose(svc *Service) error {
	conf := struct {
		RedisHost string
		RedisPort int
	}{}
	config.Load(&conf)

	addr := fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort)
	slog.Infof("Connecting to Redis at address %s", addr)
	redisClient := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        "",
		DB:              0,
		MaxRetries:      5,
		MinRetryBackoff: time.Second,
		MaxRetryBackoff: time.Second * 5,
	})

	svc.deferred = append(svc.deferred, func() error {
		err := redisClient.Close()
		if err != nil {
			slog.Errorf("Failed to close Redis connection: %v", err)
		} else {
			slog.Debugf("Closed Redis connection")
		}
		return err
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	firehoseClient := firehose.NewRedisClient(redisClient)
	svc.processes = append(svc.processes, firehoseClient)

	firehose.DefaultClient = firehoseClient

	return nil
}

func initDatabase(opts *Opts, svc *Service) error {
	conf := struct {
		MySQLHost         string
		MySQLUsername     string
		MySQLPassword     string
		MySQLDatabaseName string `envconfig:"default=home_automation"`
		MySQLCharset      string `envconfig:"default=utf8mb4"`
	}{}
	config.Load(&conf)

	// Replace hyphens and dots in the service name with underscores
	re, err := regexp.Compile(`[-.]`)
	if err != nil {
		return err
	}
	prefix := re.ReplaceAllString(opts.ServiceName, "_")

	// Remove any remaining non alphanumeric characters
	re, err = regexp.Compile(`[^a-zA-Z0-9_]+`)
	if err != nil {
		return err
	}
	prefix = re.ReplaceAllString(prefix, "")

	// Set a default table prefix
	gorm.DefaultTableNameHandler = func(_ *gorm.DB, defaultTableName string) string {
		return prefix + "_" + defaultTableName
	}

	if conf.MySQLHost == "" || conf.MySQLUsername == "" || conf.MySQLPassword == "" {
		return oops.InternalService("MySQL host, username and password not set in config")
	}

	addr := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.MySQLUsername,
		conf.MySQLPassword,
		conf.MySQLHost,
		conf.MySQLDatabaseName,
		conf.MySQLCharset)

	db, err := gorm.Open("mysql", addr)
	if err != nil {
		return err
	}

	// Always load associations
	db.InstantSet("gorm:auto_preload", true)

	svc.deferred = append(svc.deferred, func() error {
		err := db.Close()
		if err != nil {
			slog.Errorf("Failed to close MySQL connection: %v", err)
		} else {
			slog.Debugf("Closed MySQL connection")
		}
		return err
	})

	database.DefaultDB = db
	return nil
}

// Run takes a number of processes and concurrently runs them all. It will stop if all processes
// terminate or if a signal (SIGINT or SIGTERM) is received.
func (s *Service) Run(processes ...Process) {
	// os.Exit should be the last thing to happen
	var code int
	defer os.Exit(code)

	// Close all of the resources after processes have shut down
	for _, deferred := range s.deferred {
		defer func(d func() error) {
			if err := d(); err != nil {
				code = 1
			}
		}(deferred)
	}

	s.processes = append(s.processes, processes...)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	// Start all of the processes in goroutines
	for _, process := range s.processes {
		process := process

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := process.Start(); err != nil {
				slog.Errorf("Process %s stopped with error: %v", process.GetName(), err)
				code = 1
			} else {
				slog.Debugf("Process %s stopped", process.GetName())
			}
		}()
	}

	// Close the done channel when all processes return
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all processes to return or for a signal
	select {
	case <-done:
		slog.Warnf("All processes stopped prematurely")
		return
	case s := <-sig:
		slog.Infof("Received %v signal", s)
	}

	// A short timeout because Docker will kill us after 10 seconds anyway
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simultaneously stop all processes
	for _, process := range s.processes {
		process := process

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := process.Stop(ctx); err != nil {
				slog.Errorf("Failed to stop %s gracefully: %v", process.GetName(), err)
				code = 1
			}
		}()
	}

	// Wait for processes to terminate
	wg.Wait()
	slog.Infof("All processes stopped")
}
