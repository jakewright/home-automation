package taxi

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

// Middleware is a function that takes a handler and returns a new handler
type Middleware func(http.Handler) http.Handler

// Router registers routes and handlers to handle RPCs over HTTP.
type Router struct {
	router     *httprouter.Router
	middleware []Middleware
	logFunc    func(format string, v ...interface{})
}

// NewRouter returns an initialised Router
func NewRouter() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

// WithLogger sets a log function for the router to use when something goes.
// wrong. If not set, no logs will be output.
func (r *Router) WithLogger(f func(format string, v ...interface{})) *Router {
	r.logFunc = f
	return r
}

// Handle registers a new route
func (r *Router) Handle(method, path string, handler Handler) {
	r.HandleRaw(method, path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

// HandleFunc registers a new route
func (r *Router) HandleFunc(method, path string, handler func(context.Context, Decoder) (interface{}, error)) {
	r.Handle(method, path, HandlerFunc(handler))
}

// HandleRaw registers a new route with a standard http.Handler
func (r *Router) HandleRaw(method, path string, handler http.Handler) {
	r.router.Handler(method, path, handler)
}

// UseMiddleware adds the given middleware to the router. The middleware
// functions are executed in the order given.
func (r *Router) UseMiddleware(mw ...Middleware) {
	for _, m := range mw {
		r.middleware = append(r.middleware, m)
	}
}

// ServeHTTP dispatches requests to the appropriate handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Wrap the handler in the middleware functions
	var handler http.Handler = r.router
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	handler.ServeHTTP(w, req)
}

func (r *Router) log(format string, v ...interface{}) {
	if r.logFunc == nil {
		return
	}

	r.logFunc(format, v...)
}
