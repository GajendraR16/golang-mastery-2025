package main

import (
	"encoding/json"
	"net/http"
)

type ResponseStatus struct {
	Status string `json:"status"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func jsonHandler(w http.ResponseWriter, data any) {
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to Encode Json", http.StatusInternalServerError)
	}
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "Method not allowed", http.StatusInternalServerError)
	}
	jsonHandler(w, ResponseStatus{Status: "ok"})
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "Method not allowed", http.StatusInternalServerError)
		return
	}
	name := "World!"
	if len(r.URL.Path) > len("/hello/") {
		extracted := r.URL.Path[len("/hello/"):]
		if name != "" {
			name = extracted
		}
	}

	jsonHandler(w, MessageResponse{Message: "Hello, " + name + "!"})
}

func echo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	data := make(map[string]any)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonError(w, "Invalid Json", http.StatusInternalServerError)
		return
	}

	jsonHandler(w, data)
}
