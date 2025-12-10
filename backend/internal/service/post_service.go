package service

import (
	"fmt"

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
