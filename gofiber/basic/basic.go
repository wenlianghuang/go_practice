package basic

import "github.com/gofiber/fiber/v2"

func Basic() {
	app := fiber.New() // create a new Fiber instance

	// Create a new endpoint
	//app.Get("/", func(c *fiber.Ctx) error {
	//	return c.SendString("Hello, World!") // send text
	//})
	app.Get("/", func(c *fiber.Ctx) error {
		name := c.Query("name")

		response := fiber.Map{
			"message": "Hello, World",
			"name":    name,
		}

		return c.JSON(response)
	})

	// Start server on port 3000
	app.Listen(":5050")
}
