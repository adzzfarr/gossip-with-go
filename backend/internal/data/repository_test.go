/*
1. Run `docker compose up -d db` from project root to start PostgreSQL container (if not already running)
2. Run `go test -v ./internal/data/ -run TestGetAllTopicsIntegration` from backend/internal/data/

Expected Output:
=== RUN   TestGetAllTopicsIntegration
--- PASS: TestGetAllTopicsIntegration

	repository_test.go:27: Inserted test user with ID: (some number)
	repository_test.go:49: Repository function successfully executed and data verified.

PASS
ok      github.com/adzzfarr/gossip-with-go/backend/internal/data
*/

package data

import (
	"context"
	"testing"
	"time"
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
	userSQL := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING user_id"
	err = db.QueryRow(ctx, userSQL, "testuser", "hashed-pass-123").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}
	t.Logf("Inserted test user with ID: %d", userID)

	// Insert topic using userID
	topicSQL := "INSERT INTO topics (title, description, created_by, created_at) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(ctx, topicSQL, "Test Topic Title", "Test Description", userID, time.Now())
	if err != nil {
		t.Fatalf("Failed to insert test topic: %v", err)
	}

	// 3. Test Repository Function
	topics, err := repo.GetAllTopics()
	if err != nil {
		t.Errorf("GetAllTopics failed with error: %v", err)
	}

	if len(topics) != 1 {
		t.Errorf("Expected 1 topic, but got %d", len(topics))
	} else {
		if topics[0].Title != "Test Topic Title" {
			t.Errorf("Title mismatch: Expected 'Test Topic Title', got '%s'", topics[0].Title)
		}
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
