package router

import (
	"net/http"
	"runtime/debug"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/taxi"
)

func panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				stack := debug.Stack()
				err := oops.Wrap(v, oops.ErrInternalService, "recovered from panic", map[string]string{
					"stack": string(stack),
				})
				if err := taxi.WriteError(w, err); err != nil {
					slog.Errorf("Failed to write response: %v", err)
				}

				slog.Error(err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func revision(revision string) taxi.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Revision", revision)
			next.ServeHTTP(w, r)
		})
	}
}
