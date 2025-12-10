package api

import (
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// UserRegistrationRequest defines expected JSON input for new users
type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required"` // returns 400 error if missing
	Password string `json:"password" binding:"required"` // returns 400 error if missing
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
	var req UserRegistrationRequest

	// Try to parse request body JSON into UserRegistrationRequest struct format
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize user object (excluding PasswordHash) into JSON
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
