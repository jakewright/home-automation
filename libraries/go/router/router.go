package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/taxi"
)

// Router handles HTTP requests to the service
type Router struct {
	port   int
	router *taxi.Router
	server *http.Server
}

// Service represents the entire application and is used by
// the router to set up various standard endpoints that all
// services are given, such as /ping and /healthz.
type Service interface {
	Hostname() string
	Revision() string
	HealthProvider
}

// New returns a new router initialised with default middleware
func New(svc Service) *Router {
	var conf struct {
		Port int `envconfig:"default=80"`
	}
	config.Load(&conf)

	// Create the router
	router := taxi.NewRouter().WithLogger(slog.Errorf)
	router.UseMiddleware(panicRecovery, revision(svc.Revision()))
	router.HandleFunc(http.MethodGet, "/ping", PingHandler)

	router.HandleRaw(http.MethodGet, "/healthz", &HealthHandler{
		Hostname: svc.Hostname(),
		Revision: svc.Revision(),
		Provider: svc,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	r := &Router{
		port:   conf.Port,
		router: router,
		server: server,
	}

	return r
}

// GetName returns a friendly name for the process
func (r *Router) GetName() string {
	return "router"
}

// Start will listen for TCP connections on the port defined in config
func (r *Router) Start(ctx context.Context) error {
	slog.Infof("Listening on port %d", r.port)

	ch := make(chan error)

	go func() {
		ch <- r.server.ListenAndServe()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		if err := r.server.Shutdown(context.Background()); err != nil {
			return err
		}
	}

	err := <-ch

	// This error will always be returned after Shutdown is called so swallow it here
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// HandleFunc adds a route to the router with a taxi handler func
func (r *Router) HandleFunc(method, path string, handler func(context.Context, taxi.Decoder) (interface{}, error)) {
	r.router.HandleFunc(method, path, handler)
}

// HandleRaw adds a route to the router with an http.Handler
func (r *Router) HandleRaw(method, path string, handler http.Handler) {
	r.router.HandleRaw(method, path, handler)
}
