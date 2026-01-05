// backend/internal/service/user_service.go
package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"

	"github.com/jackc/pgx/v5/pgconn" // PostgreSQL driver for error handling
	"golang.org/x/crypto/bcrypt"
)

var (
	// (?=.*[a-z]) requires at least one lowercase letter
	lowercaseRegex = regexp.MustCompile(`[a-z]`)

	// (?=.*[A-Z]) requires at least one uppercase letter
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)

	// (?=.*\d) requires at least one digit
	digitRegex = regexp.MustCompile(`\d`)
)

// UserService handles business logic related to Users (*** including hashing of passwords ***) via the repository layer
type UserService struct {
	Repo *data.Repository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo *data.Repository) *UserService {
	return &UserService{Repo: repo}
}

// RegisterUser handles password hashing and delegation to the Repository
func (service *UserService) RegisterUser(username, password string) (*data.User, error) {
	// Input Validation
	if len(password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	if !lowercaseRegex.MatchString(password) {
		return nil, fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !uppercaseRegex.MatchString(password) {
		return nil, fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !digitRegex.MatchString(password) {
		return nil, fmt.Errorf("password must contain at least one digit")
	}

	// Check if username already exists
	existingUser, err := service.Repo.GetUserByUsername(username)

	if err != nil {
		// Check if "user not found" error (i.e. username is available)
		if err.Error() != fmt.Sprintf("user not found: %s", username) {
			// Some other database error occurred
			return nil, fmt.Errorf("error checking existing user: %w", err)
		}

		// Username not found means username is available => proceed
	} else if existingUser != nil {
		// User was found means username is taken
		return nil, fmt.Errorf("username '%s' is already taken", username)
	}

	// Hash Password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create User
	user := &data.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	// Delegate to the repository layer
	if _, err := service.Repo.CreateUser(user); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation error code
				// unique constraint violation means username is taken
				if strings.Contains(pgErr.ConstraintName, "username") {
					return nil, fmt.Errorf("username is already taken")
				}
			}
		}

		// Otherwise, return generic error
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	// Clear hash before returning the user to the client
	user.PasswordHash = ""
	return user, nil
}

// GetUserByID retrieves a user by their ID
func (service *UserService) GetUserByID(userID int) (*data.User, error) {
	// UserID Validation
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegate call to repository layer
	user, err := service.Repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", userID, err)
	}

	return user, nil
}

// GetUserPosts retrieves all posts created by a specific user
func (service *UserService) GetUserPosts(userID int) ([]*data.Post, error) {
	// UserID Validation
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegate call to repository layer
	posts, err := service.Repo.GetUserPosts(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts for user ID %d: %w", userID, err)
	}

	return posts, nil
}

// GetUserComments retrieves all comments made by a specific user
func (service *UserService) GetUserComments(userID int) ([]*data.Comment, error) {
	// UserID Validation
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	// Delegate call to repository layer
	comments, err := service.Repo.GetUserComments(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for user ID %d: %w", userID, err)
	}

	return comments, nil
}
