package middlewareroute

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func MyRouterMiddleware(ctx *fiber.Ctx) error {
	fmt.Println("Executing router-leve middleware before route handler")
	return ctx.Next()
}
func Middlewareroute() {
	app := fiber.New()
	apiGroup := app.Group("/api", MyRouterMiddleware)
	apiGroup.Get("/hello", func(c *fiber.Ctx) error {
		fmt.Println("Executing route handler in /api/hello")
		return c.SendString("Hello from /api/hello!")
	})

	app.Use(func(ctx *fiber.Ctx) error {
		fmt.Println("Executing global middleware before route handler")
		return ctx.Next()
	})

	app.Get("/world", func(c *fiber.Ctx) error {
		fmt.Println("Executing route handler in /world")
		return c.SendString("Hello from /world!")
	})

	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
