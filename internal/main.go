package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)

	http.ListenAndServe(":8080", mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}
