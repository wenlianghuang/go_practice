package middlewareex

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func MyMiddleware(ctx *fiber.Ctx) error {
	fmt.Println("Executing middleware before route handler")
	return ctx.Next()
}

func Middleware() {
	app := fiber.New()

	app.Use(MyMiddleware)
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Executing route handler")
		return c.SendString("Hello,Fiber!")
	})
	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
