package taxi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

const (
	contentTypeJSON = "application/json; charset=UTF-8"
	contentTypeText = "text/plain"
)

// Decoder is a function that decodes a request body into the given interface
type Decoder func(v interface{}) error

// Handler is an interface that wraps the ServeRPC method
type Handler interface {
	ServeRPC(ctx context.Context, decode Decoder) (interface{}, error)
}

// HandlerFunc is a type that allows normal functions to be used as Handlers
type HandlerFunc func(context.Context, Decoder) (interface{}, error)

// ServeRPC calls f(ctx, decode)
func (f HandlerFunc) ServeRPC(ctx context.Context, decode Decoder) (interface{}, error) {
	return f(ctx, decode)
}

// Router registers routes and handlers to handle RPCs over HTTP.
type Router struct {
	router  *mux.Router
	rw      *responseWriter
	logFunc func(format string, v ...interface{})
}

// NewRouter returns an initialised Router
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
		rw:     &responseWriter{},
	}
}

// WithLogger sets a log function for the router to use when something goes.
// wrong. If not set, no logs will be output.
func (r *Router) WithLogger(f func(format string, v ...interface{})) *Router {
	r.logFunc = f
	r.rw.logFunc = f
	return r
}

// RegisterHandler registers a new route
func (r *Router) RegisterHandler(method, path string, handler Handler) {
	var fn http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		decoder := func(v interface{}) error {
			return decodeRequest(req, v)
		}

		rsp, err := handler.ServeRPC(req.Context(), decoder)
		if err != nil {
			r.log("Failed to handle request: %v", err)
			r.rw.writeResponse(w, err)
			return
		}

		r.rw.writeResponse(w, rsp)
	}

	r.RegisterRawHandler(method, path, fn)
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
func (r *Router) UseMiddleware(mw ...mux.MiddlewareFunc) {
	r.router.Use(mw...)
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

// decodeRequest unmarshals URL parameters and the JSON body
// of the given request into the value pointed to by v.
func decodeRequest(r *http.Request, v interface{}) error {
	// This does a load of reflection to unmarshal a map into the type of v
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true,
		Result:           v,

		// Override the TagName to match the one used by the encoding/json package
		// so users of this function only have to define a single tag on struct fields
		TagName: "json",
	})
	if err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to create decoder")
	}

	// Unmarshal route parameters
	if err := decoder.Decode(mux.Vars(r)); err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to decode route parameters")
	}

	// Query parameters come out as a map[string][]string so we loop through them all
	// to remove the unnecessary slice if the parameter just has a single value
	paramSlices := r.URL.Query()
	params := map[string]interface{}{}
	for key, value := range paramSlices {
		switch len(value) {
		case 0:
			params[key] = nil
		case 1:
			params[key] = value[0]
		default:
			params[key] = value
		}
	}

	// Unmarshal query parameters
	if err := decoder.Decode(params); err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to decode query parameters")
	}

	// If there's no body, return early
	if r.Body == nil {
		return nil
	}

	// Read the body of the request
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return oops.Wrap(err, oops.ErrInternalService, "failed to read request body")
	}

	// If the body is empty, return early
	if len(body) == 0 {
		return nil
	}

	// Assume the body is JSON and unmarshal into v
	if err := json.Unmarshal(body, v); err != nil {
		return oops.Wrap(err, oops.ErrBadRequest, "failed to unmarshal request body")
	}

	return nil
}
