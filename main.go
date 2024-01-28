package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendChatMessage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, x-sveltekit-action")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Print(req)
	w.WriteHeader(http.StatusOK)
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Hello json",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/chat", sendChatMessage)

	fmt.Printf("listening on port {8080}")
	http.ListenAndServe(":8080", nil)
}
