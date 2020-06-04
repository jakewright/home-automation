package test

import (
	"context"
	"testing"

	"github.com/jakewright/home-automation/libraries/go/rpc"
)

// Mock provides various helpers to unit tests
type Mock struct {
	t          *testing.T
	rpcFixture *rpc.Fixture
	rpcStubs   []*Stub
}

// NewMock sets up a new test environment
func NewMock(t *testing.T) (*Mock, context.Context) {
	f, ctx := rpc.NewFixture(context.Background(), t)

	return &Mock{
		t:          t,
		rpcFixture: f,
	}, ctx
}

// RunAssertions asserts that the expectations have been met
func (m *Mock) RunAssertions() {
	for _, stub := range m.rpcStubs {
		if stub.expect >= 0 && stub.expect != stub.count {
			m.t.Fatalf("expected %d requests to %s but got %d", stub.expect, stub.req.URL, stub.count)
		}
	}
}

// Stop performs cleanup and should be deferred in unit tests
func (m *Mock) Stop() {
	m.rpcFixture.Stop()
}

// RPC is an interface that svcdef request types conform to
type RPC interface {
	Request() *rpc.Request
}

// ExpectN sets up a stub for the RPC that expects n matching requests
func (m *Mock) ExpectN(n int, rpc RPC) *Stub {
	stub := &Stub{
		t:      m.t,
		expect: n,
		req:    rpc.Request(),
	}

	m.rpcFixture.AddMatcher(stub)
	m.rpcStubs = append(m.rpcStubs, stub)
	return stub
}

// ExpectOne sets up a stub for the RPC that expects one matching request
func (m *Mock) ExpectOne(rpc RPC) *Stub {
	return m.ExpectN(1, rpc)
}

// Allow sets up a stub for the RPC that allows any number of requests
func (m *Mock) Allow(rpc RPC) *Stub {
	return m.ExpectN(-1, rpc)
}
