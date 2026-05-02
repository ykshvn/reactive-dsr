// Package middleware is used for creating middlewares
package middleware

import (
	"net/http"
)

func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i -= 1 {
			next = middlewares[i](next)
		}
		return next
	}
}
