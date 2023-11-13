package getroutes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Getroutes() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	app.Get("/api/users", func(c *fiber.Ctx) error {
		return c.SendString("API: List of users")
	})

	app.Post("/api/users", func(c *fiber.Ctx) error {
		return c.SendString("API: Create a new user")
	})

	routes := app.GetRoutes()

	for _, route := range routes {
		fmt.Printf("Method: %s, Path: %s, Handler:%v\n", route.Method, route.Path, route.Handlers)
	}

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
