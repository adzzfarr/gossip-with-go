package data

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenDB returns a connection pool to the PostgreSQL database
func OpenDB() (*pgxpool.Pool, error) {
	var dsn string

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		dsn = databaseURL
	} else {
		// Fallback to hardcoded parameters if DATABASE_URL is not set
		dbHost := "localhost"
		dbPort := "5432"
		dbUser := "user"
		dbPass := "password"
		dbName := "forum_db"

		// Construct data source name
		// DSN format for pgx: postgres://username:password@host:port/dbname
		sslMode := "disable"
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName,
			sslMode,
		)
	}

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection pool successfully created.")

	// Return connection pool pointer
	return pool, nil
}
