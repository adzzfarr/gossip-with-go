package api

import (
	"net/http"
	"strconv"
	"strings"

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

	// Return created topic
	ctx.JSON(http.StatusCreated, topic)
}

// UpdateTopicRequest defines expected JSON input for updating topics
type UpdateTopicRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateTopic handles PUT requests for updating existing topics
func (handler *TopicHandler) UpdateTopic(ctx *gin.Context) {
	// Get authenticated user's ID from context (set by AuthMiddleware)
	userID, exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}

	// Get topicID from URL parameter
	topicIDStr := ctx.Param("topicID")
	topicID, err := strconv.Atoi(topicIDStr)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid topic ID"})
		return
	}

	// Parse request body JSON into UpdateTopicRequest struct
	var req UpdateTopicRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"},
		)
		return
	}

	// Call service layer to update topic
	updatedTopic, err := handler.TopicService.UpdateTopic(
		topicID,
		req.Title,
		req.Description,
		userID.(int),
	)

	if err != nil {
		errMsg := err.Error()

		// Check for not found errors (Not Found 404)
		if strings.Contains(errMsg, "not found") {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": errMsg},
			)
			return
		}

		// Check for authorization errors (Forbidden 403)
		if strings.Contains(errMsg, "not authorized") {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{"error": errMsg},
			)
			return
		}

		// Check for validation errors (Bad Request 400)
		if strings.Contains(errMsg, "cannot be empty") ||
			strings.Contains(errMsg, "exceeds maximum length") {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": errMsg},
			)
			return
		}

		// Otherwise, return ISE
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update topic"},
		)
		return
	}

	// Return updated topic
	ctx.JSON(http.StatusOK, updatedTopic)
}
