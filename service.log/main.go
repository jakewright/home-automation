package main

import (
	"home-automation/libraries/go/bootstrap"
	"home-automation/libraries/go/router"

	"home-automation/service.log/routes"
)

func main() {
	if err := bootstrap.Init("service.log"); err != nil {
		panic(err)
	}

	router.Get("/read", routes.HandleReadLogs)
	router.ListenAndServe()
}
