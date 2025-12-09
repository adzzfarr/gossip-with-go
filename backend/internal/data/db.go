package data

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL driver with connection pooling
)

// OpenDB returns a connection pool to the PostgreSQL database
func OpenDB() (*pgxpool.Pool, error) {
	// 1. Connection parameters (Hardcoded for initial testing)
	// **In production, these should be sourced from environment variables or a config file**
	const (
		dbHost = "localhost"
		dbPort = "5432"
		dbUser = "user"
		dbPass = "password"
		dbName = "forum_db"
	)

	// 2. Construct Data Source Name (DSN)
	// DSN format for pgx: postgres://username:password@host:port/dbname
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// 3. Create connection pool (empty context for initial connection)
	// Connection pooling is better than a single connection for handling multiple concurrent requests (e.g. multiple users commenting at the same time)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// 4. Test  connection
	// Ping database to ensure credentials and connection are valid
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection pool successfully created.")

	// Return connection pool pointer and nil error
	return pool, nil
}
