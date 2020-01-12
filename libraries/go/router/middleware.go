package router

import (
	"net/http"
	"runtime/debug"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/response"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

func panicRecovery(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if v := recover(); v != nil {
			stack := debug.Stack()
			err := errors.Wrap(v, errors.ErrPanic, "recovered from panic", map[string]string{
				"stack": string(stack),
			})
			slog.Error(err)
			response.WriteJSON(w, err)
		}
	}()

	next(w, r)
}
