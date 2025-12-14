package service

import (
	"fmt"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
)

// TopicService handles business logic related to Topics
type TopicService struct {
	// Dependency injection of Repository into service layer
	Repo *data.Repository
}

// NewTopicService creates a new instance of TopicService
func NewTopicService(repo *data.Repository) *TopicService {
	return &TopicService{Repo: repo}
}

// GetAllTopics retrieves all topics using the repository layer
func (topicService *TopicService) GetAllTopics() ([]*data.Topic, error) {
	// Delegate call to repository layer
	return topicService.Repo.GetAllTopics()
}

// CreateTopic creates a new topic
func (topicService *TopicService) CreateTopic(title, description string, createdBy int) (*data.Topic, error) {
	// Title Validation
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if len(title) > 200 {
		return nil, fmt.Errorf("title cannot exceed 200 characters")
	}

	// Description Validation
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}

	if len(description) > 1000 {
		return nil, fmt.Errorf("description cannot exceed 1000 characters")
	}

	// UserID Validation
	if createdBy <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", createdBy)
	}

	// Delegate call to repository layer
	topic, err := topicService.Repo.CreateTopic(title, description, createdBy)

	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return topic, nil
}
