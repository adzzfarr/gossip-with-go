// Run `go test -v ./internal/data` in /backend

package data

import (
	"context"
	"testing"
)

// Test database connection and repository function
func TestGetAllTopicsIntegration(t *testing.T) {
	// 1. Establish DB Connection and Repository
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close() // Ensure DB connection is closed after test

	repo := NewRepository(db)
	ctx := context.Background()

	// 2. Insert Test Data (Simulate new user and a topic)
	// Need to create user first due to user_id foreign key constraint on the 'topics' table
	var userID int
	userSQL := "INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id"
	err = db.QueryRow(ctx, userSQL, "testuser", "hashed-pass-123").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}
	t.Logf("Inserted test user with ID: %d", userID)

	// Insert topic using userID
	topicSQL := "INSERT INTO topics (title, description, created_by, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())"
	_, err = db.Exec(ctx, topicSQL, "Test Topic Title", "Test Description", userID)
	if err != nil {
		t.Fatalf("Failed to insert test topic: %v", err)
	}

	// 3. Test Repository Function
	topics, err := repo.GetAllTopics()
	if err != nil {
		t.Errorf("GetAllTopics failed with error: %v", err)
	}

	if len(topics) == 0 {
		t.Fatalf("Expected at least 1 topic, but got 0")
	}

	// Find our test topic
	var foundTopic *Topic
	for i := range topics {
		if topics[i].Title == "Test Topic Title" {
			foundTopic = topics[i]
			break
		}
	}

	if foundTopic == nil {
		t.Fatalf("Test topic not found in results")
	}

	if foundTopic.Title != "Test Topic Title" {
		t.Errorf("Title mismatch: Expected 'Test Topic Title', got '%s'", foundTopic.Title)
	}

	if foundTopic.Description != "Test Description" {
		t.Errorf("Description mismatch: Expected 'Test Description', got '%s'", foundTopic.Description)
	}

	if foundTopic.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}

	t.Log("Repository function successfully executed and data verified.")

	// 5. Delete Test Data
	// Delete topic first due to user_id foreign key constraint
	_, err = db.Exec(ctx, "DELETE FROM topics WHERE created_by = $1", userID)
	if err != nil {
		t.Fatalf("Failed to delete test topics: %v", err)
	}

	// Delete user
	_, err = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}

// TestRegisterUser tests user registration
func TestRegisterUser(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	testUsername := "test_register_user"
	testPasswordHash := "hashed_password_123"

	// Cleanup before test
	_, _ = db.Exec(ctx, "DELETE FROM users WHERE username = $1", testUsername)

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE username = $1", testUsername)
	}()

	t.Run("successful registration", func(t *testing.T) {
		user, err := repo.CreateUser(&User{
			Username:     testUsername,
			PasswordHash: testPasswordHash,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.UserID == 0 {
			t.Error("expected user_id to be set")
		}
		if user.Username != testUsername {
			t.Errorf("expected username %s, got %s", testUsername, user.Username)
		}
		if user.PasswordHash != testPasswordHash {
			t.Errorf("expected password hash %s, got %s", testPasswordHash, user.PasswordHash)
		}
		if user.CreatedAt.IsZero() {
			t.Error("expected created_at to be set")
		}
		if user.UpdatedAt.IsZero() {
			t.Error("expected updated_at to be set")
		}
	})

	// Test: Duplicate username should fail
	t.Run("duplicate username", func(t *testing.T) {
		_, err := repo.CreateUser(&User{
			Username:     testUsername,
			PasswordHash: "another_hash",
		})
		if err == nil {
			t.Error("expected error for duplicate username, got nil")
		}
	})
}

// TestGetUserByUsername tests retrieving a user by username
func TestGetUserByUsername(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	testUsername := "test_get_user"
	testPasswordHash := "hashed_password_456"

	// Cleanup and insert test user
	_, _ = db.Exec(ctx, "DELETE FROM users WHERE username = $1", testUsername)

	var userID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id",
		testUsername,
		testPasswordHash,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}()

	// Test: Successfully retrieve existing user
	t.Run("retrieve existing user", func(t *testing.T) {
		user, err := repo.GetUserByUsername(testUsername)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Username != testUsername {
			t.Errorf("expected username %s, got %s", testUsername, user.Username)
		}
		if user.PasswordHash != testPasswordHash {
			t.Errorf("expected password hash to match")
		}
		if user.CreatedAt.IsZero() {
			t.Error("expected created_at to be set")
		}
		if user.UpdatedAt.IsZero() {
			t.Error("expected updated_at to be set")
		}
	})

	// Test: Non-existent user should return error
	t.Run("non-existent user", func(t *testing.T) {
		_, err := repo.GetUserByUsername("nonexistent_user_12345")
		if err == nil {
			t.Error("expected error for non-existent user, got nil")
		}
	})
}

// TestCreateTopic tests topic creation
func TestCreateTopic(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	// Create test user
	testUsername := "test_create_topic_user"
	var userID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id",
		testUsername,
		"hash123",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM topics WHERE created_by = $1", userID)
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}()

	// Test: Successful topic creation
	t.Run("successful topic creation", func(t *testing.T) {
		topic, err := repo.CreateTopic("Test Topic", "Test Description", userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if topic.TopicID == 0 {
			t.Error("expected topic_id to be set")
		}
		if topic.Title != "Test Topic" {
			t.Errorf("expected title 'Test Topic', got '%s'", topic.Title)
		}
		if topic.Description != "Test Description" {
			t.Errorf("expected description 'Test Description', got '%s'", topic.Description)
		}
		if topic.CreatedBy != userID {
			t.Errorf("expected created_by %d, got %d", userID, topic.CreatedBy)
		}
		if topic.CreatedAt.IsZero() {
			t.Error("expected created_at to be set")
		}
		if topic.UpdatedAt.IsZero() {
			t.Error("expected updated_at to be set")
		}
	})
}

// TestCreatePost tests post creation
func TestCreatePost(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	// Create test user
	testUsername := "test_create_post_user"
	var userID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id",
		testUsername,
		"hash123",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO topics (title, description, created_by, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING topic_id",
		"Test Topic",
		"Test Description",
		userID,
	).Scan(&topicID)
	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM posts WHERE created_by = $1", userID)
		_, _ = db.Exec(ctx, "DELETE FROM topics WHERE topic_id = $1", topicID)
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}()

	// Test: Successful post creation
	t.Run("successful post creation", func(t *testing.T) {
		post, err := repo.CreatePost(topicID, "Test Post", "Test Content", userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if post.PostID == 0 {
			t.Error("expected post_id to be set")
		}
		if post.TopicID != topicID {
			t.Errorf("expected topic_id %d, got %d", topicID, post.TopicID)
		}
		if post.Title != "Test Post" {
			t.Errorf("expected title 'Test Post', got '%s'", post.Title)
		}
		if post.Content != "Test Content" {
			t.Errorf("expected content 'Test Content', got '%s'", post.Content)
		}
		if post.CreatedBy != userID {
			t.Errorf("expected created_by %d, got %d", userID, post.CreatedBy)
		}
		if post.CreatedAt.IsZero() {
			t.Error("expected created_at to be set")
		}
		if post.UpdatedAt.IsZero() {
			t.Error("expected updated_at to be set")
		}
	})
}

// TestGetPostsByTopicID tests retrieving posts by topic ID
func TestGetPostsByTopicID(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	// Create test user
	testUsername := "test_get_posts_user"
	var userID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id",
		testUsername,
		"hash123",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO topics (title, description, created_by, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING topic_id",
		"Test Topic",
		"Test Description",
		userID,
	).Scan(&topicID)
	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING post_id",
		topicID,
		"Test Post",
		"Test Content",
		userID,
	).Scan(&postID)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM posts WHERE post_id = $1", postID)
		_, _ = db.Exec(ctx, "DELETE FROM topics WHERE topic_id = $1", topicID)
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}()

	// Test: Successfully retrieve posts
	t.Run("retrieve posts by topic", func(t *testing.T) {
		posts, err := repo.GetPostsByTopicID(topicID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(posts) == 0 {
			t.Fatal("expected at least 1 post, got 0")
		}

		found := false
		for _, post := range posts {
			if post.PostID == postID {
				found = true
				if post.Title != "Test Post" {
					t.Errorf("expected title 'Test Post', got '%s'", post.Title)
				}
				if post.UpdatedAt.IsZero() {
					t.Error("expected updated_at to be set")
				}
				break
			}
		}

		if !found {
			t.Error("test post not found in results")
		}
	})

	// Test: Non-existent topic
	t.Run("non-existent topic", func(t *testing.T) {
		posts, err := repo.GetPostsByTopicID(999999)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(posts) != 0 {
			t.Errorf("expected 0 posts for non-existent topic, got %d", len(posts))
		}
	})
}

// TestCreateComment tests comment creation
func TestCreateComment(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	// Create test user
	testUsername := "test_create_comment_user"
	var userID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id",
		testUsername,
		"hash123",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO topics (title, description, created_by, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING topic_id",
		"Test Topic",
		"Test Description",
		userID,
	).Scan(&topicID)
	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING post_id",
		topicID,
		"Test Post",
		"Test Content",
		userID,
	).Scan(&postID)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM comments WHERE created_by = $1", userID)
		_, _ = db.Exec(ctx, "DELETE FROM posts WHERE post_id = $1", postID)
		_, _ = db.Exec(ctx, "DELETE FROM topics WHERE topic_id = $1", topicID)
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}()

	// Test: Successful comment creation
	t.Run("successful comment creation", func(t *testing.T) {
		comment, err := repo.CreateComment(postID, "Test Comment", userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if comment.CommentID == 0 {
			t.Error("expected comment_id to be set")
		}
		if comment.PostID != postID {
			t.Errorf("expected post_id %d, got %d", postID, comment.PostID)
		}
		if comment.Content != "Test Comment" {
			t.Errorf("expected content 'Test Comment', got '%s'", comment.Content)
		}
		if comment.CreatedBy != userID {
			t.Errorf("expected created_by %d, got %d", userID, comment.CreatedBy)
		}
		if comment.CreatedAt.IsZero() {
			t.Error("expected created_at to be set")
		}
		if comment.UpdatedAt.IsZero() {
			t.Error("expected updated_at to be set")
		}
	})
}

// TestGetCommentsByPostID tests retrieving comments by post ID
func TestGetCommentsByPostID(t *testing.T) {
	// Setup
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	// Create test user
	testUsername := "test_get_comments_user"
	var userID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING user_id",
		testUsername,
		"hash123",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test topic
	var topicID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO topics (title, description, created_by, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING topic_id",
		"Test Topic",
		"Test Description",
		userID,
	).Scan(&topicID)
	if err != nil {
		t.Fatalf("Failed to create test topic: %v", err)
	}

	// Create test post
	var postID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING post_id",
		topicID,
		"Test Post",
		"Test Content",
		userID,
	).Scan(&postID)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Create test comment
	var commentID int
	err = db.QueryRow(
		ctx,
		"INSERT INTO comments (post_id, content, created_by, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING comment_id",
		postID,
		"Test Comment",
		userID,
	).Scan(&commentID)
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	// Cleanup after test
	defer func() {
		_, _ = db.Exec(ctx, "DELETE FROM comments WHERE comment_id = $1", commentID)
		_, _ = db.Exec(ctx, "DELETE FROM posts WHERE post_id = $1", postID)
		_, _ = db.Exec(ctx, "DELETE FROM topics WHERE topic_id = $1", topicID)
		_, _ = db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}()

	// Test: Successfully retrieve comments
	t.Run("retrieve comments by post", func(t *testing.T) {
		comments, err := repo.GetCommentsByPostID(postID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(comments) == 0 {
			t.Fatal("expected at least 1 comment, got 0")
		}

		found := false
		for _, comment := range comments {
			if comment.CommentID == commentID {
				found = true
				if comment.Content != "Test Comment" {
					t.Errorf("expected content 'Test Comment', got '%s'", comment.Content)
				}
				if comment.UpdatedAt.IsZero() {
					t.Error("expected updated_at to be set")
				}
				break
			}
		}

		if !found {
			t.Error("test comment not found in results")
		}
	})

	// Test: Non-existent post
	t.Run("non-existent post", func(t *testing.T) {
		comments, err := repo.GetCommentsByPostID(999999)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(comments) != 0 {
			t.Errorf("expected 0 comments for non-existent post, got %d", len(comments))
		}
	})
}
