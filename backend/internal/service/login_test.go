// Run `go test -v ./internal/service -run TestLoginService` in /backend
package service

import (
	"context"
	"testing"
	"time"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginService(t *testing.T) {
	// Set up database connection
	dbPool, err := data.OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	repo := data.NewRepository(dbPool)
	loginService := NewLoginService(repo)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "testuser"
	testPassword := "testpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = repo.DB.Exec(
		ctx,
		`INSERT INTO users (username, password_hash) 
		VALUES ($1, $2)`,
		testUsername,
		string(hashedPassword),
	)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Cleanup
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = repo.DB.Exec(
			ctx,
			`DELETE FROM users WHERE username = $1`,
			testUsername,
		)
	}()

	// 1. Successful login
	t.Run("SuccessfulLogin", func(t *testing.T) {
		user, err := loginService.Login(testUsername, testPassword)

		if err != nil {
			t.Fatalf("Expected successful login, got error: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user object, got nil")
		}

		if user.Username != testUsername {
			t.Fatalf("Expected username %s, got %s", testUsername, user.Username)
		}
	})

	// 2. Wrong password
	t.Run("WrongPassword", func(t *testing.T) {
		user, err := loginService.Login(testUsername, "wrongpassword")

		if err == nil {
			t.Fatal("Expected error for wrong password, got nil")
		}

		if user != nil {
			t.Fatalf("Expected nil user, got %+v", user)
		}

		if err.Error() != "invalid password" {
			t.Fatalf("Expected 'invalid password' error, got: %v", err)
		}
	})

	// 3. Non-existent user
	t.Run("NonExistentUser", func(t *testing.T) {
		user, err := loginService.Login("nonexistentuser", testPassword)

		if err == nil {
			t.Fatal("Expected error for non-existent user, got nil")
		}

		if user != nil {
			t.Fatalf("Expected nil user, got %+v", user)
		}

		if err.Error() != "user not found" {
			t.Fatalf("Expected 'user not found' error, got: %v", err)
		}
	})

	// 4. Empty Username
	t.Run("EmptyUsername", func(t *testing.T) {
		user, err := loginService.Login("", testPassword)

		if err == nil {
			t.Fatal("Expected error for empty username, got nil")
		}

		if user != nil {
			t.Fatalf("Expected nil user, got %+v", user)
		}

		if err.Error() != "username and password cannot be empty" {
			t.Fatalf("Expected 'username and password cannot be empty' error, got: %v", err)
		}
	})

	// 5. Empty Password
	t.Run("EmptyPassword", func(t *testing.T) {
		user, err := loginService.Login(testUsername, "")

		if err == nil {
			t.Fatal("Expected error for empty password, got nil")
		}

		if user != nil {
			t.Fatalf("Expected nil user, got %+v", user)
		}

		if err.Error() != "username and password cannot be empty" {
			t.Fatalf("Expected 'username and password cannot be empty' error, got: %v", err)
		}
	})
}
