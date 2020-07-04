package taxi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RPC represents a remote procedure call
type RPC struct {
	Method string
	URL    string
	Body   interface{}
}

// ToRequest converts the RPC into an http.Request
func (r *RPC) ToRequest(ctx context.Context) (*http.Request, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(r.Body); err != nil {
		return nil, fmt.Errorf("failed to encode body as JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	return req, nil
}

// RPC returns itself so an RPC implements the
// rpcProvider interface used by the TestFixture
func (r *RPC) RPC() *RPC {
	return r
}
