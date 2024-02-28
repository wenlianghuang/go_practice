package bodyparser

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Field names should start with an uppercase letter
type Person struct {
	Name string `json:"name" xml:"name" form:"name"`
	Pass string `json:"pass" xml:"pass" form:"pass"`
}

func Bodyparser() {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		p := new(Person)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		fmt.Println(p.Name) // Output: john
		fmt.Println(p.Pass) // Output: doe

		// Additional logic...

		return c.SendString("Received and processed the data")
	})

	// Start server on port 3000
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
