package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func postJSON(filename string, targetUrl string) error {
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the JSON data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Create a buffer to write the multipart form data
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	// Add the JSON file to the multipart form data
	filePart, err := writer.CreateFormFile("localfile", filename)
	if err != nil {
		return err
	}
	_, err = filePart.Write(data)
	if err != nil {
		return err
	}

	// Close the multipart writer
	writer.Close()

	// Create the HTTP request with the multipart form data
	request, err := http.NewRequest("POST", targetUrl, requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the response body
	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Print the response status and content
	fmt.Println(response.Status)
	fmt.Println(string(responseContent))

	return nil
}

func main() {
	// Set the target URL and the JSON file path
	targetUrl := "http://localhost:9090/upload"
	filename := "./example.json"

	// Upload the JSON file
	err := postJSON(filename, targetUrl)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
