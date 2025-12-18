package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// CommentHandler handles HTTP requests related to comments
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
		// Send ISE status to client
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
		// Check for validation errors (Bad Request 400)
		if err.Error() == "content cannot be empty" ||
			err.Error() == "content exceeds maximum length of 2000 characters" {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		// Otherwise, send ISE status to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create comment"},
		)
		return
	}

	// Return created comment
	ctx.JSON(http.StatusCreated, comment)
}

// UpdateCommentRequest defines expected JSON input for updating comments
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// UpdateComment handles PUT requests for updating existing comments
func (handler *CommentHandler) UpdateComment(ctx *gin.Context) {
	// Get authenticated user's ID from context (set by AuthMiddleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}

	// Get commentID from URL parameter
	commentIDStr := ctx.Param("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid comment ID"},
		)
		return
	}

	// Parse request body JSON into UpdateCommentRequest struct
	var req UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"},
		)
		return
	}

	// Call service layer to update comment
	updatedComment, err := handler.CommentService.UpdateComment(commentID, req.Content, userID.(int))

	if err != nil {
		errMsg := err.Error()

		// Check for not found errors (Not Found 404)
		if strings.Contains(errMsg, "not found") {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "Comment not found"},
			)
			return
		}

		// Check for authorization errors (Forbidden 403)
		if strings.Contains(errMsg, "not authorized") {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{"error": "Not authorized to update this comment"},
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

		// Otherwise, send ISE error to client
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update comment"},
		)
		return
	}

	// Return updated comment
	ctx.JSON(http.StatusOK, updatedComment)
}
