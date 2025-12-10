package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

// FindUserByUsername fetches user by their unique username
// Used to check if a user exists (during registration) and to retrieve credentials (during login).
func (repo *Repository) FindUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var user User
	query := `
        SELECT user_id, username, password_hash, created_at
        FROM users
        WHERE username = $1`

	err := repo.DB.QueryRow(ctx, query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query to find user failed: %w", err)
	}

	// Return pointer to the found User
	return &user, nil
}

// CreateUser inserts a new user record into the database
// NOTE: The password MUST already be hashed before this function is called (handled by service layer)
func (repo *Repository) CreateUser(user *User) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
        INSERT INTO users (username, password_hash)
        VALUES ($1, $2)
        RETURNING user_id, created_at`

	err := repo.DB.QueryRow(ctx, query, user.Username, user.PasswordHash).Scan(
		&user.UserID,
		&user.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetPostsByTopicID fetches all posts for a given topic ID
func (repo *Repository) GetPostsByTopicID(topicID int) ([]*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		SELECT post_id, topic_id, title, content, created_by, created_at
		FROM posts
		WHERE topic_id = $1
		ORDER BY created_at DESC`

	rows, err := repo.DB.Query(ctx, query, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		var post Post

		err := rows.Scan( // Match order in SELECT statement
			&post.PostID,
			&post.TopicID,
			&post.Title,
			&post.Content,
			&post.CreatedBy,
			&post.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return posts, nil
}
