package api

import (
	"net/http"
	"strconv"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	VoteService    *service.VoteService
	PostService    *service.PostService
	CommentService *service.CommentService
}

func NewVoteHandler(voteService *service.VoteService, postService *service.PostService, commentService *service.CommentService) *VoteHandler {
	return &VoteHandler{
		VoteService:    voteService,
		PostService:    postService,
		CommentService: commentService,
	}
}

type VoteRequest struct {
	VoteType int `json:"voteType" binding:"required"` // 1 for upvote, -1 for downvote
}

// VoteOnPost handles POST requests for voting on a post
func (handler *VoteHandler) VoteOnPost(ctx *gin.Context) {
	// Get postID from URL parameter
	postID, err := strconv.Atoi(ctx.Param("postID"))
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid post ID"},
		)
		return
	}

	// Get userID from context (must be authenticated)
	uid, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}
	userID := uid.(int)

	// Parse request body JSON into VoteRequest struct
	var req VoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	// Call service layer
	err = handler.VoteService.VoteOnPost(userID, postID, req.VoteType)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	// Fetch updated post to return current vote count
	updatedPost, err := handler.PostService.GetPostByID(postID, &userID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch updated post"},
		)
		return
	}

	// Return success response with updated vote count
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":   "Vote recorded successfully",
			"voteCount": updatedPost.VoteCount,
			"userVote":  updatedPost.UserVote,
		},
	)
}

// RemoveVoteFromPost handles DELETE requests to remove a user's vote from a post
func (handler *VoteHandler) RemoveVoteFromPost(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("postID"))
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid post ID"},
		)
		return
	}

	uid, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}
	userID := uid.(int)

	err = handler.VoteService.RemoveVoteFromPost(userID, postID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	updatedPost, err := handler.PostService.GetPostByID(postID, &userID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch updated post"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":   "Vote removed successfully",
			"voteCount": updatedPost.VoteCount,
			"userVote":  updatedPost.UserVote,
		},
	)
}

// VoteOnComment handles POST requests for voting on a comment
func (handler *VoteHandler) VoteOnComment(ctx *gin.Context) {
	commentID, err := strconv.Atoi(ctx.Param("commentID"))
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid comment ID"},
		)
		return
	}

	uid, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}
	userID := uid.(int)

	var req VoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	err = handler.VoteService.VoteOnComment(userID, commentID, req.VoteType)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	updatedComment, err := handler.CommentService.GetCommentByID(commentID, &userID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch updated comment"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":   "Vote recorded successfully",
			"voteCount": updatedComment.VoteCount,
			"userVote":  updatedComment.UserVote,
		},
	)
}

// RemoveVoteFromComment handles DELETE requests to remove a user's vote from a comment
func (handler *VoteHandler) RemoveVoteFromComment(ctx *gin.Context) {
	commentID, err := strconv.Atoi(ctx.Param("commentID"))
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid comment ID"},
		)
		return
	}

	uid, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		return
	}
	userID := uid.(int)

	err = handler.VoteService.RemoveVoteFromComment(userID, commentID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	updatedComment, err := handler.CommentService.GetCommentByID(commentID, &userID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch updated comment"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":   "Vote removed successfully",
			"voteCount": updatedComment.VoteCount,
			"userVote":  updatedComment.UserVote,
		},
	)
}
