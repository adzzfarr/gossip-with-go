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
	// Context and Deferred Cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensures context is cleaned up when function returns

	// SQL Query
	query := `
        SELECT topic_id, title, description, created_by, created_at, updated_at
        FROM topics
        ORDER BY created_at DESC`

	// Execute Query
	rows, err := repo.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all topics failed: %w", err)
	}
	defer rows.Close() // Close rows after processing

	// Scan Results
	topics := []*Topic{} // Initialize empty slice of Topic pointers
	for rows.Next() {
		var t Topic

		// Scan column values from current row into fields of the Topic struct
		err := rows.Scan(
			&t.TopicID,
			&t.Title,
			&t.Description,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
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

// GetUserByUsername fetches user by their unique username
// Used to check if a user exists (during registration) and to retrieve credentials (during login).
func (repo *Repository) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var user User
	query := `
        SELECT user_id, username, password_hash, created_at, updated_at
        FROM users
        WHERE username = $1`

	err := repo.DB.QueryRow(ctx, query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", username)
		}
		return nil, fmt.Errorf("query to find user failed: %w", err)
	}

	// Return pointer to the found User
	return &user, nil
}

// CreateUser inserts a new user record into the database
// NOTE: The password MUST already be hashed before this function is called (handled by service layer)
func (repo *Repository) CreateUser(user *User) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
        INSERT INTO users (username, password_hash, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING user_id, created_at, updated_at`

	err := repo.DB.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
	).Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// CreateTopic inserts a new topic into the database
func (repo *Repository) CreateTopic(title, description string, createdBy int) (*Topic, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO topics (title, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING topic_id, title, description, created_by, created_at, updated_at
	`

	// Scan returned row into Topic struct
	var topic Topic
	err := repo.DB.QueryRow(
		ctx,
		query,
		title,
		description,
		createdBy,
	).Scan(
		&topic.TopicID,
		&topic.Title,
		&topic.Description,
		&topic.CreatedBy,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return &topic, nil
}

// GetPostsByTopicID fetches all posts for a given topic ID
func (repo *Repository) GetPostsByTopicID(topicID int) ([]*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		SELECT post_id, topic_id, title, content, created_by, created_at, updated_at
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
			&post.UpdatedAt,
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

// GetCommentsByPostID fetches all comments for a given post ID
func (repo *Repository) GetCommentsByPostID(postID int) ([]*Comment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		SELECT comment_id, post_id, content, created_by, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at DESC`

	rows, err := repo.DB.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	comments := []*Comment{}
	for rows.Next() {
		var comment Comment

		err := rows.Scan( // Match order in SELECT statement
			&comment.CommentID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedBy,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return comments, nil
}

// CreatePost inserts a new post into the database
func (repo *Repository) CreatePost(topicID int, title, content string, createdBy int) (*Post, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING post_id, topic_id, title, content, created_by, created_at, updated_at`

	var post Post
	err := repo.DB.QueryRow(
		ctx,
		query,
		topicID,
		title,
		content,
		createdBy,
	).Scan(
		&post.PostID,
		&post.TopicID,
		&post.Title,
		&post.Content,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &post, nil
}

// CreateComment inserts a new comment into the database
func (repo *Repository) CreateComment(postID int, content string, createdBy int) (*Comment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO comments (post_id, content, created_by)
		VALUES ($1, $2, $3)
		RETURNING comment_id, post_id, content, created_by, created_at, updated_at`

	var comment Comment
	err := repo.DB.QueryRow(
		ctx,
		query,
		postID,
		content,
		createdBy,
	).Scan(
		&comment.CommentID,
		&comment.PostID,
		&comment.Content,
		&comment.CreatedBy,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

/*
// UpdateTopic updates an existing topic's title and description
func (repo *Repository) UpdateTopic(topicID int, title, description string, userID int) (*Topic, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		UPDATE topics
		SET title = $1, description = $2, updated_at = NOW()
		WHERE topic_id = $3 AND created_by = $4
		RETURNING topic_id, title, description, created_by, created_at, updated_at`

	var topic Topic
	err := repo.DB.QueryRow(
		ctx,
		query,
		title,
		description,
		topicID,
		userID,
	).Scan(
		&topic.TopicID,
		&topic.Title,
		&topic.Description,
		&topic.CreatedBy,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)

}

*/
