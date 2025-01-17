package router

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler returns the http handler that handles all requests
func Handler(
	ctx context.Context,
	db *pgxpool.Pool,
	timeout time.Duration,
) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Top-level middlewares
	r.Use(middleware.RequestID)
	r.Use(requestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout))
	r.Use(corsMiddleware([]string{"*"}))

	r.Get("/ping", pingHandler)

	return r, nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
