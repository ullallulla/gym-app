package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"status": "ok",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	})

	log.Println("API listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}
