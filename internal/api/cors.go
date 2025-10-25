package api

import (
	"net/http"

	"github.com/modelcontextprotocol/registry/internal/config"
)

// CORSMiddleware adds CORS headers to allow cross-origin requests
func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CORS headers if disabled
			if !cfg.CORSEnabled {
				next.ServeHTTP(w, r)
				return
			}

			// Set CORS headers for all requests
			w.Header().Set("Access-Control-Allow-Origin", cfg.CORSAllowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Content-Length")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			// Handle preflight OPTIONS requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
