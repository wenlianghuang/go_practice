package main

import (
	"log"
	"os"

	"golangAPI_construct/routes"
)

func main() {
	router := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Println("Starting server on", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
