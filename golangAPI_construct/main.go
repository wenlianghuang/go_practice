package main

import "golangAPI_construct/routes"

func main() {
	router := routes.SetupRoutes()
	router.Run("localhost:8080")
}
