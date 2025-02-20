// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"fmt"
	"log"

	"go.khulnasoft.com/velocity"
)

func main() {
	// Velocity instance
	app := velocity.New()

	// Routes
	app.Post("/", func(c *velocity.Ctx) error {
		// Get first file from form field "document":
		file, err := c.FormFile("document")
		if err != nil {
			return err
		}
		// Save file inside uploads folder under current working directory:
		return c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
	})

	app.Post("/temp", func(c *velocity.Ctx) error {
		// Get first file from form field "document":
		file, err := c.FormFile("document")
		if err != nil {
			return err
		}
		//(uploads_relative)folder must be created before hand:
		//Save file using a relative path:
		return c.SaveFile(file, fmt.Sprintf("/tmp/uploads_relative/%s", file.Filename))
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
