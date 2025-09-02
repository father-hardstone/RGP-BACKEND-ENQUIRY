package middleware

import (
	"net/http"
	"time"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/utils"
)

// CorsMiddleware handles Cross-Origin Resource Sharing (CORS) headers
// Allows the API to be accessed from different origins (domains)
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all HTTP requests with detailed information
// Provides request tracking and debugging capabilities
func LoggingMiddleware(logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer that captures the status code
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the next handler in the chain
			next.ServeHTTP(lrw, r)

			// Calculate request duration
			duration := time.Since(start)

			// Create log entry
			entry := utils.LogEntry{
				Timestamp:  start,
				Method:     r.Method,
				RemoteAddr: r.RemoteAddr,
				Path:       r.URL.Path,
				Protocol:   r.Proto,
				Duration:   duration,
				StatusCode: lrw.statusCode,
				UserAgent:  r.UserAgent(),
			}

			// Log the request using the enhanced logger
			logger.LogRequest(entry)
		})
	}
}

// loggingResponseWriter is a custom ResponseWriter that captures the status code
// Used by LoggingMiddleware to log response status codes
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// SecurityMiddleware adds basic security headers to responses
// Helps protect against common web vulnerabilities
func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
