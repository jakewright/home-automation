package router

import (
	"net/http"
	"runtime/debug"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/response"
	routerproto "github.com/jakewright/home-automation/libraries/go/router/proto"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/libraries/go/util"
)

func panicRecovery(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if v := recover(); v != nil {
			stack := debug.Stack()
			err := errors.Wrap(v, errors.ErrPanic, "recovered from panic", map[string]string{
				"stack": string(stack),
			})
			response.WriteJSON(w, err)

			if util.IsProd() {
				slog.Error(err)
			} else {
				// Panicking is useful in dev for the pretty-printed stack trace in terminal
				panic(err)
			}
		}
	}()

	next(w, r)
}

func ping(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == http.MethodGet && r.URL.Path == "/ping" {
		response.WriteJSON(w, &routerproto.PingResponse{
			Ping: "pong",
		})
		return
	}

	next(w, r)
}
