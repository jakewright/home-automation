package healthz

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

var (
	checks map[string]HealthCheck
	mu     sync.RWMutex // guards checks
)

func init() {
	checks = make(map[string]HealthCheck)
}

// HealthCheck is a function that returns an error if a
// component is unhealthy. The application's overall health
// depends on the set of registered HealthChecks returning
// nil errors.
type HealthCheck func(ctx context.Context) error

// RegisterCheck adds a new check to the global register.
// If a check with the given name already exists, the
// function will panic.
func RegisterCheck(name string, check HealthCheck) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := checks[name]; exists {
		panic(oops.InternalService("check %q already exists", name))
	}

	checks[name] = check
}

// Status runs all of the currently registered health
// checks and returns a map of the results.
func Status(ctx context.Context) map[string]error {
	results := make(map[string]error, len(checks))
	for name, check := range checks {
		results[name] = safeRunHealthCheck(ctx, check)
	}
	return results
}

// safeRunHealthCheck runs the given HealthCheck but
// recovers any panics. In the case of a panic, up to 1 MB
// of stack trace is returned in the error.
func safeRunHealthCheck(ctx context.Context, check HealthCheck) (err error) {
	defer func() {
		if v := recover(); v != nil {
			buf := make([]byte, 1<<20) // 1 MB
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = fmt.Errorf("panic: %v\n\n%s", v, buf)
		}
	}()

	err = check(ctx)
	return
}
