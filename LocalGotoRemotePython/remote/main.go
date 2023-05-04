package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Get request from local
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				http.Error(w, "can't read body", http.StatusBadRequest)
			}

			fmt.Println(string(body))

			//response some more detail
			response := map[string]interface{}{
				"status":  "success",
				"message": "You request has been successfully processed.",
			}
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
