/*
Run `go test -v ./internal/api` in /backend to run all tests
Run `go test -v ./internal/api -run {test name e.g. 'TestGetAllTopics'}` to run a specific test
Run `docker compose exec db psql -U user -d forum_db` to access database, then run
SELECT 'users' as table_name, COUNT(*) FROM users UNION ALL SELECT 'topics', COUNT(*) FROM topics UNION ALL SELECT 'posts', COUNT(*) FROM posts UNION ALL SELECT 'comments', COUNT(*) FROM comments;
to check if db has been cleared
*/
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

	"golang.org/x/crypto/bcrypt"

	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// setupRouter initializes the Gin router and all dependencies for testing
func setupRouter(t *testing.T) (*gin.Engine, *data.Repository) {
	dbPool, err := data.OpenDB()
	if err != nil {
		t.Fatalf("Failed to open test database connection: %v", err)
	}
	t.Cleanup(func() { dbPool.Close() }) // Close pool after tests

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

	jwtService := service.NewJWTService("test-secret-key", 1*time.Hour)

	loginService := service.NewLoginService(repo)
	loginHandler := NewLoginHandler(loginService, jwtService)

	// Set up router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{

		v1.GET("/topics", topicHandler.GetAllTopics)
		v1.POST("/users", userHandler.RegisterUser)
		v1.GET("/topics/:topicID/posts", postHandler.GetPostsByTopicID)
		v1.GET("/posts/:postID/comments", commentHandler.GetCommentsByPostID)
		v1.POST("/login", loginHandler.LoginUser)

		// Protected Routes
		protected := v1.Group("")
		protected.Use(AuthMiddleware(jwtService))
		{
			protected.POST("/topics", topicHandler.CreateTopic)
			protected.PUT("/topics/:topicID", topicHandler.UpdateTopic)
			protected.DELETE("/topics/:topicID", topicHandler.DeleteTopic)

			protected.POST("/topics/:topicID/posts", postHandler.CreatePost)
			protected.PUT("/posts/:postID", postHandler.UpdatePost)
			protected.DELETE("/posts/:postID", postHandler.DeletePost)

			protected.POST("/posts/:postID/comments", commentHandler.CreateComment)
			protected.PUT("/comments/:commentID", commentHandler.UpdateComment)
			protected.DELETE("/comments/:commentID", commentHandler.DeleteComment)
		}
	}

	return router, repo
}

func clearTestData(t *testing.T, repo *data.Repository, usernames []string, topicIDs []int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete Topics (CASCADE deletes Posts and Comments automatically)
	for _, topicID := range topicIDs {
		_, err := repo.DB.Exec(
			ctx,
			`DELETE FROM topics 
			WHERE topic_id = $1`,
			topicID,
		)

		if err != nil {
			t.Logf("Warning: Failed to delete test topic ID %d during teardown: %v", topicID, err)
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
			t.Logf("Warning: Failed to delete test user %s during teardown: %v", username, err)
		}
	}
}

func TestUserRegistration(t *testing.T) {
	router, repo := setupRouter(t)
	testUsername := "test_register_user"

	defer clearTestData(t, repo, []string{testUsername}, nil)

	// 1. Test Success Case (HTTP 201 Created)
	t.Run("Success_UserCreated", func(t *testing.T) {
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
			if username, ok := user["username"].(string); ok && username == testUsername {
				// Username matches expected value
			} else {
				t.Errorf("Expected username %s, got %v", testUsername, user["username"])
			}
		} else {
			t.Fatal("Response body missing 'user' object.")
		}
	})

	// 2. Test Failure Case (Weak Password)
	t.Run("Failure_WeakPassword", func(t *testing.T) {
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
	t.Run("Failure_DuplicateUsername", func(t *testing.T) {
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

func TestGetAllTopics(t *testing.T) {
	router, repo := setupRouter(t)
	testUsername := "test_get_topics_user"
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
		t.Fatalf("Failed to setup user for topic test: %v", err)
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
		t.Fatal(err)
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
		t.Fatal(err)
	}

	topicIDs = append(topicIDs, topicID1, topicID2)

	defer clearTestData(t, repo, []string{testUsername}, topicIDs)

	// Execute request (GET /api/v1/topics)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response []data.Topic
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 topics, got %d", len(response))
	}

	titles := make(map[string]bool)
	for _, topic := range response {
		titles[topic.Title] = true
	}

	if !titles["Test Topic 1"] || !titles["Test Topic 2"] {
		t.Errorf("Expected topics 'Test Topic 1' and 'Test Topic 2', got titles: %v", titles)
	}
}

func TestGetPostsByTopicID(t *testing.T) {
	router, repo := setupRouter(t)

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
		t.Fatalf("Failed to setup user for post test: %v", err)
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
		t.Fatalf("Failed to create topic for post test: %v", err)
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
		t.Fatalf("Failed to create posts for post test: %v", err)
	}

	defer clearTestData(t, repo, []string{"post_test_user"}, []int{topicID})

	// Execute request (GET /api/v1/topics/:topicID/posts)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response []data.Post
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(response))
	}

	titles := make(map[string]bool)
	for _, post := range response {
		if post.TopicID != topicID {
			t.Fatalf("Expected topic_id %d on posts, got %+v", topicID, response)
		}
		titles[post.Title] = true
	}

	if !titles["Post Title 1"] || !titles["Post Title 2"] {
		t.Errorf("Expected posts 'Post Title 1' and 'Post Title 2', got titles: %v", titles)
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	router, repo := setupRouter(t)

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
		t.Fatalf("Failed to setup user for comment test: %v", err)
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
		t.Fatalf("Failed to create topic for comment test: %v", err)
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
		t.Fatalf("Failed to create post for comment test: %v", err)
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
		t.Fatalf("Failed to create comments for comment test: %v", err)
	}

	defer clearTestData(t, repo, []string{"comment_test_user"}, []int{topicID})

	// Execute request (GET /api/v1/posts/:postID/comments)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%d/comments", postID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response []data.Comment
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(response))
	}

	for _, comment := range response {
		if comment.PostID != postID {
			t.Fatalf("Expected post_id %d on comments, got %+v", postID, response)
		}
	}
}

func TestLogin(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_login_user"
	testPassword := "test_login_password"
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
	defer clearTestData(t, repo, []string{testUsername}, nil)

	// 1. Successful login
	t.Run("SuccessfulLogin", func(t *testing.T) {
		payload := map[string]string{
			"username": testUsername,
			"password": testPassword,
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, w.Body.String())
		}

		tokenString, exists := response["token"]
		if !exists || tokenString == "" {
			t.Fatal("Response missing 'token' field")
		}

		// Validate JWT token
		jwtService := service.NewJWTService("test-secret-key", 1*time.Hour)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			t.Fatalf("Failed to validate JWT token: %v", err)
		}

		if claims.Username != testUsername {
			t.Errorf("Expected token username %s, got %s", testUsername, claims.Username)
		}
	})

	// 2. Wrong Password
	t.Run("FailedLogin_WrongPassword", func(t *testing.T) {
		payload := map[string]string{
			"username": testUsername,
			"password": "wrong_password",
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for wrong password, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Non-existent User
	t.Run("FailedLogin_NonExistentUser", func(t *testing.T) {
		payload := map[string]string{
			"username": "non_existent_user",
			"password": "some_password",
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for non-existent user, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 4. Missing Fields
	t.Run("FailedLogin_MissingFields", func(t *testing.T) {
		payload := map[string]string{
			"username": testUsername,
			// Missing password
		}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for missing fields, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestAuthMiddleware(t *testing.T) {
	router, repo := setupRouter(t)

	// Create test user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testUsername := "test_auth_user"
	testPassword := "test_auth_password"

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
	defer clearTestData(t, repo, []string{testUsername}, nil)

	// Login to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}
	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed with status %d. Response: %s", w.Code, w.Body.String())
	}

	var loginResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("Failed to unmarshal login response: %v. Body: %s", err, w.Body.String())
	}

	validTokenString, exists := loginResponse["token"]
	if !exists || validTokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Add protected test route
	jwtService := service.NewJWTService("test-secret-key", 1*time.Hour)
	protected := router.Group("/api/v1/protected")
	protected.Use(AuthMiddleware(jwtService))
	{
		protected.GET("/test", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{"message": "Access granted"},
			)
		})
	}

	// 1. Access with valid token
	t.Run("AccessWithValidToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
		req.Header.Set("Authorization", "Bearer "+validTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d with valid token, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	// 2. Access without token
	t.Run("AccessWithoutToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d without token, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Access with invalid token
	t.Run("AccessWithInvalidToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d with invalid token, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 4. Access with malformed token
	t.Run("AccessWithMalformedToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
		req.Header.Set("Authorization", validTokenString) // Missing "Bearer " prefix

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d with malformed token, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})
}

func TestCreateTopic(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_create_topic_user"
	testPassword := "test_create_topic_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Login to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Store topicIDs for cleanup
	topicIDs := []int{}

	defer func() {
		clearTestData(t, repo, []string{testUsername}, topicIDs)
	}()

	// 1. Successful Topic Creation
	t.Run("SuccessfulTopicCreation", func(t *testing.T) {
		topicPayload := map[string]string{
			"title":       "New Test Topic",
			"description": "This is a test topic created during testing.",
		}
		jsonTopicPayload, _ := json.Marshal(topicPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewBuffer(jsonTopicPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status %d for topic creation, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response data.Topic
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Title != topicPayload["title"] {
			t.Errorf("Expected topic title %s, got %s", topicPayload["title"], response.Title)
		}

		if response.Description != topicPayload["description"] {
			t.Errorf("Expected topic description %s, got %s", topicPayload["description"], response.Description)
		}

		if response.CreatedBy != userID {
			t.Errorf("Expected topic created_by %d, got %d", userID, response.CreatedBy)
		}

		// Add to cleanup list
		topicIDs = append(topicIDs, response.TopicID)
	})

	// 2. Topic Creation without Authentication Token
	t.Run("TopicCreationWithoutAuthenticationToken", func(t *testing.T) {
		topicPayload := map[string]string{
			"title":       "Unauthorized Topic",
			"description": "This topic should not be created.",
		}
		jsonTopicPayload, _ := json.Marshal(topicPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewBuffer(jsonTopicPayload))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized topic creation, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Topic Creation with Missing Title
	t.Run("TopicCreationWithMissingTitle", func(t *testing.T) {
		topicPayload := map[string]string{
			// Missing title
			"description": "This topic has no title.",
		}
		jsonTopicPayload, _ := json.Marshal(topicPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewBuffer(jsonTopicPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for topic creation with missing fields, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 4. Topic Creation with Title Exceeding Max Length
	t.Run("TopicCreationWithTitleExceedingMaxLength", func(t *testing.T) {
		longTitle := ""
		for i := 0; i < 201; i++ { // Max length is 200 characters
			longTitle += "a"
		}

		topicPayload := map[string]string{
			"title":       longTitle,
			"description": "This topic has an excessively long title.",
		}
		jsonTopicPayload, _ := json.Marshal(topicPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/topics", bytes.NewBuffer(jsonTopicPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for topic creation with long title, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestCreatePost(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_create_post_user"
	testPassword := "test_create_post_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Test Topic for Post Creation",
		"Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Login to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername}, []int{topicID})
	}()

	// 1. Successful Post Creation
	t.Run("SuccessfulPostCreation", func(t *testing.T) {
		postPayload := map[string]string{
			"title":   "New Test Post",
			"content": "This is a test post created during testing.",
		}
		jsonPostPayload, _ := json.Marshal(postPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), bytes.NewBuffer(jsonPostPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status %d for post creation, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response data.Post
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Title != postPayload["title"] {
			t.Errorf("Expected post title %s, got %s", postPayload["title"], response.Title)
		}

		if response.Content != postPayload["content"] {
			t.Errorf("Expected post content %s, got %s", postPayload["content"], response.Content)
		}

		if response.TopicID != topicID {
			t.Errorf("Expected post topic_id %d, got %d", topicID, response.TopicID)
		}

		if response.CreatedBy != userID {
			t.Errorf("Expected post created_by %d, got %d", userID, response.CreatedBy)
		}
	})

	// 2. Post Creation without Authentication Token
	t.Run("PostCreationWithoutAuthenticationToken", func(t *testing.T) {
		postPayload := map[string]string{
			"title":   "Unauthorized Post",
			"content": "This post should not be created.",
		}
		jsonPostPayload, _ := json.Marshal(postPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), bytes.NewBuffer(jsonPostPayload))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized post creation, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Post Creation with Missing Title
	t.Run("PostCreationWithMissingTitle", func(t *testing.T) {
		postPayload := map[string]string{
			// Missing title
			"content": "This post has no title.",
		}
		jsonPostPayload, _ := json.Marshal(postPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), bytes.NewBuffer(jsonPostPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post creation with missing fields, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 4. Post Creation with Title Exceeding Max Length
	t.Run("PostCreationWithTitleExceedingMaxLength", func(t *testing.T) {
		longTitle := ""
		for i := 0; i < 201; i++ { // Max length is 200 characters
			longTitle += "a"
		}

		postPayload := map[string]string{
			"title":   longTitle,
			"content": "This post has an excessively long title.",
		}
		jsonPostPayload, _ := json.Marshal(postPayload)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), bytes.NewBuffer(jsonPostPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post creation with long title, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 5. Post Creation with Missing Content
	t.Run("PostCreationWithMissingContent", func(t *testing.T) {
		postPayload := map[string]string{
			"title": "Post Without Content",
			// Missing content
		}
		jsonPostPayload, _ := json.Marshal(postPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), bytes.NewBuffer(jsonPostPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post creation with missing content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 6. Post Creation with Content Exceeding Max Length
	t.Run("PostCreationWithContentExceedingMaxLength", func(t *testing.T) {
		longContent := ""
		for i := 0; i < 5001; i++ { // Max length is 5000 characters
			longContent += "a"
		}

		postPayload := map[string]string{
			"title":   "Post With Long Content",
			"content": longContent,
		}
		jsonPostPayload, _ := json.Marshal(postPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/topics/%d/posts", topicID), bytes.NewBuffer(jsonPostPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post creation with long content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 7. Post Creation under Non-existent Topic
	t.Run("PostCreationUnderNonExistentTopic", func(t *testing.T) {
		postPayload := map[string]string{
			"title":   "Post Under Non-existent Topic",
			"content": "This post is under a topic that does not exist.",
		}
		jsonPostPayload, _ := json.Marshal(postPayload)

		req := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/v1/topics/%d/posts", 9999999), // Assuming this topic ID does not exist
			bytes.NewBuffer(jsonPostPayload),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post creation under non-existent topic, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestCreateComment(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_create_comment_user"
	testPassword := "test_create_comment_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Test Topic for Comment Creation",
		"Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO posts (topic_id, title, content, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING post_id`,
		topicID,
		"Test Post for Comment Creation",
		"Post Content",
		userID,
	).Scan(&postID)

	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Login to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername}, []int{topicID})
	}()

	// 1. Successful Comment Creation
	t.Run("SuccessfulCommentCreation", func(t *testing.T) {
		commentPayload := map[string]string{
			"content": "This is a test comment created during testing.",
		}
		jsonCommentPayload, _ := json.Marshal(commentPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", postID), bytes.NewBuffer(jsonCommentPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status %d for comment creation, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response data.Comment
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Content != commentPayload["content"] {
			t.Errorf("Expected comment content %s, got %s", commentPayload["content"], response.Content)
		}

		if response.PostID != postID {
			t.Errorf("Expected comment post_id %d, got %d", postID, response.PostID)
		}

		if response.CreatedBy != userID {
			t.Errorf("Expected comment created_by %d, got %d", userID, response.CreatedBy)
		}
	})

	// 2. Comment Creation without Authentication Token
	t.Run("CommentCreationWithoutAuthenticationToken", func(t *testing.T) {
		commentPayload := map[string]string{
			"content": "This comment should not be created.",
		}
		jsonCommentPayload, _ := json.Marshal(commentPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", postID), bytes.NewBuffer(jsonCommentPayload))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized comment creation, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Comment Creation with Missing Content
	t.Run("CommentCreationWithMissingContent", func(t *testing.T) {
		commentPayload := map[string]string{
			// Missing content
		}
		jsonCommentPayload, _ := json.Marshal(commentPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", postID), bytes.NewBuffer(jsonCommentPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for comment creation with missing content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 4. Comment Creation with Content Exceeding Max Length
	t.Run("CommentCreationWithContentExceedingMaxLength", func(t *testing.T) {
		longContent := ""
		for i := 0; i < 2001; i++ { // Max length is 2000 characters
			longContent += "a"
		}

		commentPayload := map[string]string{
			"content": longContent,
		}
		jsonCommentPayload, _ := json.Marshal(commentPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", postID), bytes.NewBuffer(jsonCommentPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for comment creation with long content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 5. Comment Creation under Non-existent Post
	t.Run("CommentCreationUnderNonExistentPost", func(t *testing.T) {
		commentPayload := map[string]string{
			"content": "This comment is under a post that does not exist.",
		}
		jsonCommentPayload, _ := json.Marshal(commentPayload)

		req := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/v1/posts/%d/comments", 9999999), // Assuming this post ID does not exist
			bytes.NewBuffer(jsonCommentPayload),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Expected status %d for comment creation under non-existent post, got %d. Response: %s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})

	// 6. Comment Creation with Empty Content
	t.Run("CommentCreationWithEmptyContent", func(t *testing.T) {
		commentPayload := map[string]string{
			"content": "",
		}
		jsonCommentPayload, _ := json.Marshal(commentPayload)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/posts/%d/comments", postID), bytes.NewBuffer(jsonCommentPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for comment creation with empty content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestUpdateTopic(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_update_topic_user"
	testPassword := "test_update_topic_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create another user (for unauthorised test)
	otherUsername := "other_user"
	otherPassword := "other_user_password"

	otherHashedPassword, err := bcrypt.GenerateFromPassword([]byte(otherPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash other user password: %v", err)
	}

	var otherUserID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		otherUsername,
		string(otherHashedPassword),
	).Scan(&otherUserID)

	if err != nil {
		t.Fatalf("Failed to create other test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Original Topic Title",
		"Original Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Login as topic creator to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Login as other user to get JWT token
	otherLoginPayload := map[string]string{
		"username": otherUsername,
		"password": otherPassword,
	}

	jsonOtherLoginPayload, _ := json.Marshal(otherLoginPayload)

	otherReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonOtherLoginPayload))
	otherReq.Header.Set("Content-Type", "application/json")

	otherW := httptest.NewRecorder()

	router.ServeHTTP(otherW, otherReq)

	var otherLoginResponse map[string]string

	json.Unmarshal(otherW.Body.Bytes(), &otherLoginResponse)
	otherTokenString, otherExists := otherLoginResponse["token"]
	if !otherExists || otherTokenString == "" {
		t.Fatal("Other login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername, otherUsername}, []int{topicID})
	}()

	// 1. Successful Topic Update
	t.Run("SuccessfulTopicUpdate", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":       "Updated Topic Title",
			"description": "Updated Topic Description",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/topics/%d", topicID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d for topic update, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response data.Topic
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Title != updatePayload["title"] {
			t.Errorf("Expected updated topic title %s, got %s", updatePayload["title"], response.Title)
		}

		if response.Description != updatePayload["description"] {
			t.Errorf("Expected updated topic description %s, got %s", updatePayload["description"], response.Description)
		}
	})

	// 2. Topic Update without Authentication Token
	t.Run("TopicUpdateWithoutAuthenticationToken", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":       "Unauthorized Update Title",
			"description": "Unauthorized Update Description",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/topics/%d", topicID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized topic update, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Topic Update by Non-Owner
	t.Run("TopicUpdateByNonOwner", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":       "Non-Owner Update Title",
			"description": "Non-Owner Update Description",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/topics/%d", topicID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected status %d for topic update by non-owner, got %d. Response: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// 4. Update Non-Existent Topic
	t.Run("UpdateNonExistentTopic", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":       "Update Non-Existent Title",
			"description": "Update Non-Existent Description",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/topics/%d", 9999999), bytes.NewBuffer(jsonUpdatePayload)) // Assuming this topic ID does not exist
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status %d for update of non-existent topic, got %d. Response: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	// 5. Topic Update with Missing Title
	t.Run("TopicUpdateWithMissingTitle", func(t *testing.T) {
		updatePayload := map[string]string{
			// Missing title
			"description": "Update with Missing Title Description",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/topics/%d", topicID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for topic update with missing title, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 6. Topic Update with Title Exceeding Max Length
	t.Run("TopicUpdateWithTitleExceedingMaxLength", func(t *testing.T) {
		longTitle := ""
		for i := 0; i < 201; i++ { // Max length is 200 characters
			longTitle += "a"
		}

		updatePayload := map[string]string{
			"title":       longTitle,
			"description": "Update with Long Title Description",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/topics/%d", topicID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for topic update with long title, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestUpdatePost(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_update_post_user"
	testPassword := "test_update_post_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create another user (for unauthorised test)
	otherUsername := "other_post_user"
	otherPassword := "other_post_password"

	otherHashedPassword, err := bcrypt.GenerateFromPassword([]byte(otherPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash other user password: %v", err)
	}

	var otherUserID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		otherUsername,
		string(otherHashedPassword),
	).Scan(&otherUserID)

	if err != nil {
		t.Fatalf("Failed to create other test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Topic for Post Update",
		"Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO posts (topic_id, title, content, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING post_id`,
		topicID,
		"Original Post Title",
		"Original Post Content",
		userID,
	).Scan(&postID)

	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Login as post creator to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Login as other user to get JWT token
	otherLoginPayload := map[string]string{
		"username": otherUsername,
		"password": otherPassword,
	}

	jsonOtherLoginPayload, _ := json.Marshal(otherLoginPayload)

	otherReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonOtherLoginPayload))
	otherReq.Header.Set("Content-Type", "application/json")

	otherW := httptest.NewRecorder()

	router.ServeHTTP(otherW, otherReq)

	var otherLoginResponse map[string]string

	json.Unmarshal(otherW.Body.Bytes(), &otherLoginResponse)
	otherTokenString, otherExists := otherLoginResponse["token"]
	if !otherExists || otherTokenString == "" {
		t.Fatal("Other login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername, otherUsername}, []int{topicID})
	}()

	// 1. Successful Post Update
	t.Run("SuccessfulPostUpdate", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":   "Updated Post Title",
			"content": "Updated Post Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d for post update, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response data.Post
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Title != updatePayload["title"] {
			t.Errorf("Expected updated post title %s, got %s", updatePayload["title"], response.Title)
		}

		if response.Content != updatePayload["content"] {
			t.Errorf("Expected updated post content %s, got %s", updatePayload["content"], response.Content)
		}
	})

	// 2. Post Update without Authentication Token
	t.Run("PostUpdateWithoutAuthenticationToken", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":   "Unauthorized Update Title",
			"content": "Unauthorized Update Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized post update, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Post Update by Non-Owner
	t.Run("PostUpdateByNonOwner", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":   "Non-Owner Update Title",
			"content": "Non-Owner Update Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected status %d for post update by non-owner, got %d. Response: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// 4. Update Non-Existent Post
	t.Run("UpdateNonExistentPost", func(t *testing.T) {
		updatePayload := map[string]string{
			"title":   "Update Non-Existent Title",
			"content": "Update Non-Existent Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", 9999999), bytes.NewBuffer(jsonUpdatePayload)) // Assuming this post ID does not exist
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status %d for update of non-existent post, got %d. Response: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	// 5. Post Update with Missing Title
	t.Run("PostUpdateWithMissingTitle", func(t *testing.T) {
		updatePayload := map[string]string{
			// Missing title
			"content": "Update with Missing Title Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post update with missing title, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 6. Post Update with Title Exceeding Max Length
	t.Run("PostUpdateWithTitleExceedingMaxLength", func(t *testing.T) {
		longTitle := ""
		for i := 0; i < 201; i++ { // Max length is 200 characters
			longTitle += "a"
		}

		updatePayload := map[string]string{
			"title":   longTitle,
			"content": "Update with Long Title Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post update with long title, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 7. Post Update with Content Exceeding Max Length
	t.Run("PostUpdateWithContentExceedingMaxLength", func(t *testing.T) {
		longContent := ""
		for i := 0; i < 5001; i++ { // Max length is 5000 characters
			longContent += "a"
		}

		updatePayload := map[string]string{
			"title":   "Update With Long Content Title",
			"content": longContent,
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/posts/%d", postID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post update with long content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestUpdateComment(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_update_comment_user"
	testPassword := "test_update_comment_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create another user (for unauthorised test)
	otherUsername := "other_comment_user"
	otherPassword := "other_comment_password"

	otherHashedPassword, err := bcrypt.GenerateFromPassword([]byte(otherPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash other user password: %v", err)
	}

	var otherUserID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		otherUsername,
		string(otherHashedPassword),
	).Scan(&otherUserID)

	if err != nil {
		t.Fatalf("Failed to create other test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Topic for Comment Update",
		"Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO posts (topic_id, title, content, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING post_id`,
		topicID,
		"Post for Comment Update",
		"Post Content",
		userID,
	).Scan(&postID)

	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Create test comment
	var commentID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO comments (post_id, content, created_by)
		VALUES ($1, $2, $3)
		RETURNING comment_id`,
		postID,
		"Original Comment Content",
		userID,
	).Scan(&commentID)

	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	// Login as comment creator to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Login as other user to get JWT token
	otherLoginPayload := map[string]string{
		"username": otherUsername,
		"password": otherPassword,
	}

	jsonOtherLoginPayload, _ := json.Marshal(otherLoginPayload)

	otherReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonOtherLoginPayload))
	otherReq.Header.Set("Content-Type", "application/json")

	otherW := httptest.NewRecorder()

	router.ServeHTTP(otherW, otherReq)

	var otherLoginResponse map[string]string

	json.Unmarshal(otherW.Body.Bytes(), &otherLoginResponse)
	otherTokenString, otherExists := otherLoginResponse["token"]
	if !otherExists || otherTokenString == "" {
		t.Fatal("Other login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername, otherUsername}, []int{topicID})
	}()

	// 1. Successful Comment Update
	t.Run("SuccessfulCommentUpdate", func(t *testing.T) {
		updatePayload := map[string]string{
			"content": "Updated Comment Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/comments/%d", commentID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d for comment update, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response data.Comment
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Content != updatePayload["content"] {
			t.Errorf("Expected updated comment content %s, got %s", updatePayload["content"], response.Content)
		}
	})

	// 2. Comment Update without Authentication Token
	t.Run("CommentUpdateWithoutAuthenticationToken", func(t *testing.T) {
		updatePayload := map[string]string{
			"content": "Unauthorized Update Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/comments/%d", commentID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized comment update, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Comment Update by Non-Owner
	t.Run("CommentUpdateByNonOwner", func(t *testing.T) {
		updatePayload := map[string]string{
			"content": "Non-Owner Update Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/comments/%d", commentID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected status %d for comment update by non-owner, got %d. Response: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// 4. Update Non-Existent Comment
	t.Run("UpdateNonExistentComment", func(t *testing.T) {
		updatePayload := map[string]string{
			"content": "Update Non-Existent Content",
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/comments/%d", 9999999), bytes.NewBuffer(jsonUpdatePayload)) // Assuming this comment ID does not exist
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status %d for update of non-existent comment, got %d. Response: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	// 5. Comment Update with Missing Content
	t.Run("CommentUpdateWithMissingContent", func(t *testing.T) {
		updatePayload := map[string]string{
			// Missing content
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/comments/%d", commentID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for comment update with missing content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	// 6. Comment Update with Content Exceeding Max Length
	t.Run("CommentUpdateWithContentExceedingMaxLength", func(t *testing.T) {
		longContent := ""
		for i := 0; i < 2001; i++ { // Max length is 2000 characters
			longContent += "a"
		}

		updatePayload := map[string]string{
			"content": longContent,
		}
		jsonUpdatePayload, _ := json.Marshal(updatePayload)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/comments/%d", commentID), bytes.NewBuffer(jsonUpdatePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for comment update with long content, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestDeleteComment(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_delete_comment_user"
	testPassword := "test_delete_comment_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create another user (for unauthorised test)
	otherUsername := "other_delete_comment_user"
	otherPassword := "other_delete_comment_password"

	otherHashedPassword, err := bcrypt.GenerateFromPassword([]byte(otherPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash other user password: %v", err)
	}

	var otherUserID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		otherUsername,
		string(otherHashedPassword),
	).Scan(&otherUserID)

	if err != nil {
		t.Fatalf("Failed to create other test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Topic for Delete Comment",
		"Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO posts (topic_id, title, content, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING post_id`,
		topicID,
		"Post for Delete Comment",
		"Post Content",
		userID,
	).Scan(&postID)

	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Login as comment creator to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Login as other user to get JWT token
	otherLoginPayload := map[string]string{
		"username": otherUsername,
		"password": otherPassword,
	}

	jsonOtherLoginPayload, _ := json.Marshal(otherLoginPayload)

	otherReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonOtherLoginPayload))
	otherReq.Header.Set("Content-Type", "application/json")

	otherW := httptest.NewRecorder()

	router.ServeHTTP(otherW, otherReq)

	var otherLoginResponse map[string]string

	json.Unmarshal(otherW.Body.Bytes(), &otherLoginResponse)
	otherTokenString, otherExists := otherLoginResponse["token"]
	if !otherExists || otherTokenString == "" {
		t.Fatal("Other login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername, otherUsername}, []int{topicID})
	}()

	// 1. Successful Comment Deletion
	t.Run("SuccessfulCommentDeletion", func(t *testing.T) {
		// Create test comment
		var commentID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO comments (post_id, content, created_by)
			VALUES ($1, $2, $3)
			RETURNING comment_id`,
			postID,
			"Comment to be Deleted",
			userID,
		).Scan(&commentID)

		if err != nil {
			t.Fatalf("Failed to create test comment: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", commentID), nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("Expected status %d for comment deletion, got %d. Response: %s", http.StatusNoContent, w.Code, w.Body.String())
		}

		// Verify comment was deleted
		var commentCount int
		err = repo.DB.QueryRow(
			ctx,
			`SELECT COUNT(*) 
			FROM comments
			WHERE comment_id = $1`,
			commentID,
		).Scan(&commentCount)

		if err != nil {
			t.Fatalf("Failed to verify comment deletion: %v", err)
		}

		if commentCount != 0 {
			t.Fatalf("Comment with ID %d was not deleted", commentID)
		}
	})

	// 2. Comment Deletion without Authentication Token
	t.Run("CommentDeletionWithoutAuthenticationToken", func(t *testing.T) {
		// Create test comment
		var commentID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO comments (post_id, content, created_by)
			VALUES ($1, $2, $3)
			RETURNING comment_id`,
			postID,
			"Comment for Unauthorized Deletion",
			userID,
		).Scan(&commentID)

		if err != nil {
			t.Fatalf("Failed to create test comment: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", commentID), nil)
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized comment deletion, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Comment Deletion by Non-Owner
	t.Run("CommentDeletionByNonOwner", func(t *testing.T) {
		// Create test comment
		var commentID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO comments (post_id, content, created_by)
			VALUES ($1, $2, $3)
			RETURNING comment_id`,
			postID,
			"Comment for Non-Owner Deletion",
			userID,
		).Scan(&commentID)

		if err != nil {
			t.Fatalf("Failed to create test comment: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", commentID), nil)
		req.Header.Set("Authorization", "Bearer "+otherTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected status %d for comment deletion by non-owner, got %d. Response: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// 4. Deletion of Non-Existent Comment
	t.Run("DeletionOfNonExistentComment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/comments/%d", 9999999), nil) // Assuming this comment ID does not exist
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status %d for deletion of non-existent comment, got %d. Response: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	// 5. Comment Deletion with Invalid Comment ID
	t.Run("CommentDeletionWithInvalidCommentID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/invalid_id", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for comment deletion with invalid comment ID, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestDeletePost(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_delete_post_user"
	testPassword := "test_delete_post_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create another user (for unauthorised test)
	otherUsername := "other_delete_post_user"
	otherPassword := "other_delete_post_password"

	otherHashedPassword, err := bcrypt.GenerateFromPassword([]byte(otherPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash other user password: %v", err)
	}

	var otherUserID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		otherUsername,
		string(otherHashedPassword),
	).Scan(&otherUserID)

	if err != nil {
		t.Fatalf("Failed to create other test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO topics (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING topic_id`,
		"Topic for Delete Post",
		"Topic Description",
		userID,
	).Scan(&topicID)

	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Login as post creator to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Login as other user to get JWT token
	otherLoginPayload := map[string]string{
		"username": otherUsername,
		"password": otherPassword,
	}

	jsonOtherLoginPayload, _ := json.Marshal(otherLoginPayload)

	otherReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonOtherLoginPayload))
	otherReq.Header.Set("Content-Type", "application/json")

	otherW := httptest.NewRecorder()

	router.ServeHTTP(otherW, otherReq)

	var otherLoginResponse map[string]string

	json.Unmarshal(otherW.Body.Bytes(), &otherLoginResponse)
	otherTokenString, otherExists := otherLoginResponse["token"]
	if !otherExists || otherTokenString == "" {
		t.Fatal("Other login response missing 'token' field")
	}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername, otherUsername}, []int{topicID})
	}()

	// 1. Successful Post Deletion
	t.Run("SuccessfulPostDeletion", func(t *testing.T) {
		// Create test post
		var postID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO posts (topic_id, title, content, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING post_id`,
			topicID,
			"Post to be Deleted",
			"Post Content",
			userID,
		).Scan(&postID)

		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}

		// Create test comments under the post
		for i := 0; i < 3; i++ {
			_, err = repo.DB.Exec(
				ctx,
				`INSERT INTO comments (post_id, content, created_by)
				VALUES ($1, $2, $3)`,
				postID,
				fmt.Sprintf("Comment %d for Deletion Test", i+1),
				userID,
			)

			if err != nil {
				t.Fatalf("Failed to create test comment %d: %v", i+1, err)
			}
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", postID), nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("Expected status %d for post deletion, got %d. Response: %s", http.StatusNoContent, w.Code, w.Body.String())
		}

		// Verify post was deleted
		var postCount int
		err = repo.DB.QueryRow(
			ctx,
			`SELECT COUNT(*) 
			FROM posts
			WHERE post_id = $1`,
			postID,
		).Scan(&postCount)

		if err != nil {
			t.Fatalf("Failed to verify post deletion: %v", err)
		}

		if postCount != 0 {
			t.Fatalf("Post with ID %d was not deleted", postID)
		}

		// Verify comments under the post were deleted
		var commentCount int
		err = repo.DB.QueryRow(
			ctx,
			`SELECT COUNT(*) 
			FROM comments
			WHERE post_id = $1`,
			postID,
		).Scan(&commentCount)

		if err != nil {
			t.Fatalf("Failed to verify comments deletion: %v", err)
		}

		if commentCount != 0 {
			t.Fatalf("Comments under post ID %d were not deleted", postID)
		}
	})

	// 2. Post Deletion without Authentication Token
	t.Run("PostDeletionWithoutAuthenticationToken", func(t *testing.T) {
		// Create test post
		var postID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO posts (topic_id, title, content, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING post_id`,
			topicID,
			"Post for Unauthorized Deletion",
			"Post Content",
			userID,
		).Scan(&postID)

		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", postID), nil)
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized post deletion, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Post Deletion by Non-Owner
	t.Run("PostDeletionByNonOwner", func(t *testing.T) {
		// Create test post
		var postID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO posts (topic_id, title, content, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING post_id`,
			topicID,
			"Post for Non-Owner Deletion",
			"Post Content",
			userID,
		).Scan(&postID)

		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", postID), nil)
		req.Header.Set("Authorization", "Bearer "+otherTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected status %d for post deletion by non-owner, got %d. Response: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// 4. Deletion of Non-Existent Post
	t.Run("DeletionOfNonExistentPost", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/posts/%d", 9999999), nil) // Assuming this post ID does not exist
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status %d for deletion of non-existent post, got %d. Response: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	// 5. Post Deletion with Invalid Post ID
	t.Run("PostDeletionWithInvalidPostID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/invalid_id", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for post deletion with invalid post ID, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestDeleteTopic(t *testing.T) {
	router, repo := setupRouter(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create test user
	testUsername := "test_delete_topic_user"
	testPassword := "test_delete_topic_password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		testUsername,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create another user (for unauthorised test)
	otherUsername := "other_delete_topic_user"
	otherPassword := "other_delete_topic_password"

	otherHashedPassword, err := bcrypt.GenerateFromPassword([]byte(otherPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fatalf("Failed to hash other user password: %v", err)
	}

	var otherUserID int
	err = repo.DB.QueryRow(
		ctx,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING user_id`,
		otherUsername,
		string(otherHashedPassword),
	).Scan(&otherUserID)

	if err != nil {
		t.Fatalf("Failed to create other test user: %v", err)
	}

	// Login as topic creator to get JWT token
	loginPayload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}

	jsonLoginPayload, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonLoginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]string

	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	tokenString, exists := loginResponse["token"]
	if !exists || tokenString == "" {
		t.Fatal("Login response missing 'token' field")
	}

	// Login as other user to get JWT token
	otherLoginPayload := map[string]string{
		"username": otherUsername,
		"password": otherPassword,
	}

	jsonOtherLoginPayload, _ := json.Marshal(otherLoginPayload)

	otherReq := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(jsonOtherLoginPayload))
	otherReq.Header.Set("Content-Type", "application/json")

	otherW := httptest.NewRecorder()

	router.ServeHTTP(otherW, otherReq)

	var otherLoginResponse map[string]string

	json.Unmarshal(otherW.Body.Bytes(), &otherLoginResponse)
	otherTokenString, otherExists := otherLoginResponse["token"]
	if !otherExists || otherTokenString == "" {
		t.Fatal("Other login response missing 'token' field")
	}

	// Store topicIDs for cleanup
	topicIDs := []int{}

	// Cleanup
	defer func() {
		clearTestData(t, repo, []string{testUsername, otherUsername}, topicIDs)
	}()

	// 1. Successful Topic Deletion
	t.Run("SuccessfulTopicDeletion", func(t *testing.T) {
		// Create test topic
		var topicID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO topics (title, description, created_by)
			VALUES ($1, $2, $3)
			RETURNING topic_id`,
			"Topic to be Deleted",
			"Topic Description",
			userID,
		).Scan(&topicID)

		if err != nil {
			t.Fatalf("Failed to create test topic: %v", err)
		}

		// Create test posts and comments under the topic
		for i := 0; i < 2; i++ {
			var postID int
			err = repo.DB.QueryRow(
				ctx,
				`INSERT INTO posts (topic_id, title, content, created_by)
				VALUES ($1, $2, $3, $4)
				RETURNING post_id`,
				topicID,
				fmt.Sprintf("Post %d for Deletion Test", i+1),
				"Post Content",
				userID,
			).Scan(&postID)

			if err != nil {
				t.Fatalf("Failed to create test post %d: %v", i+1, err)
			}

			for j := 0; j < 2; j++ {
				_, err = repo.DB.Exec(
					ctx,
					`INSERT INTO comments (post_id, content, created_by)
					VALUES ($1, $2, $3)`,
					postID,
					fmt.Sprintf("Comment %d for Post %d Deletion Test", j+1, i+1),
					userID,
				)

				if err != nil {
					t.Fatalf("Failed to create test comment %d for post %d: %v", j+1, i+1, err)
				}
			}
		}

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/topics/%d", topicID), nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("Expected status %d for topic deletion, got %d. Response: %s", http.StatusNoContent, w.Code, w.Body.String())
		}

		// Verify topic was deleted
		var topicCount int
		err = repo.DB.QueryRow(
			ctx,
			`SELECT COUNT(*) 
			FROM topics
			WHERE topic_id = $1`,
			topicID,
		).Scan(&topicCount)

		if err != nil {
			t.Fatalf("Failed to verify topic deletion: %v", err)
		}

		if topicCount != 0 {
			t.Fatalf("Topic with ID %d was not deleted", topicID)
		}

		// Verify posts were deleted
		var postCount int
		err = repo.DB.QueryRow(
			ctx,
			`SELECT COUNT(*) 
			FROM posts
			WHERE topic_id = $1`,
			topicID,
		).Scan(&postCount)

		if err != nil {
			t.Fatalf("Failed to verify posts deletion: %v", err)
		}

		if postCount != 0 {
			t.Fatalf("Posts under topic ID %d were not deleted", topicID)
		}

		// Verify comments were deleted
		var commentCount int
		err = repo.DB.QueryRow(
			ctx,
			`SELECT COUNT(*) 
			FROM comments
			WHERE post_id IN (SELECT post_id FROM posts WHERE topic_id = $1)`,
			topicID,
		).Scan(&commentCount)

		if err != nil {
			t.Fatalf("Failed to verify comments deletion: %v", err)
		}

		if commentCount != 0 {
			t.Fatalf("Comments under topic ID %d were not deleted", topicID)
		}
	})

	// 2. Topic Deletion without Authentication Token
	t.Run("TopicDeletionWithoutAuthenticationToken", func(t *testing.T) {
		// Create test topic
		var topicID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO topics (title, description, created_by)
			VALUES ($1, $2, $3)
			RETURNING topic_id`,
			"Topic for Unauthorized Deletion",
			"Topic Description",
			userID,
		).Scan(&topicID)

		if err != nil {
			t.Fatalf("Failed to create test topic: %v", err)
		}

		// Add topicID to cleanup list
		topicIDs = append(topicIDs, topicID)

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/topics/%d", topicID), nil)
		// No Authorization header

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status %d for unauthorized topic deletion, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// 3. Topic Deletion by Non-Owner
	t.Run("TopicDeletionByNonOwner", func(t *testing.T) {
		// Create test topic
		var topicID int
		err = repo.DB.QueryRow(
			ctx,
			`INSERT INTO topics (title, description, created_by)
			VALUES ($1, $2, $3)
			RETURNING topic_id`,
			"Topic for Non-Owner Deletion",
			"Topic Description",
			userID,
		).Scan(&topicID)

		if err != nil {
			t.Fatalf("Failed to create test topic: %v", err)
		}

		// Add topicID to cleanup list
		topicIDs = append(topicIDs, topicID)

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/topics/%d", topicID), nil)
		req.Header.Set("Authorization", "Bearer "+otherTokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Expected status %d for topic deletion by non-owner, got %d. Response: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// 4. Deletion of Non-Existent Topic
	t.Run("DeletionOfNonExistentTopic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/topics/%d", 9999999), nil) // Assuming this topic ID does not exist
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status %d for deletion of non-existent topic, got %d. Response: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	// 5. Topic Deletion with Invalid Topic ID
	t.Run("TopicDeletionWithInvalidTopicID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/topics/invalid_id", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status %d for topic deletion with invalid topic ID, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}
