package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type RequestData struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

type ResponseData struct {
	Message   string `json:"message"`
	Modified  bool   `json:"modified"`
	NewNumber int    `json:"newNumber"`
}

var (
	requestData  = RequestData{Message: "Initial message", Number: 5}
	requestMutex sync.Mutex
)

func main() {
	http.HandleFunc("/testpy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			requestMutex.Lock()
			defer requestMutex.Unlock()

			// Convert request data to JSON
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				http.Error(w, "JSON marshaling failed", http.StatusInternalServerError)
				return
			}

			// Set response headers
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Write JSON data as response
			w.Write(jsonData)
			return
		} else if r.Method == http.MethodPost {
			// Read request body
			body, _ := ioutil.ReadAll(r.Body)

			// Parse JSON data from request
			var requestDataChange RequestData
			err := json.Unmarshal(body, &requestDataChange)
			if err != nil {
				http.Error(w, "JSON unmarshaling failed", http.StatusBadRequest)
				return
			}

			requestMutex.Lock()
			defer requestMutex.Unlock()

			// Modify request data
			requestData.Number += requestDataChange.Number
			requestData.Message = "Modified by server"
			fmt.Println(requestData.Message)

			// Convert response data to JSON
			responseData := ResponseData{
				Message:   requestData.Message,
				Modified:  true,
				NewNumber: requestData.Number,
			}
			jsonData, err := json.Marshal(responseData)
			if err != nil {
				http.Error(w, "JSON marshaling failed", http.StatusInternalServerError)
				return
			}

			// Set response headers
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Write JSON data as response
			w.Write(jsonData)
			return
		}
	})

	http.ListenAndServe("10.36.172.78:8080", nil)
}
