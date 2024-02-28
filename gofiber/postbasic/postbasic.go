package postbasic

import "github.com/gofiber/fiber/v2"

type RequestBody struct {
	Name string `json:"name"`
}

func Postbasic() {
	app := fiber.New()

	// Create a new endpoint
	app.Post("/", func(c *fiber.Ctx) error {

		// Parse JSON data from the request body
		var requestBody RequestBody
		if err := c.BodyParser(&requestBody); err != nil {
			// Handle invalid JSON format
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON format")
		}

		// Get the name from the parsed request body
		name := requestBody.Name
		name = "Matt"
		// Create a map for the JSON response
		response := fiber.Map{
			"message": "Hello, World!",
			"name":    name,
		}

		// Convert the map to JSON and send it as the response
		return c.JSON(response)
	})

	// Start server on port 5050
	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
