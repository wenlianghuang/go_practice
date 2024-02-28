package mountexample

import "github.com/gofiber/fiber/v2"

func Mountexample() {
	app := fiber.New()
	subApp := fiber.New()

	// 在子應用程式中定義路由
	subApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from SubApp!")
	})
	// 使用Mount將子應用程式掛載到主應用程式的指定路由位置
	app.Mount("/sub", subApp)

	// 在主應用程式中定義其他路由
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from MainApp")
	})
	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
