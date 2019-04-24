package router

import (
	"context"
	"fmt"
	"home-automation/libraries/go/config"
	"home-automation/libraries/go/slog"
	"net/http"
	"sync/atomic"

	"github.com/jakewright/muxinator"
)

// Router is a wrapper around muxinator.Router that conforms to the bootstrap.Process interface
type Router struct {
	r               muxinator.Router
	shutdownInvoked *int32
}

// New returns a new router
func New() *Router {
	return &Router{
		r:               muxinator.NewRouter(),
		shutdownInvoked: new(int32),
	}
}

func (r *Router) GetName() string {
	return "router"
}

// Start will listen for TCP connections on the port defined in config
func (r *Router) Start() error {
	port := config.Get("port").Int(80)
	slog.Info("Listening on port %d", port)
	err := r.r.ListenAndServe(fmt.Sprintf(":%d", port))

	// This error will always be returned after Shutdown is called so swallow it here
	if atomic.LoadInt32(r.shutdownInvoked) > 0 && err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Stop will gracefully shutdown the server
func (r *Router) Stop(ctx context.Context) error {
	atomic.StoreInt32(r.shutdownInvoked, 1)
	return r.r.Shutdown(ctx)
}

// Get is a helper function to add a GET route
func (r *Router) Get(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	r.r.Get(path, handler, middlewares...)
}

// Post is a helper function to add a POST route
func (r *Router) Post(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	r.r.Post(path, handler, middlewares...)
}

// Put is a helper function to add a PUT route
func (r *Router) Put(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	r.r.Put(path, handler, middlewares...)
}

// Patch is a helper function to add a PATCH route
func (r *Router) Patch(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	r.r.Patch(path, handler, middlewares...)
}

// Delete is a helper function to add a DELETE route
func (r *Router) Delete(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	r.r.Delete(path, handler, middlewares...)
}
