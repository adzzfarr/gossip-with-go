package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// PostHandler handles HTTP requests related to posts
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
	topicIDStr := ctx.Param("topicID")
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
		// Send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch posts for the topic"})
		return
	}

	// Gin serializes 'posts' slice into JSON
	ctx.JSON(http.StatusOK, posts)
}

// GetPostByID handles GET requests for a specific post by its ID
func (handler *PostHandler) GetPostByID(ctx *gin.Context) {
	// Get postID from URL parameter
	postID, err := strconv.Atoi(ctx.Param("postID"))
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid post ID"},
		)
		return
	}

	post, err := handler.PostService.GetPostByID(postID)
	if err != nil {
		// Check for not found errors (Not Found 404)
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "Post not found"},
			)
			return
		}

		// Otherwise, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch post"},
		)
		return
	}

	// Gin serializes 'post' object into JSON
	ctx.JSON(http.StatusOK, post)
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
	topicIDStr := ctx.Param("topicID")
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

		// Check for validation errors (Bad Request 400)
		if strings.Contains(errMsg, "cannot be empty") ||
			strings.Contains(errMsg, "exceeds maximum length") {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": errMsg},
			)
			return
		}

		// Foreign key constraint failure (i.e. topicID does not exist)
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

// UpdatePostRequest defines expected JSON input for updating posts
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UpdatePost handles PUT requests for updating existing posts
func (handler *PostHandler) UpdatePost(ctx *gin.Context) {
	// Get authenticated user's ID from context (set by AuthMiddleware)
	userID, exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}

	// Get postID from URL parameter
	postIDStr := ctx.Param("postID")
	postID, err := strconv.Atoi(postIDStr)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid post ID"})
		return
	}

	// Parse request body JSON into UpdatePostRequest struct
	var req UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"},
		)
		return
	}

	// Call service layer to update post
	updatedPost, err := handler.PostService.UpdatePost(
		postID,
		req.Title,
		req.Content,
		userID.(int),
	)

	if err != nil {
		errMsg := err.Error()

		// Check for not found errors (Not Found 404)
		if strings.Contains(errMsg, "not found") {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "Post not found"},
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

		// Otherwise, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update post"},
		)
		return
	}

	// Gin serializes updatedPost object into JSON
	ctx.JSON(http.StatusOK, updatedPost)
}

// DeletePost handles DELETE requests for deleting existing posts
func (handler *PostHandler) DeletePost(ctx *gin.Context) {
	// Get authenticated user's ID from context (set by AuthMiddleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}

	// Get postID from URL parameter
	postIDStr := ctx.Param("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid post ID"},
		)
		return
	}

	// Call service layer to delete post
	err = handler.PostService.DeletePost(postID, userID.(int))
	if err != nil {
		errMsg := err.Error()

		// Check for not found errors (Not Found 404)
		if strings.Contains(errMsg, "not found") {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "Post not found"},
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
		if strings.Contains(errMsg, "invalid user ID") ||
			strings.Contains(errMsg, "invalid post ID") {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": errMsg},
			)
			return
		}

		// Otherwise, send ISE error to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to delete post"},
		)
		return
	}

	// Return No Content status on successful deletion
	ctx.Status(http.StatusNoContent)
}
