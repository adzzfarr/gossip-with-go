package data

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL driver with connection pooling
)

// OpenDB returns a connection pool to the PostgreSQL database
func OpenDB() (*pgxpool.Pool, error) {
	// 1. Connection parameters (Hardcoded)
	// In production, get from environment variables
	const (
		dbHost = "localhost"
		dbPort = "5432"
		dbUser = "user"
		dbPass = "password"
		dbName = "forum_db"
	)

	// 2. Construct data source name
	// DSN format for pgx: postgres://username:password@host:port/dbname
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// 3. Create connection pool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// 4. Test  connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection pool successfully created.")

	// Return connection pool pointer
	return pool, nil
}
