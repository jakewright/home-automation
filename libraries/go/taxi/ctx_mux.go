package taxi

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/danielchatfield/go-randutils"
)

// ContextMultiplexer is an interface that
// wraps the Handle function of ContextMux
type ContextMultiplexer interface {
	RegisterHandler(ctx context.Context, handler http.Handler) (context.Context, func())
}

// contextKeyType is a custom type to guarantee uniqueness of the context key
type contextKeyType struct{}

var contextKey = contextKeyType{}

// ContextMux is an HTTP multiplexer that dispatches requests to handlers
// based on a context value in the request.
//
// Since ContextMux implements the http.Handler interface, it can be used as
// the handler in a MockClient. Multiple unit tests can then run in parallel,
// sharing the same MockClient but without sharing the same Handler.
type ContextMux struct {
	handlers map[string]http.Handler
	mu       *sync.RWMutex
}

// NewContextMux returns an initialised ContextMux
func NewContextMux() *ContextMux {
	return &ContextMux{
		handlers: make(map[string]http.Handler),
		mu:       &sync.RWMutex{},
	}
}

// RegisterHandler adds a new handler to the multiplexer. It returns a new
// context that RPCs should use if they are to be handled by the given handler.
// It also returns a function that removes the handler from the multiplexer.
// Typically, this function should be deferred until the end of the test.
func (c *ContextMux) RegisterHandler(ctx context.Context, handler http.Handler) (context.Context, func()) {
	id, err := randutils.String(32)
	if err != nil {
		panic(fmt.Errorf("failed to generate fixture key: %w", err))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers[id] = handler

	stop := func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.handlers, id)
	}

	ctx = context.WithValue(ctx, contextKey, id)
	return ctx, stop
}

// ServeHTTP handles HTTP requests by looking up the context ID in the
// request, and dispatching to the handler that was assigned the same ID.
func (c *ContextMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(contextKey).(string)
	if !ok {
		panic(fmt.Errorf(
			"could not find context ID in request to %s %s", r.Method, r.URL,
		))
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	h, ok := c.handlers[id]
	if !ok {
		panic(fmt.Errorf(
			"could not find handler for context ID %q in request to %s %s",
			id,
			r.Method,
			r.URL,
		))
	}

	h.ServeHTTP(w, r)
}
