// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/fiber
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"log"

	"github.com/khulnasoft/fiber/v2"
	"github.com/khulnasoft/fiber/v2/middleware/recover"
)

func main() {
	// Fiber instance
	app := fiber.New(fiber.Config{
		// ErrorHandler: func(c *fiber.Ctx, err error) error {
		// 	c.Status(fiber.StatusInternalServerError)
		// 	return c.SendString(err.Error())
		// },
	})

	// Middleware
	app.Use(recover.New())

	// Routes
	app.Get("/", hello)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Handler
func hello(c *fiber.Ctx) error {
	panic("No worries, I won't crash! 🙏")
}
