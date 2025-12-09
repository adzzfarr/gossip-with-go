package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository holds all database access methods
type Repository struct {
	DB *pgxpool.Pool // Connection pool created in ./db.go
}

// NewRepository initializes a new instance of Repository struct
func NewRepository(db *pgxpool.Pool) *Repository {
	// Return pointer to Repository struct with DB field set to provided connection pool
	return &Repository{DB: db}
}

// GetAllTopics fetches all topics from the database
func (repo *Repository) GetAllTopics() ([]*Topic, error) {
	// 1. Context and Deferred Cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensures context is cleaned up when function returns

	// 2. SQL Query
	query := `
        SELECT topic_id, title, description, created_by, created_at
        FROM topics
        ORDER BY created_at DESC`

	// 3. Execute Query
	rows, err := repo.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all topics failed: %w", err)
	}
	defer rows.Close() // Close rows after processing

	// 4. Scan Results
	topics := []*Topic{} // Initialize empty slice of Topic pointers
	for rows.Next() {
		var t Topic

		// Scan column values from current row into fields of the Topic struct
		err := rows.Scan(&t.TopicID, &t.Title, &t.Description, &t.CreatedBy, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning topic row: %w", err)
		}

		topics = append(topics, &t) // Append a pointer to the Topic to the slice
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return topics, nil
}
