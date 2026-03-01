package main

import (
	"net/http"
	"net/http/cgi"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", helloHandler)

	err := cgi.Serve(mux)
	if err != nil {
		return
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "message": "Hello from CGI"}`))
}
