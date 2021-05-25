//go:generate mockery --name RouterOption --output ../mocks

package options

import (
	"time"

	"github.com/labstack/echo/v4"
)

type (
	RouterOption interface {
		Apply(s *RouterSettings)
	}

	RouterSettings struct {
		Middlewares           []echo.MiddlewareFunc
		CustomCacheExpiration *time.Duration
		NoCache               bool
	}
)

// WithMiddlewares returns a RouterOption that specifies an array of echo.Middleware
func WithMiddlewares(middlewares ...echo.MiddlewareFunc) RouterOption {
	return withMiddlewares{middlewares}
}

type withMiddlewares struct{ middlewares []echo.MiddlewareFunc }

func (w withMiddlewares) Apply(o *RouterSettings) {
	o.Middlewares = w.middlewares
}

// WithCustomCacheExpiration returns a RouterOption that specifies custom expiration value for cache
func WithCustomCacheExpiration(cacheExpiration time.Duration) RouterOption {
	return withCustomCacheExpiration{cacheExpiration}
}

type withCustomCacheExpiration struct{ customCacheExpiration time.Duration }

func (w withCustomCacheExpiration) Apply(o *RouterSettings) {
	o.CustomCacheExpiration = &w.customCacheExpiration
}

// WithNoCache returns a RouterOption that disable cache
func WithNoCache() RouterOption {
	return withNoCache{}
}

type withNoCache struct{}

func (w withNoCache) Apply(o *RouterSettings) {
	o.NoCache = true
}

func ApplyOptions(options ...RouterOption) *RouterSettings {
	rs := &RouterSettings{}

	for _, option := range options {
		option.Apply(rs)
	}

	return rs
}
