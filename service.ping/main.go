package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)

	fmt.Printf("Listening on %s\n", addr)

	http.HandleFunc("/ping", handlePing)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func handlePing(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "pong")
}
