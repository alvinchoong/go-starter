package router

import (
	"log/slog"
	"net/http"
	"time"

	"go-starter/internal/pkg/slogr"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// corsMiddleware configures and returns a CORS (Cross-Origin Resource Sharing) middleware handler.
func corsMiddleware(allowedOrigins []string) func(next http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins:     allowedOrigins,
		AllowOriginFunc:    nil,
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   false,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
	})

	return cors.Handler
}

// requestLogger creates a middleware that logs HTTP request details:
//   - Request path and method
//   - Request ID for tracing
//   - User agent
//   - Response status code
//   - Request duration
//
// It uses structured logging via slog to ensure consistent log format
// and adds the logger to the request context for use by handlers.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		logger := slogr.FromContext(ctx).With(
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		)

		reqID := middleware.GetReqID(ctx)
		if reqID != "" {
			logger = logger.With(slog.String("request-id", reqID))
		}

		logger.Info("START",
			slog.String("user-agent", r.UserAgent()),
		)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r.WithContext(slogr.ToContext(ctx, logger)))

		logger.Info("END",
			slog.Duration("duration", time.Since(start)),
			slog.Int("status", rw.statusCode),
		)
	})
}

// responseWriter is a custom http.ResponseWriter that captures the HTTP status code.
// It embeds the standard http.ResponseWriter and overrides the WriteHeader method to record the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and delegates the call to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
