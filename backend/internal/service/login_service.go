package service

import (
	"fmt"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
	"golang.org/x/crypto/bcrypt"
)

// LoginService handles business logic related to login via the repository layer
type LoginService struct {
	Repo *data.Repository
}

// NewLoginService creates a new instance of LoginService
func NewLoginService(repo *data.Repository) *LoginService {
	return &LoginService{Repo: repo}
}

// Login authenticates a user with given username and password
func (loginService *LoginService) Login(username, password string) (*data.User, error) {
	// Validate input
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	// Delegate call to repository layer
	user, err := loginService.Repo.GetUserByUsername(username)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check password against stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
