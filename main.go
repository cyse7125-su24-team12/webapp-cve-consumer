package main

import (
	"net/http"
)

func main() {
	// Start the HTTP server in a Goroutine so it runs independently
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			// Respond with HTTP 200 OK
			w.WriteHeader(http.StatusOK)
		})
		if err := http.ListenAndServe(":80", nil); err != nil {
			// Handle potential errors, such as port already in use
			panic(err)
		}
	}()

	ConsumerKafka();
}
