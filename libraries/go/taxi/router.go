package taxi

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// Decoder is a function that decodes a request body into the given interface
type Decoder func(v interface{}) error

// Handler is an interface that wraps the ServeRPC method
type Handler interface {
	ServeRPC(ctx context.Context, decode Decoder) (interface{}, error)
}

// HandlerFunc is a type that allows normal functions to be used as Handlers
type HandlerFunc func(ctx context.Context, decode Decoder) (interface{}, error)

// ServeRPC calls f(ctx, decode)
func (f HandlerFunc) ServeRPC(ctx context.Context, decode Decoder) (interface{}, error) {
	return f(ctx, decode)
}

// Router registers routes and handlers to handle RPCs over HTTP.
type Router struct {
	router  *mux.Router
	logFunc func(format string, v ...interface{})
}

// NewRouter returns an initialised Router
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// WithLogger sets a log function for the router to use when something goes.
// wrong. If not set, no logs will be output.
func (r *Router) WithLogger(f func(format string, v ...interface{})) *Router {
	r.logFunc = f
	return r
}

// RegisterHandler registers a new route
func (r *Router) RegisterHandler(method, path string, handler Handler) {
	r.RegisterRawHandler(method, path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		decoder := func(v interface{}) error {
			return DecodeRequest(req, v)
		}

		rsp, err := handler.ServeRPC(req.Context(), decoder)
		if err != nil {
			r.log("Failed to handle request: %v", err)
			if err := WriteError(w, err); err != nil {
				r.log("Failed to write response: %v", err)
			}
			return
		}

		if err := WriteSuccess(w, rsp); err != nil {
			r.log("Failed to handle request: %v", err)
		}
	}))
}

// RegisterHandlerFunc registers a new route
func (r *Router) RegisterHandlerFunc(method, path string, handler func(context.Context, Decoder) (interface{}, error)) {
	r.RegisterHandler(method, path, HandlerFunc(handler))
}

// RegisterRawHandler registers a new route with a standard http.Handler
func (r *Router) RegisterRawHandler(method, path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(method)
}

// UseMiddleware adds a stack of middleware to the router
func (r *Router) UseMiddleware(mw ...func(http.Handler) http.Handler) {
	mws := make([]mux.MiddlewareFunc, len(mw))
	for i, mw := range mw {
		mws[i] = mw
	}
	r.router.Use(mws...)
}

// ServeHTTP dispatches requests to the appropriate handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) log(format string, v ...interface{}) {
	if r.logFunc == nil {
		return
	}

	r.logFunc(format, v...)
}
