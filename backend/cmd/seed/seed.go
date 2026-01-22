package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Run seeds the database with test data
func Run(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create test users
	log.Println("🌱 Creating test users...")
	userIDs, err := createUsers(ctx, pool, 10)
	if err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}
	log.Printf("   ✓ Created %d users\n", len(userIDs))

	// Create topics
	log.Println("🌱 Creating topics...")
	topicIDs, err := createTopics(ctx, pool, userIDs, 15)
	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}
	log.Printf("   ✓ Created %d topics\n", len(topicIDs))

	// Create posts for each topic
	log.Println("🌱 Creating posts...")
	postCount := 0
	for i, topicID := range topicIDs {
		postIDs, err := createPosts(ctx, pool, topicID, userIDs, 20)
		if err != nil {
			return fmt.Errorf("failed to create posts: %w", err)
		}
		postCount += len(postIDs)

		// Create comments for each post
		commentCount := 0
		for _, postID := range postIDs {
			count, err := createComments(ctx, pool, postID, userIDs, 25)
			if err != nil {
				return fmt.Errorf("failed to create comments: %w", err)
			}
			commentCount += count
		}

		if (i+1)%5 == 0 {
			log.Printf("   ✓ Progress: %d/%d topics processed\n", i+1, len(topicIDs))
		}
	}
	log.Printf("   ✓ Created %d total posts\n", postCount)

	// Add some votes
	log.Println("🌱 Adding votes...")
	if err := addVotes(ctx, pool, userIDs); err != nil {
		return fmt.Errorf("failed to add votes: %w", err)
	}
	log.Println("   ✓ Votes added")

	return nil
}

func createUsers(ctx context.Context, pool *pgxpool.Pool, count int) ([]int, error) {
	userIDs := make([]int, 0, count)

	usernames := []string{
		"alice_wonder", "bob_builder", "charlie_brown", "diana_prince",
		"evan_almighty", "fiona_apple", "george_curious", "hannah_montana",
		"isaac_newton", "julia_child",
	}

	// Hash a simple password for testing (Password123)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	for i := 0; i < count && i < len(usernames); i++ {
		var userID int
		err := pool.QueryRow(ctx, `
            INSERT INTO users (username, password_hash, created_at, updated_at)
            VALUES ($1, $2, NOW(), NOW())
            ON CONFLICT (username) DO NOTHING
            RETURNING user_id
        `, usernames[i], string(hashedPassword)).Scan(&userID)

		if err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", usernames[i], err)
		}

		if userID != 0 {
			userIDs = append(userIDs, userID)
		}
	}

	return userIDs, nil
}

func createTopics(ctx context.Context, pool *pgxpool.Pool, userIDs []int, count int) ([]int, error) {
	topicIDs := make([]int, 0, count)

	topics := []struct {
		title       string
		description string
	}{
		{"Technology Trends 2026", "Discuss the latest innovations in AI, blockchain, and quantum computing"},
		{"Cooking Tips & Recipes", "Share your favorite recipes and cooking techniques"},
		{"Travel Adventures", "Stories and tips from around the world"},
		{"Gaming Community", "All things video games - reviews, tips, and multiplayer sessions"},
		{"Book Club", "Monthly book discussions and recommendations"},
		{"Fitness & Health", "Workout routines, nutrition advice, and wellness tips"},
		{"Movie Reviews", "Latest releases and classic film discussions"},
		{"Music Production", "Production techniques, gear reviews, and collaboration"},
		{"Photography Showcase", "Share your photos and get feedback"},
		{"DIY & Home Improvement", "Projects, tools, and home renovation ideas"},
		{"Pet Lovers", "Cute pet photos and care advice"},
		{"Career Advice", "Job hunting, interviews, and professional development"},
		{"Science & Space", "Latest discoveries and space exploration news"},
		{"Art & Design", "Creative works, techniques, and inspiration"},
		{"Cryptocurrency", "Trading strategies, market analysis, and blockchain tech"},
	}

	for i := 0; i < count; i++ {
		topic := topics[i%len(topics)]
		userID := userIDs[rand.Intn(len(userIDs))]

		// Random creation time in the past 30 days
		daysAgo := rand.Intn(30)
		createdAt := time.Now().AddDate(0, 0, -daysAgo)

		var topicID int
		err := pool.QueryRow(ctx, `
            INSERT INTO topics (title, description, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $4)
            RETURNING topic_id
        `, topic.title, topic.description, userID, createdAt).Scan(&topicID)

		if err != nil {
			return nil, fmt.Errorf("failed to create topic %s: %w", topic.title, err)
		}
		topicIDs = append(topicIDs, topicID)
	}

	return topicIDs, nil
}

func createPosts(ctx context.Context, pool *pgxpool.Pool, topicID int, userIDs []int, count int) ([]int, error) {
	postIDs := make([]int, 0, count)

	postTitles := []string{
		"Getting Started Guide", "My Experience So Far", "Tips and Tricks",
		"Common Mistakes to Avoid", "Advanced Techniques", "Beginner's Questions",
		"Discussion: What do you think?", "Tutorial: Step by Step", "Review and Analysis",
		"Comparison: A vs B", "Latest Updates", "Personal Story",
		"Help Needed!", "Amazing Discovery", "Controversial Opinion",
		"Meta Discussion", "Resource Collection", "Quick Question",
		"Detailed Breakdown", "Future Predictions",
	}

	postContent := "This is a sample post with detailed content. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	for i := 0; i < count; i++ {
		title := postTitles[i%len(postTitles)]
		userID := userIDs[rand.Intn(len(userIDs))]

		hoursAgo := rand.Intn(24 * 7)
		createdAt := time.Now().Add(-time.Duration(hoursAgo) * time.Hour)

		var postID int
		err := pool.QueryRow(ctx, `
            INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $5)
            RETURNING post_id
        `, topicID, title, postContent, userID, createdAt).Scan(&postID)

		if err != nil {
			return nil, fmt.Errorf("failed to create post: %w", err)
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func createComments(ctx context.Context, pool *pgxpool.Pool, postID int, userIDs []int, count int) (int, error) {
	comments := []string{
		"Great post! This really helped me.",
		"Thanks for sharing!",
		"Interesting perspective.",
		"I disagree, but well written.",
		"Can you elaborate more?",
		"This is exactly what I needed.",
		"Have you considered this alternative?",
		"Well researched!",
		"I tried this and it worked.",
		"Could you provide more details?",
	}

	for i := 0; i < count; i++ {
		content := comments[i%len(comments)]
		userID := userIDs[rand.Intn(len(userIDs))]

		minutesAgo := rand.Intn(60 * 24)
		createdAt := time.Now().Add(-time.Duration(minutesAgo) * time.Minute)

		_, err := pool.Exec(ctx, `
            INSERT INTO comments (post_id, content, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $4)
        `, postID, content, userID, createdAt)

		if err != nil {
			return 0, fmt.Errorf("failed to create comment: %w", err)
		}
	}

	return count, nil
}

func addVotes(ctx context.Context, pool *pgxpool.Pool, userIDs []int) error {
	// Get all posts
	rows, err := pool.Query(ctx, "SELECT post_id FROM posts")
	if err != nil {
		return err
	}
	defer rows.Close()

	postIDs := []int{}
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return err
		}
		postIDs = append(postIDs, postID)
	}

	// Add random votes to posts (30% chance per user)
	for _, postID := range postIDs {
		for _, userID := range userIDs {
			if rand.Float32() < 0.3 {
				voteType := "upvote"
				if rand.Float32() < 0.2 {
					voteType = "downvote"
				}

				_, err := pool.Exec(ctx, `
                    INSERT INTO votes (user_id, post_id, vote_type, created_at)
                    VALUES ($1, $2, $3, NOW())
                    ON CONFLICT (user_id, post_id) DO NOTHING
                `, userID, postID, voteType)

				if err != nil {
					return err
				}
			}
		}
	}

	// Get all comments
	rows, err = pool.Query(ctx, "SELECT comment_id FROM comments")
	if err != nil {
		return err
	}
	defer rows.Close()

	commentIDs := []int{}
	for rows.Next() {
		var commentID int
		if err := rows.Scan(&commentID); err != nil {
			return err
		}
		commentIDs = append(commentIDs, commentID)
	}

	// Add random votes to comments (20% chance per user)
	for _, commentID := range commentIDs {
		for _, userID := range userIDs {
			if rand.Float32() < 0.2 {
				voteType := "upvote"
				if rand.Float32() < 0.15 {
					voteType = "downvote"
				}

				_, err := pool.Exec(ctx, `
                    INSERT INTO votes (user_id, comment_id, vote_type, created_at)
                    VALUES ($1, $2, $3, NOW())
                    ON CONFLICT (user_id, comment_id) DO NOTHING
                `, userID, commentID, voteType)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
