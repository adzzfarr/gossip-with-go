package api

import (
	"fmt"
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens of incoming requests; authentication required
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from Authorization header
		authHeader := ctx.GetHeader("Authorization")

		// ← ADD THIS LOG
		fmt.Printf("🔍 AuthMiddleware: authHeader = '%s'\n", authHeader)

		if authHeader == "" {
			fmt.Println("❌ AuthMiddleware: No Authorization header")
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header required"},
			)
			ctx.Abort()
			return
		}

		// Expected format: "Bearer <token>"
		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)

		if err != nil || tokenString == "" {
			fmt.Println("❌ AuthMiddleware: Invalid Authorization header format")
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid Authorization header format"},
			)
			ctx.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)

		// ← ADD THIS LOG
		fmt.Printf("🔍 AuthMiddleware: ValidateToken err = %v\n", err)

		if err != nil {
			fmt.Printf("❌ AuthMiddleware: Invalid token: %v\n", err)
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			ctx.Abort()
			return
		}

		fmt.Printf("✅ AuthMiddleware: Valid token for userID = %d\n", claims.UserID)
		// Store claims in context for other handlers
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)

		// Proceed to next handler
		ctx.Next()
	}
}

// OptionalAuthMiddleware processes JWT tokens if present; allows unauthenticated access (guests)
func OptionalAuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			// Continue as guest (no token)
			ctx.Next()
			return
		}

		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil || tokenString == "" {
			// Continue as guest (invalid header)
			ctx.Next()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			// Continue as guest (invalid token)
			ctx.Next()
			return
		}

		// Valid tken; set user info in context
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
