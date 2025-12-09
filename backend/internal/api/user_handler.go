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
	Service *service.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// RegisterUser handles POST requests for user registration
func (handler *UserHandler) RegisterUser(c *gin.Context) {
	var req UserRegistrationRequest

	// Try to parse request body JSON into UserRegistrationRequest struct format
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	// Call Service Layer
	user, err := handler.Service.RegisterUser(req.Username, req.Password)

	if err != nil {
		// Since service layer handles input validation (password length, complexity)
		// and unique username checks, errors here are likely client-related (Bad Request 400)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize user object (excluding PasswordHash) into JSON
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
