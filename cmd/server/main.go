package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-starter/cmd/server/router"
	"go-starter/internal/pkg/buildinfo"
	"go-starter/internal/pkg/db"
	"go-starter/internal/pkg/envvar"
	"go-starter/internal/pkg/slogr"

	"golang.org/x/sync/errgroup"
)

func main() {
	// Setup context with cancellation on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if err := run(ctx); err != nil {
		cancel()
		slog.Error("something happen", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Load and validate configuration from environment variables
	config, err := newConfig()
	if err != nil {
		return fmt.Errorf("newConfig: %w", err)
	}

	// Initialize structured logger with build information
	slogr.SetDefaultJSON(config.logLevel)
	logger := slog.Default().With(
		slog.String("version", buildinfo.Version),
		slog.String("build-time", buildinfo.BuildTime),
	)
	ctx = slogr.ToContext(ctx, logger)

	// Initialize database connection with configured parameters
	db, err := db.Connect(ctx, config.databaseURL,
		db.WithMaxConnIdleTime(config.databaseIdleConnTimeout),
		db.WithMinConns(config.databaseConns),
		db.WithMaxConns(config.databaseConns))
	if err != nil {
		return fmt.Errorf("db.Connect: %w", err)
	}

	// Setup HTTP router with configured timeout
	handler, err := router.Handler(ctx, db, config.serverReadTimeout+config.serverWriteTimeout)
	if err != nil {
		return fmt.Errorf("router.Handler: %w", err)
	}

	// Initialize HTTP server
	server := &http.Server{
		Addr:                         config.serverAddr,
		Handler:                      handler,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  config.serverReadTimeout,
		ReadHeaderTimeout:            10 * time.Second,
		WriteTimeout:                 config.serverWriteTimeout,
		IdleTimeout:                  config.serverIdleTimeout,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
		HTTP2:                        nil,
		Protocols:                    nil,
	}

	// Start server and handle graceful shutdown
	err = runHTTPServer(ctx, server)
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("runHTTPServer: %w", err)
	}

	return nil
}

// runHTTPServer starts the HTTP server and handles graceful shutdown.
// It uses errgroup to manage concurrent operations and ensure proper cleanup.
func runHTTPServer(ctx context.Context, server *http.Server) error {
	logger := slogr.FromContext(ctx)

	// Create errgroup for managing server goroutine
	g, gctx := errgroup.WithContext(ctx)

	// Start server in a goroutine
	g.Go(func() error {
		logger.Info("starting HTTP server", slog.String("addr", server.Addr))

		// ListenAndServe returns ErrServerClosed on graceful shutdown
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server.ListenAndServe: %w", err)
		}
		return nil
	})

	// Setup signal handling in another goroutine
	g.Go(func() error {
		// Wait for the context to be canceled e.g., via signal.NotifyContext
		<-gctx.Done()

		logger.Info("[server] shutting down...")

		defer server.Close()

		// Initiate graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown the server gracefully
		//nolint:contextcheck // We need a fresh context here since the parent context is cancelled
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("errgroup.Wait: %w", err)
	}

	return nil
}

// config holds all the configuration options for the server.
// These are typically set through environment variables.
type config struct {
	// logLevel determines the minimum severity level for log output
	logLevel slog.Level

	// Database configuration
	databaseURL             string        // PostgreSQL connection URL
	databaseIdleConnTimeout time.Duration // Maximum time a connection can be idle
	databaseConns           int32         // Maximum number of database connections

	// HTTP server configuration
	serverAddr         string        // Address to listen on (e.g., ":8080")
	serverReadTimeout  time.Duration // Maximum duration for reading entire request
	serverWriteTimeout time.Duration // Maximum duration for writing response
	serverIdleTimeout  time.Duration // Maximum duration for idle keep-alive connections
}

// newConfig loads and validates configuration from environment variables.
// It returns an error if any required variables are missing or invalid.
func newConfig() (config, error) {
	var logLevel slog.Level
	if err := envvar.UnmarshalText("LOG_LEVEL", &logLevel); err != nil {
		return config{}, fmt.Errorf("fail to parse LOG_LEVEL: %w", err)
	}

	dbConns, err := envvar.ParseInt32("DATABASE_CONNS")
	if err != nil {
		return config{}, fmt.Errorf("fail to parse DATABASE_MAX_CONNS: %w", err)
	}

	readTimeout, err := envvar.ParseDuration("SERVER_READ_TIMEOUT")
	if err != nil {
		return config{}, fmt.Errorf("fail to parse SERVER_READ_TIMEOUT: %w", err)
	}
	writeTimeout, err := envvar.ParseDuration("SERVER_WRITE_TIMEOUT")
	if err != nil {
		return config{}, fmt.Errorf("fail to parse SERVER_WRITE_TIMEOUT: %w", err)
	}
	idleTimeout, err := envvar.ParseDuration("SERVER_IDLE_TIMEOUT")
	if err != nil {
		return config{}, fmt.Errorf("fail to parse SERVER_IDLE_TIMEOUT: %w", err)
	}

	return config{
		logLevel:                logLevel,
		databaseURL:             os.Getenv("DATABASE_URL"),
		databaseIdleConnTimeout: 30 * time.Minute,
		databaseConns:           dbConns,
		serverAddr:              os.Getenv("SERVER_ADDR"),
		serverReadTimeout:       readTimeout,
		serverWriteTimeout:      writeTimeout,
		serverIdleTimeout:       idleTimeout,
	}, nil
}
