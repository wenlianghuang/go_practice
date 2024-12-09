package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

	// register an sample API endpoint
	router.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Hello, World!"}`)
	}).Methods("GET")

	//
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://example.com"}), // 允許的來源
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),                 // 允許的 HTTP 方法
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),                // 允許的標頭
		handlers.AllowCredentials(),
	)
	http.ListenAndServe(":8080", cors(router))
}
