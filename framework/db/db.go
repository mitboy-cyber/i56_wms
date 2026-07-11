// Package db provides a PostgreSQL connection pool for the I56 framework.
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is the shared PostgreSQL connection pool (nil when DB is unavailable).
var Pool *pgxpool.Pool

// Connect creates the connection pool. Call once at startup.
func Connect(dsn string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Printf("[DB] PostgreSQL connect failed (falling back to in-memory): %v", err)
		return err
	}
	if err := Pool.Ping(context.Background()); err != nil {
		log.Printf("[DB] PostgreSQL ping failed (falling back to in-memory): %v", err)
		Pool.Close()
		Pool = nil
		return err
	}
	log.Println("[DB] PostgreSQL connected")
	return nil
}

// Close shuts down the pool cleanly.
func Close() {
	if Pool != nil {
		Pool.Close()
		Pool = nil
	}
}
