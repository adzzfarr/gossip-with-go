package api

import (
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// TopicHandler holds instance of TopicService to perform business logic
type TopicHandler struct {
	TopicService *service.TopicService
}

// NewTopicHandler creates a new instance of TopicHandler
func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{TopicService: topicService}
}

// GetAllTopics handles GET requests for topics
func (handler *TopicHandler) GetAllTopics(ctx *gin.Context) {
	// Call service layer
	topics, err := handler.TopicService.GetAllTopics()

	if err != nil {
		// Log error, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch topics"})
		return
	}

	// Gin serializes 'topics' slice into JSON
	ctx.JSON(http.StatusOK, topics)
}

// CreateTopicRequest defines expected JSON input for new topics
type CreateTopicRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateTopic handles POST requests for creating new topics
func (handler *TopicHandler) CreateTopic(ctx *gin.Context) {
	// Get authenticated user's ID from context (set by AuthMiddleware)
	userID, exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}

	// Parse request body JSON into CreateTopicRequest struct
	var req CreateTopicRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"},
		)
		return
	}

	// Call service layer to create topic
	topic, err := handler.TopicService.CreateTopic(
		req.Title,
		req.Description,
		userID.(int),
	)

	if err != nil {
		// Log error, send bad request status to client (if validation fails)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	// Return created topic with 201 status
	ctx.JSON(http.StatusCreated, topic)
}
