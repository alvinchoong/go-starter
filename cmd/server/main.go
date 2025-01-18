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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if err := errmain(ctx); err != nil {
		cancel()
		slog.Error("something happen", slog.Any("err", err))
		os.Exit(1)
	}
}

func errmain(ctx context.Context) error {
	// parse env vars
	flags, err := newFlags()
	if err != nil {
		return fmt.Errorf("newFlags: %w", err)
	}

	// Create and set logger with build info in context
	slogr.SetDefaultJSON(flags.logLevel)
	logger := slog.Default().With(
		slog.String("version", buildinfo.Version),
		slog.String("build-time", buildinfo.BuildTime),
	)
	ctx = slogr.ToContext(ctx, logger)

	// connect to DB
	db, err := db.Connect(ctx, flags.databaseURL,
		db.WithMaxConnIdleTime(flags.databaseIdleConnTimeout),
		db.WithMinConns(flags.databaseConns),
		db.WithMaxConns(flags.databaseConns))
	if err != nil {
		return fmt.Errorf("db.Connect: %w", err)
	}

	handler, err := router.Handler(ctx, db, flags.serverReadTimeout+flags.serverWriteTimeout)
	if err != nil {
		return fmt.Errorf("router.Handler: %w", err)
	}

	// create a new server
	server := &http.Server{
		Addr:                         flags.serverAddr,
		Handler:                      handler,
		ReadHeaderTimeout:            10 * time.Second,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  flags.serverReadTimeout,
		WriteTimeout:                 flags.serverWriteTimeout,
		IdleTimeout:                  flags.serverIdleTimeout,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	err = runHTTPServer(ctx, server)
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("runHTTPServer: %w", err)
	}

	return nil
}

func runHTTPServer(ctx context.Context, server *http.Server) error {
	logger := slogr.FromContext(ctx)

	// Create an errgroup
	g, gctx := errgroup.WithContext(ctx)

	// Start the server in a goroutine
	g.Go(func() error {
		logger.Info("[server] starting...", slog.String("addr", server.Addr))

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

		// Create a new shutdown context with a timeout
		// shutdownCtx, shutdownCancel := cmd.NewShutdownContext(ctx, 10*time.Second)
		// defer shutdownCancel()
		defer server.Close()

		// Shutdown the server gracefully
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	})

	// Wait for all goroutines to finish and return any errors encountered
	if err := g.Wait(); err != nil {
		return fmt.Errorf("g.Wait: %w", err)
	}

	return nil
}

type flags struct {
	logLevel slog.Level

	// database
	databaseURL             string
	databaseIdleConnTimeout time.Duration
	databaseConns           int32

	// http.Server
	serverAddr         string
	serverReadTimeout  time.Duration
	serverWriteTimeout time.Duration
	serverIdleTimeout  time.Duration
}

func newFlags() (flags, error) {
	var logLevel slog.Level
	if err := envvar.UnmarshalText("LOG_LEVEL", &logLevel); err != nil {
		return flags{}, fmt.Errorf("fail to parse LOG_LEVEL: %w", err)
	}

	dbConns, err := envvar.ParseInt32("DATABASE_CONNS")
	if err != nil {
		return flags{}, fmt.Errorf("fail to parse DATABASE_MAX_CONNS: %w", err)
	}

	readTimeout, err := envvar.ParseDuration("SERVER_READ_TIMEOUT")
	if err != nil {
		return flags{}, fmt.Errorf("fail to parse SERVER_READ_TIMEOUT: %w", err)
	}
	writeTimeout, err := envvar.ParseDuration("SERVER_WRITE_TIMEOUT")
	if err != nil {
		return flags{}, fmt.Errorf("fail to parse SERVER_WRITE_TIMEOUT: %w", err)
	}
	idleTimeout, err := envvar.ParseDuration("SERVER_IDLE_TIMEOUT")
	if err != nil {
		return flags{}, fmt.Errorf("fail to parse SERVER_IDLE_TIMEOUT: %w", err)
	}

	return flags{
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
