package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/khulnasoft/template/html/v2"
	"go.khulnasoft.com/velocity"

	"geoip/handlers"
)

func main() {
	// Create new velocity instance
	app := velocity.New(velocity.Config{
		// Pass view engine
		Views: html.New("./views", ".html"),
		// Pass global error handler
		ErrorHandler: handlers.Errors("./public/500.html"),
	})

	// Render index template with IP input value
	app.Get("/", handlers.Render())

	// Serve static assets
	app.Static("/", "./public", velocity.Static{
		Compress: true,
	})

	// Main GEO handler that is cached for 10 minutes
	app.Get("/geo", handlers.Cache(10*time.Minute), handlers.GEO())

	// Handle 404 errors
	app.Use(handlers.NotFound("./public/404.html"))

	// Listen on environment specified PORT, "3000" otherwise
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
