package bootstrap

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danielchatfield/go-randutils"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/distsync"
	"github.com/jakewright/home-automation/libraries/go/firehose"
	"github.com/jakewright/home-automation/libraries/go/healthz"
	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/taxi"
)

// Revision is the service's revision and should be
// set at build time to the current git commit hash.
var Revision string

// Service represents a collection of processes
type Service struct {
	name     string
	hostname string
	revision string
	id       string // Used to identify the instance in a Firehose consumer group
	router   *router.Router
	runner   *runner

	// Long-lived connections shared across the whole application.
	// Do not access these variables directly. Use the getX()
	// functions which will initialise them if necessary.
	mysqlCon       *gorm.DB
	redisClient    *redis.Client
	firehoseClient *firehose.Client
}

// Opts defines basic initialisation options for a service
type Opts struct {
	// Config is a pointer to a struct which, if not nil, will
	// be populated with config from environment variables.
	Config interface{}

	// ServiceName is the name of the service e.g. service.foo
	ServiceName string
}

// Init performs standard service startup tasks and returns a Service
func Init(opts *Opts) *Service {
	svc, err := initService(opts)
	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}
	return svc
}

// Hostname returns the hostname as reported by the
// operating system's kernel.
func (s *Service) Hostname() string {
	return s.hostname
}

// Revision returns the current revision of the code.
func (s *Service) Revision() string {
	return s.revision
}

// Health returns a map representing the result of all of
// the health checks that have been registered.
func (s *Service) Health(ctx context.Context) map[string]error {
	return healthz.Status(ctx)
}

// Database is a helper function that returns a cached database.
// The first time it is called, a new connection to the database
// is established. Closing the connection when the program ends
// is handled automatically.
func (s *Service) Database() database.Database {
	db, err := s.getMySQL()
	if err != nil {
		panic(err)
	}

	return database.NewGorm(db)
}

// FirehosePublisher is a helper function that returns a cached
// firehose client. The first time it is called, a new firehose
// client is set up.
func (s *Service) FirehosePublisher() firehose.Publisher {
	client, err := s.getFirehoseClient()
	if err != nil {
		panic(err)
	}

	return client
}

// HandleFunc registers a new taxi-style handler with the
// application's router for the specified method and path.
func (s *Service) HandleFunc(method, path string, handler func(context.Context, taxi.Decoder) (interface{}, error)) {
	s.router.HandleFunc(method, path, handler)
}

// HandleRaw registers a new http.Handler with the application's
// router for the specified method and path.
func (s *Service) HandleRaw(method, path string, handler http.Handler) {
	s.router.HandleRaw(method, path, handler)
}

// Run starts all processes that have already been registered
// with the service, plus any extra ones passed in as arguments.
// This function blocks until either all processes return, or
// an interrupt signal is received, at which point all processes
// are signalled to end. If any processes do not, at this point,
// return, then Run() may hang indefinitely.
func (s *Service) Run(processes ...Process) {
	for _, process := range processes {
		s.runner.addProcess(process)
	}

	s.runner.Run()
}

func initService(opts *Opts) (*Service, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	id, err := randutils.String(5)
	if err != nil {
		return nil, err
	}

	service := &Service{
		name:     opts.ServiceName,
		hostname: hostname,
		revision: Revision,
		id:       id,
		runner:   &runner{},
	}

	// Load config if requested
	if opts.Config != nil {
		config.Load(opts.Config)
	}

	// Set up locking
	if err := initLock(opts, service); err != nil {
		return nil, err
	}

	// Set up router
	initRouter(service)

	return service, nil
}

func initLock(opts *Opts, svc *Service) error {
	conf := struct {
		LockMode    string        `envconfig:"LOCK_MODE"`
		LockTimeout time.Duration `envconfig:"optional,LOCK_TIMEOUT"`
		LockTTL     time.Duration `envconfig:"optional,LOCK_TTL"`
	}{}
	config.Load(&conf)

	switch strings.ToLower(conf.LockMode) {
	case "local":
		distsync.DefaultLocksmith = distsync.NewLocalLocksmith()

	case "shared":
		redisClient, err := svc.getRedisClient()
		if err != nil {
			return err
		}

		distsync.DefaultLocksmith = &distsync.RedisLocksmith{
			ServiceName: opts.ServiceName,
			Client:      redisClient,
			Timeout:     conf.LockTimeout,
			Expiration:  conf.LockTTL,
		}

	default:
		return oops.InternalService("unknown lock mode %q", conf.LockMode)
	}

	return nil
}

func initRouter(svc *Service) {
	svc.router = router.New(svc)
	svc.runner.addProcess(svc.router)
}

func (s *Service) getFirehoseClient() (*firehose.Client, error) {
	if s.firehoseClient == nil {
		redisClient, err := s.getRedisClient()
		if err != nil {
			return nil, err
		}

		s.firehoseClient = firehose.NewClient(&firehose.ClientOptions{
			Group:          s.name,
			Consumer:       s.id,
			Redis:          redisClient,
			HandlerTimeout: 0, // TODO: let services override this
		})
		s.runner.addProcess(s.firehoseClient)
	}

	return s.firehoseClient, nil
}
