// Run `go test -v ./internal/api` in /backend to run all tests
// Run `go test -v ./internal/api -run {test name e.g. 'TestGetAllTopics'}` to run a specific test
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

	postService := service.NewPostService(repo)
	postHandler := NewPostHandler(postService)

	commentService := service.NewCommentService(repo)
	commentHandler := NewCommentHandler(commentService)

	// Set up router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/topics", topicHandler.GetAllTopics)
		v1.POST("/users", userHandler.RegisterUser)
		v1.GET("/topics/:topicId/posts", postHandler.GetPostsByTopicID)
		v1.GET("/posts/:postID/comments", commentHandler.GetCommentsByPostID)
	}

	return router, repo
}

func clearTestData(test *testing.T, repo *data.Repository, usernames []string, topicIDs []int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete Posts
	for _, topicID := range topicIDs {
		_, err := repo.DB.Exec(
			ctx,
			`DELETE FROM posts
			WHERE topic_id = $1`,
			topicID,
		)

		if err != nil {
			test.Logf("Warning: Failed to delete posts for test topic ID %d during teardown: %v", topicID, err)
		}
	}

	// Delete Topics
	for _, topicID := range topicIDs {
		_, err := repo.DB.Exec(
			ctx,
			`DELETE FROM topics 
			WHERE topic_id = $1`,
			topicID,
		)

		if err != nil {
			test.Logf("Warning: Failed to delete test topic ID %d during teardown: %v", topicID, err)
		}
	}

	// Delete Users
	for _, username := range usernames {
		_, err := repo.DB.Exec(
			ctx,
			`DELETE FROM users 
			WHERE username = $1`,
			username,
		)

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

	// Create User
	var userID int
	err := repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash) 
		VALUES ($1, $2) 
		RETURNING user_id`,
		testUsername,
		"fakehash",
	).Scan(&userID)

	if err != nil {
		test.Fatalf("Failed to setup user for topic test: %v", err)
	}

	// Create Topics
	var topicID1, topicID2 int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by) 
		VALUES ($1, $2, $3) 
		RETURNING topic_id`,
		"Test Topic 1",
		"Description 1",
		userID,
	).Scan(&topicID1)

	if err != nil {
		test.Fatal(err)
	}
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by) 
		VALUES ($1, $2, $3) 
		RETURNING topic_id`,
		"Test Topic 2",
		"Description 2",
		userID,
	).Scan(&topicID2)

	if err != nil {
		test.Fatal(err)
	}

	topicIDs = append(topicIDs, topicID1, topicID2)

	defer clearTestData(test, repo, []string{testUsername}, topicIDs)

	// Execute request (GET /api/v1/topics)
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

func TestGetPostsByTopicID(test *testing.T) {
	router, repo := setupRouter(test)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create user
	var userID int
	err := repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash) 
		VALUES ($1, $2) 
		RETURNING user_id`,
		"post_test_user",
		"fakehash",
	).Scan(&userID)

	if err != nil {
		test.Fatalf("Failed to setup user for post test: %v", err)
	}

	// Create topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by) 
		VALUES ($1, $2, $3) 
		RETURNING topic_id`,
		"Post Test Topic",
		"Topic for Post Test",
		userID,
	).Scan(&topicID)

	if err != nil {
		test.Fatalf("Failed to create topic for post test: %v", err)
	}

	// Create posts
	_, err = repo.DB.Exec(
		ctx,
		`INSERT INTO posts (topic_id, title, content, created_by) 
		VALUES ($1, $2, $3, $4), ($1, $5, $6, $7)`,
		topicID,
		"Post Title 1",
		"Post Content 1",
		userID,
		"Post Title 2",
		"Post Content 2",
		userID,
	)

	if err != nil {
		test.Fatalf("Failed to create posts for post test: %v", err)
	}

	defer clearTestData(test, repo, []string{"post_test_user"}, []int{topicID})

	// Execute request (GET /api/v1/topics/:topicId/posts)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		test.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response []data.Post
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		test.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}

	if len(response) != 2 {
		test.Errorf("Expected 2 posts, got %d", len(response))
	}

	titles := make(map[string]bool)
	for _, post := range response {
		if post.TopicID != topicID {
			test.Fatalf("Expected topic_id %d on posts, got %+v", topicID, response)
		}
		titles[post.Title] = true
	}

	if !titles["Post Title 1"] || !titles["Post Title 2"] {
		test.Errorf("Expected posts 'Post Title 1' and 'Post Title 2', got titles: %v", titles)
	}
}

func TestGetCommentsByPostID(test *testing.T) {
	router, repo := setupRouter(test)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create user
	var userID int
	err := repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		"comment_test_user",
		"fakehash",
	).Scan(&userID)

	if err != nil {
		test.Fatalf("Failed to setup user for comment test: %v", err)
	}

	// Create topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Comment Test Topic",
		"Topic for Comment Test",
		userID,
	).Scan(&topicID)

	if err != nil {
		test.Fatalf("Failed to create topic for comment test: %v", err)
	}

	// Create post
	var postID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO posts (topic_id, title, content, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING post_id`,
		topicID,
		"Comment Test Post",
		"Post for Comment Test",
		userID,
	).Scan(&postID)

	if err != nil {
		test.Fatalf("Failed to create post for comment test: %v", err)
	}

	// Create comments
	_, err = repo.DB.Exec(
		ctx,
		`INSERT INTO comments (post_id, content, created_by)
		VALUES ($1, $2, $3), ($1, $4, $5)`,
		postID,
		"Comment Content 1",
		userID,
		"Comment Content 2",
		userID,
	)

	if err != nil {
		test.Fatalf("Failed to create comments for comment test: %v", err)
	}

	defer clearTestData(test, repo, []string{"comment_test_user"}, []int{topicID})

	// Execute request (GET /api/v1/posts/:postID/comments)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d/comments", postID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		test.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response []data.Comment
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		test.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}

	if len(response) != 2 {
		test.Errorf("Expected 2 comments, got %d", len(response))
	}

	for _, comment := range response {
		if comment.PostID != postID {
			test.Fatalf("Expected post_id %d on comments, got %+v", postID, response)
		}
	}
}
