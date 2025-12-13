package api

import (
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// LoginHandler handles HTTP requests related to Login
type LoginHandler struct {
	LoginService *service.LoginService
	JWTService   *service.JWTService
}

// NewLoginHandler creates a new instance of LoginHandler
func NewLoginHandler(loginService *service.LoginService, jwtService *service.JWTService) *LoginHandler {
	return &LoginHandler{
		LoginService: loginService,
		JWTService:   jwtService,
	}
}

// LoginRequest struct
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUser handles POST requests for user login
func (handler *LoginHandler) LoginUser(ctx *gin.Context) {
	var req LoginRequest

	// Try to parse request body JSON into LoginRequest struct format
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	// Call service ayer to authenticate user
	user, err := handler.LoginService.Login(req.Username, req.Password)

	if err != nil {
		// Authentication failed
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid username or password"},
		)
		return
	}

	// Generate JWT token
	token, err := handler.JWTService.GenerateToken(user.UserID, user.Username)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
		return
	}

	// Return token to client
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Login successful",
			"token":   token,
		},
	)
}
