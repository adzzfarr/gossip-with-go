package api

import (
	"net/http"
	"strconv"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// PostHandler handles HTTP requests related to Posts
type PostHandler struct {
	PostService *service.PostService
}

// NewPostHandler creates a new instance of PostHandler
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{PostService: postService}
}

// GetPostsByTopicID handles GET requests for posts in a specific topic
func (handler *PostHandler) GetPostsByTopicID(ctx *gin.Context) {
	// Get topicID from URL parameter
	topicIdStr := ctx.Param("topicId")
	topicID, err := strconv.Atoi(topicIdStr)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid topic ID"})
		return
	}

	// Call service layer
	posts, err := handler.PostService.GetPostsByTopicID((topicID))

	if err != nil {
		// Log error, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch posts for the topic"})
		return
	}

	// Gin serializes 'posts' slice into JSON
	ctx.JSON(http.StatusOK, posts)
}
