package router

import (
	"context"
	"net/http"
	"time"

	"go-starter/internal/models"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/plainresp"
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

	q := models.New()

	// Post CRUD API
	NewPostHandler(db, q).Mount(r)

	// User API
	NewUserHandler(newHTTPClient(), "https://jsonplaceholder.typicode.com/users").Mount(r)

	r.Get("/ping", httphandler.Handle(pingHandler))

	return r, nil
}

func pingHandler(_ *http.Request) httphandler.Responder {
	return plainresp.Success("pong")
}

// newHTTPClient returns a new HTTP client
// by default it is configured
// - with a 15-second timeout to prevent requests from hanging indefinitely
// - not follow redirects automatically to prevent redirect-related attacks
func newHTTPClient(opts ...func(*http.Client)) *http.Client {
	c := &http.Client{
		Transport: http.DefaultTransport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     nil,
		Timeout: 15 * time.Second,
	}
	for _, o := range opts {
		o(c)
	}

	return c
}
