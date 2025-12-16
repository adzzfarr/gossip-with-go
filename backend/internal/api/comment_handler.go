package api

import (
	"net/http"
	"strconv"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// CommentHandler handles HTTP requests related to Comments
type CommentHandler struct {
	CommentService *service.CommentService
}

// NewCommentHandler creates a new instance of CommentHandler
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}

// GetCommentsByPostID handles GET requests for comments on a specific post
func (handler *CommentHandler) GetCommentsByPostID(ctx *gin.Context) {
	// Get postID from URL parameter
	postIdStr := ctx.Param("postID")
	postID, err := strconv.Atoi(postIdStr)

	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid post ID"})
		return
	}

	// Call service layer
	comments, err := handler.CommentService.GetCommentsByPostID(postID)

	if err != nil {
		// Log error, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch comments for the post"})
		return
	}

	// Gin serializes 'comments' slice into JSON
	ctx.JSON(http.StatusOK, comments)
}

// CreateCommentRequest defines expected JSON input for new comments
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment handles POST requests for creating new comments
func (handler *CommentHandler) CreateComment(ctx *gin.Context) {
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

	// Parse request body JSON into CreateCommentRequest struct
	var req CreateCommentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"},
		)
		return
	}

	// Call service layer to create comment
	comment, err := handler.CommentService.CreateComment(
		postID,
		req.Content,
		userID.(int),
	)

	if err != nil {
		// Check for validation errors
		if err.Error() == "content cannot be empty" ||
			err.Error() == "content exceeds maximum length of 2000 characters" {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		// Otherwise, log error and send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create comment"},
		)
		return
	}

	// Return created comment
	ctx.JSON(http.StatusCreated, comment)
}
