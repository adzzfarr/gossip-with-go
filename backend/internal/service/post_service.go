package service

import (
	"fmt"
	"strings"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
)

// PostService handles business logic related to Posts
type PostService struct {
	// Dependency injection of Repository into service layer
	Repo *data.Repository
}

// NewPostService creates a new instance of PostService
func NewPostService(repo *data.Repository) *PostService {
	return &PostService{Repo: repo}
}

// GetPostsByTopicID retrieves all posts for a given topic ID using the repository layer
func (service *PostService) GetPostsByTopicID(topicID int) ([]*data.Post, error) {
	// Validate topicID
	if topicID <= 0 {
		return nil, fmt.Errorf("invalid topic ID: %d", topicID)
	}

	// Delegate call to repository layer
	posts, err := service.Repo.GetPostsByTopicID(topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts for topic ID %d: %w", topicID, err)
	}

	if posts == nil {
		// Return empty slice
		posts = []*data.Post{}
	}

	return posts, nil
}

// CreatePost creates a new post
func (postService *PostService) CreatePost(topicID int, title, content string, createdBy int) (*data.Post, error) {
	// TopicID Validation
	if topicID <= 0 {
		return nil, fmt.Errorf("invalid topic ID: %d", topicID)
	}

	// Title Validation
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if len(title) > 200 {
		return nil, fmt.Errorf("title exceeds maximum length of 200 characters")
	}

	// Content Validation
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}
	if len(content) > 5000 {
		return nil, fmt.Errorf("content exceeds maximum length of 5000 characters")
	}

	// UserID Validation
	if createdBy <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", createdBy)
	}

	// Delegate call to repository layer
	post, err := postService.Repo.CreatePost(topicID, title, content, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// UpdatePost updates an existing post
func (postService *PostService) UpdatePost(postID int, title, content string, userID int) (*data.Post, error) {
	// Title Validation
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if len(title) > 200 {
		return nil, fmt.Errorf("title exceeds maximum length of 200 characters")
	}

	// Content Validation
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	if len(content) > 5000 {
		return nil, fmt.Errorf("content exceeds maximum length of 5000 characters")
	}

	// UserID Validation
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegate call to repository layer
	updatedPost, err := postService.Repo.UpdatePost(postID, title, content, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return updatedPost, nil
}
