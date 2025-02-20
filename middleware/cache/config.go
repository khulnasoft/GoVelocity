package cache

import (
	"time"

	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/lib/utils"
)

// Config defines the config for middleware.
type Config struct {
	// Store is used to store the state of the middleware
	//
	// Default: an in memory store for this process only
	Storage velocity.Storage

	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c velocity.Ctx) bool

	// CacheInvalidator defines a function to invalidate the cache when returned true
	//
	// Optional. Default: nil
	CacheInvalidator func(velocity.Ctx) bool

	// Key allows you to generate custom keys, by default c.Path() is used
	//
	// Default: func(c velocity.Ctx) string {
	//   return utils.CopyString(c.Path())
	// }
	KeyGenerator func(velocity.Ctx) string

	// allows you to generate custom Expiration Key By Key, default is Expiration (Optional)
	//
	// Default: nil
	ExpirationGenerator func(velocity.Ctx, *Config) time.Duration

	// CacheHeader header on response header, indicate cache status, with the following possible return value
	//
	// hit, miss, unreachable
	//
	// Optional. Default: X-Cache
	CacheHeader string

	// You can specify HTTP methods to cache.
	// The middleware just caches the routes of its methods in this slice.
	//
	// Default: []string{velocity.MethodGet, velocity.MethodHead}
	Methods []string

	// Expiration is the time that an cached response will live
	//
	// Optional. Default: 1 * time.Minute
	Expiration time.Duration

	// Max number of bytes of response bodies simultaneously stored in cache. When limit is reached,
	// entries with the nearest expiration are deleted to make room for new.
	// 0 means no limit
	//
	// Default: 0
	MaxBytes uint

	// CacheControl enables client side caching if set to true
	//
	// Optional. Default: false
	CacheControl bool

	// allows you to store additional headers generated by next middlewares & handler
	//
	// Default: false
	StoreResponseHeaders bool
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:             nil,
	Expiration:       1 * time.Minute,
	CacheHeader:      "X-Cache",
	CacheControl:     false,
	CacheInvalidator: nil,
	KeyGenerator: func(c velocity.Ctx) string {
		return utils.CopyString(c.Path())
	},
	ExpirationGenerator:  nil,
	StoreResponseHeaders: false,
	Storage:              nil,
	MaxBytes:             0,
	Methods:              []string{velocity.MethodGet, velocity.MethodHead},
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if int(cfg.Expiration.Seconds()) == 0 {
		cfg.Expiration = ConfigDefault.Expiration
	}
	if cfg.CacheHeader == "" {
		cfg.CacheHeader = ConfigDefault.CacheHeader
	}
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = ConfigDefault.KeyGenerator
	}
	if len(cfg.Methods) == 0 {
		cfg.Methods = ConfigDefault.Methods
	}
	return cfg
}
