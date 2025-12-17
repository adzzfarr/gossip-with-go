package api

import (
	"net/http"
	"strconv"
	"strings"

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
	topicIDStr := ctx.Param("topicId")
	topicID, err := strconv.Atoi(topicIDStr)

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

// CreatePostRequest defines expected JSON input for new posts
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreatePost handles POST requests for creating new posts
func (handler *PostHandler) CreatePost(ctx *gin.Context) {
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
	topicIDStr := ctx.Param("topicId")
	topicID, err := strconv.Atoi(topicIDStr)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid topic ID"})
		return
	}

	// Parse request body JSON into CreatePostRequest struct
	var req CreatePostRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"},
		)
		return
	}

	// Call service layer to create post
	post, err := handler.PostService.CreatePost(
		topicID,
		req.Title,
		req.Content,
		userID.(int),
	)

	if err != nil {
		errMsg := err.Error()

		// Check for validation errors
		if strings.Contains(errMsg, "cannot be empty") ||
			strings.Contains(errMsg, "exceeds maximum length") {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": errMsg},
			)
			return
		}

		// Foreign key constraint failure (topicID does not exist)
		if strings.Contains(errMsg, "foreign key constraint") ||
			strings.Contains(errMsg, "foreign key") {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Topic does not exist"},
			)
			return
		}

		// Otherwise, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create post"},
		)
		return
	}

	// Gin serializes post object into JSON
	ctx.JSON(http.StatusCreated, post)
}
