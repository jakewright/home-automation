package dmx

import (
	"context"
	"sync"

	"github.com/jakewright/home-automation/services/dmx/domain"
)

// GetSetter is an interface for interacting with a DMX universe
type GetSetter interface {
	GetValues(ctx context.Context) ([512]byte, error)
	SetValues(ctx context.Context, values [512]byte) error
}

// Client can get and set DMX values for all universes
type Client struct {
	getSetters map[domain.UniverseNumber]GetSetter
	mu         *sync.Mutex
}

// NewClient returns a new client
func NewClient() *Client {
	return &Client{
		getSetters: make(map[domain.UniverseNumber]GetSetter),
		mu:         &sync.Mutex{},
	}
}

// AddGetSetter adds a GetSetter to the client
func (c *Client) AddGetSetter(un domain.UniverseNumber, gs GetSetter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.getSetters[un] = gs
}

// GetValues returns the DMX values for the specified universe
func (c *Client) GetValues(ctx context.Context, un domain.UniverseNumber) ([512]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getSetters[un].GetValues(ctx)
}

// SetValues sets the DMX values for the specified universe
func (c *Client) SetValues(ctx context.Context, un domain.UniverseNumber, values [512]byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getSetters[un].SetValues(ctx, values)
}
