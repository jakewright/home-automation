package router

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/taxi"
)

var (
	defaultRouter *Router
	once          = &sync.Once{}
)

// SetDefaultRouter sets the global Router instance
func SetDefaultRouter(r *Router) {
	once.Do(func() { defaultRouter = r })
}

func mustGetDefaultRouter() *Router {
	if defaultRouter == nil {
		panic(fmt.Errorf("no default router set"))
	}

	return defaultRouter
}

// RegisterHandler registers a Taxi handler on the default router
func RegisterHandler(method, path string, handler taxi.HandlerFunc) {
	mustGetDefaultRouter().RegisterHandler(method, path, handler)
}

// Router sets up a Taxi server
type Router struct {
	port            int
	router          *taxi.Router
	server          *http.Server
	shutdownInvoked *int32
}

// New returns a new router initialised with default middleware
func New() *Router {
	var conf struct {
		Port int `envconfig:"default=80"`
	}

	config.Load(&conf)

	router := taxi.NewRouter().WithLogger(slog.Errorf)
	router.UseMiddleware(panicRecovery, revision)
	router.RegisterHandlerFunc(http.MethodGet, "/ping", pingHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	r := &Router{
		port:            conf.Port,
		router:          router,
		server:          server,
		shutdownInvoked: new(int32),
	}

	return r
}

// GetName returns a friendly name for the process
func (r *Router) GetName() string {
	return "router"
}

// Start will listen for TCP connections on the port defined in config
func (r *Router) Start() error {
	slog.Infof("Listening on port %d", r.port)
	err := r.server.ListenAndServe()

	// This error will always be returned after Shutdown is called so swallow it here
	if atomic.LoadInt32(r.shutdownInvoked) > 0 && err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Stop will gracefully shutdown the server
func (r *Router) Stop(ctx context.Context) error {
	atomic.StoreInt32(r.shutdownInvoked, 1)
	return r.server.Shutdown(ctx)
}

// RegisterHandler adds a route to the router
func (r *Router) RegisterHandler(method, path string, handler taxi.HandlerFunc) {
	r.router.RegisterHandler(method, path, handler)
}
