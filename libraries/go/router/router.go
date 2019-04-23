package router

import (
	"context"
	"fmt"
	"home-automation/libraries/go/config"
	"home-automation/libraries/go/slog"
	"net/http"

	"github.com/jakewright/muxinator"
)

var router muxinator.Router

func init() {
	router = muxinator.NewRouter()
}

// ListenAndServe listens for TCP connections on the port defined in config
func ListenAndServe() error {
	//return errors.InternalService("testing something")
	port := config.Get("port").Int(80)
	slog.Info("Listening on port %d", port)
	return router.ListenAndServe(fmt.Sprintf(":%d", port))
}

// Shutdown gracefully shuts down the server
func Shutdown(ctx context.Context) error {
	return router.Shutdown(ctx)
}

// Get is a helper function to add a GET route
func Get(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	router.Get(path, handler, middlewares...)
}

// Post is a helper function to add a POST route
func Post(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	router.Post(path, handler, middlewares...)
}

// Put is a helper function to add a PUT route
func Put(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	router.Put(path, handler, middlewares...)
}

// Patch is a helper function to add a PATCH route
func Patch(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	router.Patch(path, handler, middlewares...)
}

// Delete is a helper function to add a DELETE route
func Delete(path string, handler http.HandlerFunc, middlewares ...muxinator.Middleware) {
	router.Delete(path, handler, middlewares...)
}
