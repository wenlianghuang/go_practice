package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("About Page"))
	})

	http.ListenAndServe(":8080", nil)
}
