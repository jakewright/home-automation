package muxinator

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

// Middleware replicates the negroni.HandlerFunc type but decouples the code from the library
type Middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Router is a wrapper around the gorilla mux router and the negroni middleware library.
// It has some convenience functions to make it easier to do per-route middleware
type Router struct {
	n *negroni.Negroni
	m *mux.Router
}

// NewRouter returns a new Router instance with some defaults
func NewRouter() Router {
	n := negroni.New()
	m := mux.NewRouter().StrictSlash(true)
	return Router{n, m}
}

// BuildHandler returns an http.Handler that can be used as the argument to http.ListenAndServe.
func (router *Router) BuildHandler() http.Handler {
	// The mux router needs to be the last item of middleware added to the negroni instance.
	router.n.UseHandler(router.m)
	return router.n
}

// AddMiddleware adds middleware that will be applied to every request.
// Middleware handlers are executed in the order defined.
func (router *Router) AddMiddleware(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		router.n.UseFunc(middleware)
	}
}

// handle registers a route with the router. Internally, gorilla mux is used - see https://github.com/gorilla/mux for
// options available for the path, including variables.
func (router *Router) handle(method string, path string, handler http.HandlerFunc, middlewares ...Middleware) {
	// A slice to hold all of the middleware once it's converted (including the handler itself)
	var stack []negroni.Handler

	// The middleware functions have type Middleware but they need to conform to the negroni.Handler interface.
	// By using the negroni.HandlerFunc adapter, they will be given the method required by the interface.
	for _, middleware := range middlewares {
		stack = append(stack, negroni.HandlerFunc(middleware))
	}

	// The handler needs to be treated like middleware. The negroni.Wrap function can convert an http.Handler into
	// a negroni.Handler, but first we need to transform the handler into an http.Handler by using the adapter.
	stack = append(stack, negroni.Wrap(http.HandlerFunc(handler)))

	// Handle this path using a new instance of negroni with all of the middleware in our stack
	router.m.Handle(path, negroni.New(stack...)).Methods(method)
}

// Get is a helper function to add a GET route
func (router *Router) Get(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.handle("GET", path, handler, middlewares...)
}

// Post is a helper function to add a POST route
func (router *Router) Post(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.handle("POST", path, handler, middlewares...)
}

// Put is a helper function to add a PUT route
func (router *Router) Put(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.handle("PUT", path, handler, middlewares...)
}

// Patch is a helper function to add a PATCH route
func (router *Router) Patch(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.handle("PATCH", path, handler, middlewares...)
}

// Delete is a helper function to add a DELETE route
func (router *Router) Delete(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	router.handle("DELETE", path, handler, middlewares...)
}
