// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Velocity instance
	app := velocity.New()

	// Routes
	app.Get("/hello", hello)

	// 404 Handler
	app.Use(func(c *velocity.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Handler
func hello(c *velocity.Ctx) error {
	return c.SendString("I made a ☕ for you!")
}
