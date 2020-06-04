package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/danielchatfield/go-randutils"
	"github.com/stretchr/testify/require"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// fixtureIDKeyType is a custom type to guarantee uniqueness of the context key
type fixtureIDKeyType string

const fixtureIDKey fixtureIDKeyType = "mock"

// Matcher matches RPC requests and returns a canned response
type Matcher interface {
	Match(*Request) bool
	Serve(*Request) (*Response, error)
}

// MockResponse is a canned response for a Matcher to return
type MockResponse struct {
	Status int
	Body   []byte
}

// TestClient can be used by unit tests to mock RPCs
type TestClient struct {
	fixtures map[string]*Fixture
	mu       *sync.RWMutex
}

// Do simulates performing an RPC
func (c *TestClient) Do(ctx context.Context, request *Request, resultData interface{}) (*Response, error) {
	id, ok := ctx.Value(fixtureIDKey).(string)
	if !ok {
		return nil, oops.InternalService("RPC made without using mock context")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	f, ok := c.fixtures[id]
	if !ok {
		panic(fmt.Errorf("no fixture found for request to %s", request.URL))
	}

	return f.handle(request, resultData)
}

// Fixture is a set of matchers for a specific test run
type Fixture struct {
	id       string
	t        *testing.T
	matchers []Matcher
	mu       *sync.RWMutex
}

// NewFixture returns a new fixture
func NewFixture(ctx context.Context, t *testing.T) (*Fixture, context.Context) {
	SetDefaultClient(&TestClient{
		fixtures: make(map[string]*Fixture),
		mu:       &sync.RWMutex{},
	})

	id, err := randutils.String(32)
	require.NoError(t, err, "failed to generate fixture key")

	m, ok := defaultClient.(*TestClient)
	require.True(t, ok, "defaultClient was unexpected type")

	m.mu.Lock()
	defer m.mu.Unlock()

	_, existing := m.fixtures[id]
	require.False(t, existing, "fixture with ID already exists")

	m.fixtures[id] = &Fixture{
		id:       id,
		t:        t,
		matchers: make([]Matcher, 0),
		mu:       &sync.RWMutex{},
	}

	ctx = context.WithValue(ctx, fixtureIDKey, id)
	return m.fixtures[id], ctx
}

// Stop removes the fixture from the mock client
func (f *Fixture) Stop() {
	m, ok := defaultClient.(*TestClient)
	require.True(f.t, ok, "defaultClient was unexpected type")

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.fixtures, f.id)
}

// AddMatcher adds a matcher to the fixture
func (f *Fixture) AddMatcher(m Matcher) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.matchers = append(f.matchers, m)
}

func (f *Fixture) handle(request *Request, resultData interface{}) (*Response, error) {
	rsp, err := f.match(request).Serve(request)
	if err != nil {
		return rsp, err
	}

	if rsp != nil && len(rsp.Body) == 0 {
		if err := json.Unmarshal(rsp.Body, resultData); err != nil {
			return rsp, err
		}
	}

	return rsp, nil
}

func (f *Fixture) match(request *Request) Matcher {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, m := range f.matchers {
		if m.Match(request) {
			return m
		}
	}

	f.t.Fatalf("No matcher for request to %s %+v", request.URL, request.Body)
	return nil
}
