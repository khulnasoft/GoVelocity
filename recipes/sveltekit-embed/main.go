package main

import (
	"app/template"

	"github.com/khulnasoft/fiber/v2"
	"github.com/khulnasoft/fiber/v2/middleware/cors"
	"github.com/khulnasoft/fiber/v2/middleware/filesystem"
	"github.com/khulnasoft/fiber/v2/middleware/logger"
)

const (
	appName    = "Sveltekit Embed App"
	apiVersion = "v1"
	port       = ":3000"
)

func main() {
	// Create new Fiber Instance
	app := fiber.New(fiber.Config{
		AppName: appName,
	})
	defer app.Shutdown()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	// Serve static files
	app.All("/*", filesystem.New(filesystem.Config{
		Root:         template.Dist(),
		NotFoundFile: "index.html",
		Index:        "index.html",
	}))

	// Start the server
	if err := app.Listen(port); err != nil {
		panic(err)
	}
}
