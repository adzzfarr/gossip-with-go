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
	// (?=.*[A-Z]) requires at least one uppercase letter
	// (?=.*\d) requires at least one digit
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	digitRegex     = regexp.MustCompile(`\d`)
)

// UserService handles business logic related to Users, including hashing of passwords
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

	// Check if username already exists using GetUserByUsername (in repository.go)
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

	// Delegate to the Repository
	if _, err := service.Repo.CreateUser(user); err != nil {
		// Need to check if error is due to a unique constraint violation
		// (i.e. username already taken) for clarity and error translation

		// Try to cast the generic Go error into the specific pgconn.PgError type
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation error code
				// Translate system error into functional error for the API
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
