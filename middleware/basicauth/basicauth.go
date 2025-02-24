package basicauth

import (
	"encoding/base64"
	"strings"

	"github.com/khulnasoft/velocity"
	"github.com/khulnasoft/velocity/utils"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// The keys for the values in context
const (
	usernameKey contextKey = iota
	passwordKey
)

// New creates a new middleware handler
func New(config Config) velocity.Handler {
	// Set default config
	cfg := configDefault(config)

	// Return new handler
	return func(c velocity.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get authorization header
		auth := c.Get(velocity.HeaderAuthorization)

		// Check if the header contains content besides "basic".
		if len(auth) <= 6 || !utils.EqualFold(auth[:6], "basic ") {
			return cfg.Unauthorized(c)
		}

		// Decode the header contents
		raw, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return cfg.Unauthorized(c)
		}

		// Get the credentials
		creds := utils.UnsafeString(raw)

		// Check if the credentials are in the correct form
		// which is "username:password".
		index := strings.Index(creds, ":")
		if index == -1 {
			return cfg.Unauthorized(c)
		}

		// Get the username and password
		username := creds[:index]
		password := creds[index+1:]

		if cfg.Authorizer(username, password) {
			c.Locals(usernameKey, username)
			c.Locals(passwordKey, password)
			return c.Next()
		}

		// Authentication failed
		return cfg.Unauthorized(c)
	}
}

// UsernameFromContext returns the username found in the context
// returns an empty string if the username does not exist
func UsernameFromContext(c velocity.Ctx) string {
	username, ok := c.Locals(usernameKey).(string)
	if !ok {
		return ""
	}
	return username
}

// PasswordFromContext returns the password found in the context
// returns an empty string if the password does not exist
func PasswordFromContext(c velocity.Ctx) string {
	password, ok := c.Locals(passwordKey).(string)
	if !ok {
		return ""
	}
	return password
}
