package data

import "time"

// json and db tags are used for serialisation and database mapping respectively

// User struct
type User struct {
	UserID       int       `json:"userID" db:"user_id"` // Primary key
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
	Username    string    `json:"username" db:"username"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// Post struct
type Post struct {
	PostID     int       `json:"postID" db:"post_id"`   // Primary key
	TopicID    int       `json:"topicID" db:"topic_id"` // Foreign key to Topic
	TopicTitle string    `json:"topicTitle" db:"topic_title"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	CreatedBy  int       `json:"createdBy" db:"created_by"`
	Username   string    `json:"username" db:"username"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	VoteCount  int       `json:"voteCount" db:"vote_count"`
	UserVote   *int      `json:"userVote,omitempty" db:"user_vote"` // Current user's vote on post
}

// Comment struct
type Comment struct {
	CommentID int       `json:"commentID" db:"comment_id"` // Primary key
	PostID    int       `json:"postID" db:"post_id"`       // Foreign key to Post
	PostTitle string    `json:"postTitle" db:"post_title"`
	Content   string    `json:"content" db:"content"`
	CreatedBy int       `json:"createdBy" db:"created_by"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	VoteCount int       `json:"voteCount" db:"vote_count"`
	UserVote  *int      `json:"userVote,omitempty" db:"user_vote"` // Current user's vote on comment
}

// Vote struct
type Vote struct {
	VoteID    int       `json:"voteID" db:"vote_id"`                 // Primary key
	UserID    int       `json:"userID" db:"user_id"`                 // Foreign key to User
	PostID    *int      `json:"postID,omitempty" db:"post_id"`       // Foreign key to Post (nullable); vote can be for either post or comment, not both
	CommentID *int      `json:"commentID,omitempty" db:"comment_id"` // Foreign key to Comment (nullable)
	VoteType  int       `json:"voteType" db:"vote_type"`             // +1 for upvote, -1 for downvote
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
