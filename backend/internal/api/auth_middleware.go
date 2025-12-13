package api

import (
	"fmt"
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens of incoming requests
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from Authorization header
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header required"},
			)
			ctx.Abort()
			return
		}

		// Try to extract token string
		// Expected format: "Bearer <token>"
		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)

		if err != nil || tokenString == "" {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid Authorization header format"},
			)
			ctx.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)

		if err != nil {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			ctx.Abort()
			return
		}

		// Store claims in context for further handlers
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)

		// Proceed to next handler
		ctx.Next()
	}
}
