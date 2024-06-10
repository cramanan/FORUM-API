package server

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// ServerParameters configures the server parameters
func ServerParameters(handler http.Handler, Request int) *http.Server {
	return &http.Server{
		Addr:              "127.0.0.1:8080",
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Handler:           MiddlewareRateLimiting(handler, Request),
	}
}

// MiddlewareRateLimiting limits the number of requests per second
func MiddlewareRateLimiting(next http.Handler, Request int) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(Request), Request)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
