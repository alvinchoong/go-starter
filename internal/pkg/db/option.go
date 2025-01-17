package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Option is a function that configures a pgxpool.Config
type Option func(*pgxpool.Config)

// WithBeforeConnect sets the BeforeConnect function for the pool
func WithBeforeConnect(beforeConnect func(context.Context, *pgx.ConnConfig) error) Option {
	return func(c *pgxpool.Config) {
		c.BeforeConnect = beforeConnect
	}
}

// WithAfterConnect sets the AfterConnect function for the pool
func WithAfterConnect(afterConnect func(context.Context, *pgx.Conn) error) Option {
	return func(c *pgxpool.Config) {
		c.AfterConnect = afterConnect
	}
}

// WithBeforeAcquire sets the BeforeAcquire function for the pool
func WithBeforeAcquire(beforeAcquire func(context.Context, *pgx.Conn) bool) Option {
	return func(c *pgxpool.Config) {
		c.BeforeAcquire = beforeAcquire
	}
}

// WithAfterRelease sets the AfterRelease function for the pool
func WithAfterRelease(afterRelease func(*pgx.Conn) bool) Option {
	return func(c *pgxpool.Config) {
		c.AfterRelease = afterRelease
	}
}

// WithBeforeClose sets the BeforeClose function for the pool
func WithBeforeClose(beforeClose func(*pgx.Conn)) Option {
	return func(c *pgxpool.Config) {
		c.BeforeClose = beforeClose
	}
}

// WithMaxConnLifetime sets the MaxConnLifetime for the pool
func WithMaxConnLifetime(maxConnLifetime time.Duration) Option {
	return func(c *pgxpool.Config) {
		c.MaxConnLifetime = maxConnLifetime
	}
}

// WithMaxConnLifetimeJitter sets the MaxConnLifetimeJitter for the pool
func WithMaxConnLifetimeJitter(maxConnLifetimeJitter time.Duration) Option {
	return func(c *pgxpool.Config) {
		c.MaxConnLifetimeJitter = maxConnLifetimeJitter
	}
}

// WithMaxConnIdleTime sets the MaxConnIdleTime for the pool
func WithMaxConnIdleTime(maxConnIdleTime time.Duration) Option {
	return func(c *pgxpool.Config) {
		c.MaxConnIdleTime = maxConnIdleTime
	}
}

// WithMaxConns sets the MaxConns for the pool
func WithMaxConns(maxConns int32) Option {
	return func(c *pgxpool.Config) {
		c.MaxConns = maxConns
	}
}

// WithMinConns sets the MinConns for the pool
func WithMinConns(minConns int32) Option {
	return func(c *pgxpool.Config) {
		c.MinConns = minConns
	}
}

// WithHealthCheckPeriod sets the HealthCheckPeriod for the pool
func WithHealthCheckPeriod(healthCheckPeriod time.Duration) Option {
	return func(c *pgxpool.Config) {
		c.HealthCheckPeriod = healthCheckPeriod
	}
}
