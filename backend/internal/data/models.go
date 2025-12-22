package data

import "time"

// json and db tags are used for serialisation and database mapping respectively

// User struct
type User struct {
	UserID       int       `json:"userID" db:"userID"` // Primary key
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"` // Exclude from JSON output for security
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// Topic struct
type Topic struct {
	TopicID     int       `json:"topicID" db:"topic_id"` // Primary key
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedBy   int       `json:"createdBy" db:"created_by"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// Post struct
type Post struct {
	PostID    int       `json:"postID" db:"post_id"`   // Primary key
	TopicID   int       `json:"topicID" db:"topic_id"` // Foreign key to Topic
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedBy int       `json:"createdBy" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Comment struct
type Comment struct {
	CommentID int       `json:"commentID" db:"comment_id"` // Primary key
	PostID    int       `json:"postID" db:"post_id"`       // Foreign key to Post
	Content   string    `json:"content" db:"content"`
	CreatedBy int       `json:"createdBy" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
