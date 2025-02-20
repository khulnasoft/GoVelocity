// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.io

package main

import (
	"crypto/tls"
	"log"

	"go.khulnasoft.com/velocity"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Velocity instance
	app := velocity.New()

	// Routes
	app.Get("/", func(c *velocity.Ctx) error {
		return c.SendString("This is a secure server 👮")
	})

	// Let’s Encrypt has rate limits: https://letsencrypt.org/docs/rate-limits/
	// It's recommended to use it's staging environment to test the code:
	// https://letsencrypt.org/docs/staging-environment/

	// Certificate manager
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Replace with your domain
		HostPolicy: autocert.HostWhitelist("example.com"),
		// Folder to store the certificates
		Cache: autocert.DirCache("./certs"),
	}

	// TLS Config
	cfg := &tls.Config{
		// Get Certificate from Let's Encrypt
		GetCertificate: m.GetCertificate,
		// By default NextProtos contains the "h2"
		// This has to be removed since Fasthttp does not support HTTP/2
		// Or it will cause a flood of PRI method logs
		// http://webconcepts.info/concepts/http-method/PRI
		NextProtos: []string{
			"http/1.1", "acme-tls/1",
		},
	}
	ln, err := tls.Listen("tcp", ":443", cfg)
	if err != nil {
		panic(err)
	}

	// Start server
	log.Fatal(app.Listener(ln))
}
