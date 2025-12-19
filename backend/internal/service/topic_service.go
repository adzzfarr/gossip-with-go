package service

import (
	"fmt"
	"strings"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
)

// TopicService handles business logic related to Topics via the repository layer
type TopicService struct {
	Repo *data.Repository
}

// NewTopicService creates a new instance of TopicService
func NewTopicService(repo *data.Repository) *TopicService {
	return &TopicService{Repo: repo}
}

// GetAllTopics retrieves all topics
func (topicService *TopicService) GetAllTopics() ([]*data.Topic, error) {
	return topicService.Repo.GetAllTopics()
}

// CreateTopic creates a new topic
func (topicService *TopicService) CreateTopic(title, description string, userID int) (*data.Topic, error) {
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
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegate call to repository layer
	topic, err := topicService.Repo.CreateTopic(title, description, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return topic, nil
}

// UpdateTopic updates an existing topic
func (topicService *TopicService) UpdateTopic(topicID int, title, description string, userID int) (*data.Topic, error) {
	// Title Validation
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if len(title) > 200 {
		return nil, fmt.Errorf("title exceeds maximum length of 200 characters")
	}

	// Description Validation
	if strings.TrimSpace(description) == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}

	if len(description) > 1000 {
		return nil, fmt.Errorf("description exceeds maximum length of 1000 characters")
	}

	// UserID Validation
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegeate call to repository layer
	updatedTopic, err := topicService.Repo.UpdateTopic(
		topicID,
		title,
		description,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return updatedTopic, nil
}

// DeleteTopic deletes an existing topic
func (topicService *TopicService) DeleteTopic(topicID, userID int) error {
	// UserID Validation
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	// TopicID Validation
	if topicID <= 0 {
		return fmt.Errorf("invalid topic ID: %d", topicID)
	}

	// Delegate call to repository layer
	err := topicService.Repo.DeleteTopic(topicID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	return nil
}
