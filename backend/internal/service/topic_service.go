package service

import (
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
func (service *TopicService) GetAllTopics() ([]*data.Topic, error) {
	// Delegate call to repository layer
	return service.Repo.GetAllTopics()
}
