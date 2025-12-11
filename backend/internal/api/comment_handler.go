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
