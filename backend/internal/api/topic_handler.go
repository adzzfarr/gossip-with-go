package api

import (
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// TopicHandler holds instance of TopicService to perform business logic
type TopicHandler struct {
	Service *service.TopicService
}

// NewTopicHandler creates a new instance of TopicHandler
func NewTopicHandler(service *service.TopicService) *TopicHandler {
	return &TopicHandler{Service: service}
}

// GetAllTopics handles GET requests for topics
func (handler *TopicHandler) GetAllTopics(c *gin.Context) {
	// Call service layer
	topics, err := handler.Service.GetAllTopics()

	if err != nil {
		// Log error, send ISE status to client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
		return
	}

	// Gin serializes 'topics' slice into JSON
	c.JSON(http.StatusOK, topics)
}
