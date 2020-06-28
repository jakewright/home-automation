package taxi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFixture is an http.Handler that can be used in testing. It dispatches
// requests to one of a set of Stubs that match based on arbitrary details
// about the request. Assertions can be run to assert that expected requests
// were received.
type TestFixture struct {
	t     *testing.T
	stubs []*Stub
	mu    *sync.Mutex
}

// NewTestFixture returns an initialised TestFixture
func NewTestFixture(t *testing.T) *TestFixture {
	return &TestFixture{
		t:  t,
		mu: &sync.Mutex{},
	}
}

// RegisterTestFixture sets a mock client using a ContextMux as the default
// Dispatcher and registers the given TestFixture as a handler. A modified
// context is returned that can be used by RPCs to have them be handled
// by the TestFixture.
func RegisterTestFixture(ctx context.Context, f *TestFixture) context.Context {
	SetDefaultDispatcher(NewMockClient(NewContextMux()))

	dispatcher, ok := mustGetDefaultDispatcher().(*MockClient)
	require.True(f.t, ok, "default Dispatcher was unexpected type")

	mux, ok := dispatcher.handler.(ContextMultiplexer)
	require.True(f.t, ok, "mock client's handler was unexpected type")

	ctx, stop := mux.RegisterHandler(ctx, f)
	f.t.Cleanup(stop)
	return ctx
}

// ServeHTTP dispatches requests to the first stub that matches
func (f *TestFixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, s := range f.stubs {
		if s.Match(r) {
			s.ServeHTTP(w, r)
			return
		}
	}

	f.t.Fatalf("could not find handler for request to %s %s", r.Method, r.URL)
}

// RunAssertions runs all of the stubs' assertions
func (f *TestFixture) RunAssertions() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, s := range f.stubs {
		s.RunAssertions()
	}
}

// Expect returns a new stub that matches on the method and path, and asserts
// that it is called n times.
func (f *TestFixture) Expect(n int, method, path string) *Stub {
	f.mu.Lock()
	defer f.mu.Unlock()

	s := NewStub(f.t).Expect(n).MatchMethod(method).MatchPath(path)
	f.stubs = append(f.stubs, s)
	return s
}

// Allow returns a new stub that matches on the method and path but does
// not care how many times it is called.
func (f *TestFixture) Allow(method, path string) *Stub {
	f.mu.Lock()
	defer f.mu.Unlock()

	s := NewStub(f.t).MatchMethod(method).MatchPath(path)
	f.stubs = append(f.stubs, s)
	return s
}

// Stub is an http.Handler to use with a TestFixture. It can be configured to
// match requests based on arbitrary rules, and make assertions based on the
// requests that were received.
type Stub struct {
	t          *testing.T
	requests   []*http.Request
	mu         *sync.Mutex
	matchers   []func(r *http.Request) bool
	assertions []func(t *testing.T, requests []*http.Request)
	serve      func(w http.ResponseWriter, r *http.Request)
	rw         *responseWriter
}

// NewStub returns an initialised Stub
func NewStub(t *testing.T) *Stub {
	return &Stub{
		t:  t,
		mu: &sync.Mutex{},
		rw: &responseWriter{
			logFunc: t.Logf,
		},
	}
}

// ServeHTTP tracks that the request was received by the stub, and calls
// the stub's serve function if not nil.
func (s *Stub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, r)
	if s.serve != nil {
		s.serve(w, r)
	}
}

// Match returns whether this stub should handle the given request
func (s *Stub) Match(r *http.Request) bool {
	if len(s.matchers) < 1 {
		s.t.Fatal("stub created with no matchers")
	}

	for _, m := range s.matchers {
		if !m(r) {
			return false
		}
	}

	return true
}

// RunAssertions runs the stub's assertions
func (s *Stub) RunAssertions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.assertions {
		a(s.t, s.requests)
	}
}

// MatchMethod adds a matcher that matches on the request method
func (s *Stub) MatchMethod(method string) *Stub {
	return s.WithMatcher(func(r *http.Request) bool {
		return r.Method == method
	})
}

// MatchPath adds a matcher that matches on the request path
func (s *Stub) MatchPath(path string) *Stub {
	return s.WithMatcher(func(r *http.Request) bool {
		return r.URL.Path == path
	})
}

// MatchBody adds a matcher that matches on the request body
func (s *Stub) MatchBody(body interface{}) *Stub {
	expectBytes, err := json.Marshal(body)
	require.NoError(s.t, err)

	var expect interface{}
	err = json.Unmarshal(expectBytes, &expect)
	require.NoError(s.t, err)

	return s.WithMatcher(func(r *http.Request) bool {
		actualBytes, err := ioutil.ReadAll(r.Body)
		require.NoError(s.t, err)
		defer func() { _ = r.Body.Close() }()

		var actual interface{}
		err = json.Unmarshal(actualBytes, &actual)
		require.NoError(s.t, err)

		return reflect.DeepEqual(&expect, &actual)
	})
}

// MatchPartialBody adds a matcher that matches on the request body but
// only compares fields that are set in the fields slice.
func (s *Stub) MatchPartialBody(body interface{}, fields []string) *Stub {
	expectBytes, err := json.Marshal(body)
	require.NoError(s.t, err)

	var expect map[string]interface{}
	err = json.Unmarshal(expectBytes, &expect)
	require.NoError(s.t, err)

	return s.WithMatcher(func(r *http.Request) bool {
		actualBytes, err := ioutil.ReadAll(r.Body)
		require.NoError(s.t, err)
		defer func() { _ = r.Body.Close() }()

		var actual map[string]interface{}
		err = json.Unmarshal(actualBytes, &actual)
		require.NoError(s.t, err)

		for _, f := range fields {
			if !reflect.DeepEqual(expect[f], actual[f]) {
				return false
			}
		}

		return true
	})
}

// WithMatcher adds a custom matcher
func (s *Stub) WithMatcher(f func(r *http.Request) bool) *Stub {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.matchers = append(s.matchers, f)
	return s
}

// Expect adds an assertion that the stub receives n requests
func (s *Stub) Expect(n int) *Stub {
	return s.WithAssertion(func(t *testing.T, requests []*http.Request) {
		got := len(requests)
		require.Equal(t, n, got, "Expected %d requests but got %d", n, got)
	})
}

// WithAssertion adds a custom assertion
func (s *Stub) WithAssertion(f func(t *testing.T, requests []*http.Request)) *Stub {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.assertions = append(s.assertions, f)
	return s
}

// RespondWith sets the data that should be returned by the stub when handling
// requests. The interface will be marshaled to JSON and wrapped in a data
// field. If v is an error, the string value will be put in an error field
// in JSON. The logic is the same as that of Taxi router handlers.
func (s *Stub) RespondWith(v interface{}) *Stub {
	s.serve = func(w http.ResponseWriter, r *http.Request) {
		s.rw.writeResponse(w, v)
	}

	return s
}
