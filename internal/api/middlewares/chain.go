package middlewares

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

type ChainConfig struct {
	Handler     http.Handler
	Middlewares []Middleware
}

func Chain(config ChainConfig) http.Handler {
	handler := config.Handler
	for _, middleware := range slices.Backward(config.Middlewares) {
		handler = middleware(handler)
	}
	return handler
}
