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

	app.Post("/", func(c *velocity.Ctx) error {
		// Parse the multipart form:
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		// => *multipart.Form

		// Get all files from "documents" key:
		files := form.File["documents"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
			// Check for errors
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
