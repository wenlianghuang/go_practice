package groupexample

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Groupexmaple() {
	// Build Fiber application\
	app := fiber.New()

	//Build a group of route
	apiGroup := app.Group("/api")

	// Add a middleware in the group of route
	apiGroup.Use(func(ctx *fiber.Ctx) error {
		fmt.Println("Executing middleware for /api group")
		return ctx.Next()
	})

	// Add a new route in the group
	apiGroup.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello from /api/hello!")
	})

	// Add lots of route in the group
	apiGroup.Get("/greet/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString(fmt.Sprintf("Greetings, %s!", name))
	})

	apiGroup.Post("/add", func(c *fiber.Ctx) error {
		data := new(struct {
			Number1 int `json:"number1"`
			Number2 int `json:"number2"`
		})
		if err := c.BodyParser(data); err != nil {
			return err
		}

		sum := data.Number1 + data.Number2
		return c.JSON(fiber.Map{"sum": sum})
	})

	// Add a route in the main application
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from the main app!")
	})

	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
