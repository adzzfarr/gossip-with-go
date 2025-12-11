package service

import (
	"fmt"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
)

// CommentService handles business logic related to Comments
type CommentService struct {
	// Dependency injection of Repository into service layer
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
