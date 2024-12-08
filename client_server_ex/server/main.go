package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("localfile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a new file in the server to save the uploaded file
	dst, err := os.Create("./uploaded_" + handler.Filename)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the server file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server started at :9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
