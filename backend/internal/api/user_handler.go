package api

import (
	"net/http"
	"strconv"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// UserRegistrationRequest defines expected JSON input for new users
type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserHandler holds UserService instance to perform business logic
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// RegisterUser handles POST requests for user registration
func (handler *UserHandler) RegisterUser(ctx *gin.Context) {
	// Parse request body JSON into UserRegistrationRequest struct format
	var req UserRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	// Call Service Layer
	user, err := handler.UserService.RegisterUser(req.Username, req.Password)

	if err != nil {
		// Since service layer handles input validation (password length, complexity)
		// and unique username checks, errors here are likely client-related (Bad Request 400)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	// Serialize user object (excluding PasswordHash) into JSON
	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"message": "User registered successfully",
			"user":    user,
		},
	)
}

// GetUserByID handles GET requests to fetch user profile by userID
func (handler *UserHandler) GetUserByID(ctx *gin.Context) {
	// Extract userID from URL parameters
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user ID"},
		)
		return
	}

	// Call Service Layer
	user, err := handler.UserService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "User not found"},
			)
			return
		}

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	// Serialize user object (excluding PasswordHash) into JSON
	ctx.JSON(http.StatusOK, user)
}

// GetUserPosts handles GET requests to fetch all posts by a specific user
func (handler *UserHandler) GetUserPosts(ctx *gin.Context) {
	// Extract userID from URL parameters
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user ID"},
		)
		return
	}

	// Call Service Layer
	posts, err := handler.UserService.GetUserPosts(userID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

// GetUserComments handles GET requests to fetch all comments by a specific user
func (handler *UserHandler) GetUserComments(ctx *gin.Context) {
	// Extract userID from URL parameters
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user ID"},
		)
		return
	}

	// Call Service Layer
	comments, err := handler.UserService.GetUserComments(userID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
