package limiter

import "github.com/ximispot/woody"

const (
	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

type LimiterHandler interface {
	New(config Config) woody.Handler
}

// New creates a new middleware handler
func New(config ...Config) woody.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return the specified middleware handler.
	return cfg.LimiterMiddleware.New(cfg)
}
