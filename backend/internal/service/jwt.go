package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	SecretKey     string
	TokenDuration time.Duration
}

// NewJWTService creates a new instance of JWTService
func NewJWTService(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}
}

// JWTClaims struct
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a user
func (jwtService *JWTService) GenerateToken(userId int, username string) (string, error) {
	now := time.Now()

	claims := &JWTClaims{
		UserID:   userId,   // Private (custom) claim
		Username: username, // Private (custom) claim
		RegisteredClaims: jwt.RegisteredClaims{ // Standard claims
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtService.TokenDuration)),
		},
	}

	// Unsigned token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(jwtService.SecretKey))

	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken verifies JWT token, returns claims if valid
func (jwtService *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token with claims
	token, err := jwt.ParseWithClaims(
		tokenString,  // String to validate
		&JWTClaims{}, // Empty struct to fill with parsed claims
		func(token *jwt.Token) (interface{}, error) { // Key function
			// Verify signing method (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtService.SecretKey), nil
		})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)

	// Check token validity
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
