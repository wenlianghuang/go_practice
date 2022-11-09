package gethttp

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Gethttp() {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Fatalf("Response failed with status code %d and \n body: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}
