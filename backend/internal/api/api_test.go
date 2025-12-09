// Run `go test -v ./...` in /backend
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// setupRouter initializes the Gin router and all dependencies for testing
func setupRouter(test *testing.T) (*gin.Engine, *data.Repository) {
	dbPool, err := data.OpenDB()
	if err != nil {
		test.Fatalf("Failed to open test database connection: %v", err)
	}
	test.Cleanup(func() { dbPool.Close() }) // Close pool after tests

	repo := data.NewRepository(dbPool)

	// Initialize handlers and services
	topicService := service.NewTopicService(repo)
	topicHandler := NewTopicHandler(topicService)

	userService := service.NewUserService(repo)
	userHandler := NewUserHandler(userService)

	// Set up router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/topics", topicHandler.GetAllTopics)
		v1.POST("/users", userHandler.RegisterUser)
	}

	return router, repo
}

func clearTestData(test *testing.T, repo *data.Repository, usernames []string, topicIDs []int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete Topics
	for _, id := range topicIDs {
		_, err := repo.DB.Exec(ctx, "DELETE FROM topics WHERE topic_id = $1", id)
		if err != nil {
			test.Logf("Warning: Failed to delete test topic ID %d during teardown: %v", id, err)
		}
	}

	// Delete Users
	for _, username := range usernames {
		_, err := repo.DB.Exec(ctx, "DELETE FROM users WHERE username = $1", username)
		if err != nil {
			test.Logf("Warning: Failed to delete test user %s during teardown: %v", username, err)
		}
	}
}

func TestUserRegistration(test *testing.T) {
	router, repo := setupRouter(test)
	testUsername := "test_register_user_1"

	defer clearTestData(test, repo, []string{testUsername}, nil)

	// 1. Test Success Case (HTTP 201 Created)
	test.Run("Success_UserCreated", func(t *testing.T) {
		payload := map[string]string{
			"username": testUsername,
			"password": "SecurePassword123",
		}
		jsonPayload, _ := json.Marshal(payload)

		// httptest simulates network request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// httptest records response
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req) // Fire handler

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if user, ok := response["user"].(map[string]interface{}); ok {
			if user["username"] != testUsername {
				t.Errorf("Expected username %s, got %s", testUsername, user["username"])
			}
		} else {
			t.Fatal("Response body missing 'user' object.")
		}
	})

	// 2. Test Failure Case (Weak Password)
	test.Run("Failure_WeakPassword", func(t *testing.T) {
		payload := map[string]string{
			"username": "weakpass",
			"password": "a",
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for weak password, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// 3. Test Failure Case (Duplicate Username)
	test.Run("Failure_DuplicateUsername", func(t *testing.T) {
		payload := map[string]string{
			"username": testUsername, // This user was already created in the first test
			"password": "AnotherPassword456",
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for duplicate username, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["error"] != fmt.Sprintf("username '%s' is already taken", testUsername) {
			t.Errorf("Unexpected error message: %s", response["error"])
		}
	})
}

func TestGetAllTopics(test *testing.T) {
	router, repo := setupRouter(test)
	testUsername := "topic_test_user_2"
	topicIDs := make([]int, 0, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert User
	var userID int
	err := repo.DB.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING user_id", testUsername, "fakehash").Scan(&userID)
	if err != nil {
		test.Fatalf("Failed to setup user for topic test: %v", err)
	}

	// Insert Topics
	var topicID1, topicID2 int
	err = repo.DB.QueryRow(ctx, "INSERT INTO topics (title, description, created_by) VALUES ($1, $2, $3) RETURNING topic_id", "Test Topic 1", "Description 1", userID).Scan(&topicID1)
	if err != nil {
		test.Fatal(err)
	}
	err = repo.DB.QueryRow(ctx, "INSERT INTO topics (title, description, created_by) VALUES ($1, $2, $3) RETURNING topic_id", "Test Topic 2", "Description 2", userID).Scan(&topicID2)
	if err != nil {
		test.Fatal(err)
	}

	topicIDs = append(topicIDs, topicID1, topicID2)

	// Move defer HERE - right after setup, before the test executes
	defer clearTestData(test, repo, []string{testUsername}, topicIDs)

	// Execute Test (GET /api/v1/topics)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		test.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response []data.Topic
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		test.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}

	if len(response) != 2 {
		test.Errorf("Expected 2 topics, got %d", len(response))
	}

	titles := make(map[string]bool)
	for _, topic := range response {
		titles[topic.Title] = true
	}

	if !titles["Test Topic 1"] || !titles["Test Topic 2"] {
		test.Errorf("Expected topics 'Test Topic 1' and 'Test Topic 2', got titles: %v", titles)
	}
}
