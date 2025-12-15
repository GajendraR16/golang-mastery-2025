package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/hello/", hello) // Note the trailing slash
	http.HandleFunc("/echo", echo)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
