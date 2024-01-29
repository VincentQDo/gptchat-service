package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Define Enum
type Role string

// Define Enum values
const (
	System    Role = "system"
	User      Role = "user"
	Assistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type OpenAiRequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func sendChatMessage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Unsupported method.", http.StatusMethodNotAllowed)
	}

	url := "https://api.openai.com/v1/chat/completions"
	apiKey := "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"

	var incomingMessages []Message
	err := json.NewDecoder(req.Body).Decode(&incomingMessages)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("%v", incomingMessages)
	body := OpenAiRequestBody{
		Model:    "gpt-4-turbo-preview", // Update if needed with the correct model
		Messages: incomingMessages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing json: %s", err), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, "Error creating request to OpenAI", http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	res, err := client.Do(request)
	if err != nil {
		http.Error(w, "Error sending request to OpenAI", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		http.Error(w, fmt.Sprintf("OpenAI request failed: %s", body), res.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(res.Body)
	for {
		var message Message
		if err := decoder.Decode(&message); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("error decoding OpenAI response: %v", err)
			// You may want to handle the error more gracefully
			break
		}

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("error marshalling message: %v", err)
			// You may want to handle the error more gracefully
			break
		}

		// Writes an SSE event to the client, with the message as the data field
		fmt.Fprintf(w, "data: %s\n\n", jsonMessage)

		flusher.Flush() // Send to client immediately
	}

	// Handle remaining business logic here if needed
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
