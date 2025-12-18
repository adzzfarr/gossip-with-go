package service

import (
	"fmt"
	"strings"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
)

// CommentService handles business logic related to comments via the repository layer
type CommentService struct {
	Repo *data.Repository
}

// NewCommentService creates a new instance of CommentService
func NewCommentService(repo *data.Repository) *CommentService {
	return &CommentService{Repo: repo}
}

// GetCommentsByPostID retrieves all comments for a given post
func (commentService *CommentService) GetCommentsByPostID(postID int) ([]*data.Comment, error) {
	// Validate post ID
	if postID <= 0 {
		return nil, fmt.Errorf("invalid post ID: %d", postID)
	}

	// Delegate call to repository layer
	comments, err := commentService.Repo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for post ID %d: %w", postID, err)
	}

	if comments == nil {
		// Return empty slice
		comments = []*data.Comment{}
	}

	return comments, nil
}

// CreateComment creates a new comment on a post
func (commentService *CommentService) CreateComment(postID int, content string, userID int) (*data.Comment, error) {
	// Content Validation
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}
	if len(content) > 2000 {
		return nil, fmt.Errorf("content exceeds maximum length of 2000 characters")
	}

	// Create comment
	createdComment, err := commentService.Repo.CreateComment(postID, content, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return createdComment, nil
}

// UpdateComment updates an existing comment
func (commentService *CommentService) UpdateComment(commentID int, content string, userID int) (*data.Comment, error) {
	// Content Validation
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}
	if len(content) > 2000 {
		return nil, fmt.Errorf("content exceeds maximum length of 2000 characters")
	}

	// UserID Validation
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegate call to repository layer
	updatedComment, err := commentService.Repo.UpdateComment(commentID, content, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return updatedComment, nil
}
