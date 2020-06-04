package test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jakewright/home-automation/libraries/go/rpc"
)

// Stub returns a static response to a request
type Stub struct {
	t          *testing.T
	expect     int
	count      int
	req        *rpc.Request
	ignoreBody bool
	rsp        *rpc.Response
}

// IgnoreBody will match a request based on method and URL only
func (s *Stub) IgnoreBody() *Stub {
	s.ignoreBody = true
	return s
}

// RespondWith will return a response that contains the
// argument marshaled to JSON.
func (s *Stub) RespondWith(v interface{}) *Stub {
	body, err := json.Marshal(v)
	require.NoError(s.t, err)

	s.rsp = &rpc.Response{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			// TODO: fill in more fields if needed
		},
		Body: body,
	}

	return s
}

// Match conforms to the Matcher interface
func (s *Stub) Match(request *rpc.Request) bool {
	if request.Method != s.req.Method {
		return false
	}

	if request.URL != s.req.URL {
		return false
	}

	if s.ignoreBody {
		return true
	}

	return reflect.DeepEqual(s.req.Body, request.Body)
}

// Serve conforms to the Matcher interface
func (s *Stub) Serve(_ *rpc.Request) (*rpc.Response, error) {
	s.count++
	return s.rsp, nil
}
