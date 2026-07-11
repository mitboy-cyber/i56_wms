// Package database provides PostgreSQL connection pool management.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/i56/framework/core/config"
)

// DB wraps sql.DB with helper methods.
type DB struct {
	*sql.DB
	driver string
}

// Open creates a new database connection pool.
func Open(cfg config.DatabaseConfig) (*DB, error) {
	dsn := buildDSN(cfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database: open: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	return &DB{DB: db, driver: cfg.Driver}, nil
}

// OpenSQLite creates an in-memory SQLite database for testing.
func OpenSQLite(path string) (*DB, error) {
	if path == "" {
		path = ":memory:"
	}
	dsn := path + "?_journal_mode=WAL&_foreign_keys=on"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("database: open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite serialized
	return &DB{DB: db, driver: "sqlite3"}, nil
}

// Driver returns the database driver name.
func (db *DB) Driver() string { return db.driver }

// IsPostgres returns true if the driver is postgres.
func (db *DB) IsPostgres() bool { return db.driver == "postgres" }

// IsSQLite returns true if the driver is sqlite3.
func (db *DB) IsSQLite() bool { return db.driver == "sqlite3" }

// RunMigrations executes SQL migration files in order.
func (db *DB) RunMigrations(ctx context.Context, migrations []string) error {
	for i, m := range migrations {
		if _, err := db.ExecContext(ctx, m); err != nil {
			return fmt.Errorf("database: migration %d: %w", i+1, err)
		}
	}
	return nil
}

// InTransaction runs a function within a database transaction.
func (db *DB) InTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("database: begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
}
