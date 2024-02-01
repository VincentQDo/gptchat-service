package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendChatMessage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Unsupported method.", http.StatusMethodNotAllowed)
	}

	// Set the headers here
	w.Header().Set("Content-Type", "application/json")

	// WriteHeader will return the correct status
	w.WriteHeader(http.StatusOK)
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Hello json",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, x-sveltekit-action")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/chat", sendChatMessage)
	// Wrap the http handler in cors middle ware
	handlerWithCors := corsMiddleware(multiplexer)

	fmt.Printf("listening on port {8080}")
	http.ListenAndServe(":8080", handlerWithCors)
}
