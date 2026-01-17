package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Connect to database
	dbURL := "postgres://user:password@localhost:5432/forum_db?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	log.Println("Connected to database successfully")

	// Seed data
	if err := seedData(pool); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("✅ Database seeded successfully!")
}

func seedData(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create test users
	log.Println("Creating test users...")
	userIDs, err := createUsers(ctx, pool, 10)
	if err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}
	log.Printf("Created %d users\n", len(userIDs))

	// Create topics
	log.Println("Creating topics...")
	topicIDs, err := createTopics(ctx, pool, userIDs, 15)
	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}
	log.Printf("Created %d topics\n", len(topicIDs))

	// Create posts for each topic
	log.Println("Creating posts...")
	postCount := 0
	for _, topicID := range topicIDs {
		postIDs, err := createPosts(ctx, pool, topicID, userIDs, 20)
		if err != nil {
			return fmt.Errorf("failed to create posts: %w", err)
		}
		postCount += len(postIDs)

		// Create comments for each post
		log.Printf("Creating comments for topic %d...\n", topicID)
		commentCount := 0
		for _, postID := range postIDs {
			count, err := createComments(ctx, pool, postID, userIDs, 25)
			if err != nil {
				return fmt.Errorf("failed to create comments: %w", err)
			}
			commentCount += count
		}
		log.Printf("  Created %d comments for %d posts\n", commentCount, len(postIDs))
	}
	log.Printf("Created %d total posts\n", postCount)

	// Add some votes
	log.Println("Adding votes...")
	if err := addVotes(ctx, pool, userIDs); err != nil {
		return fmt.Errorf("failed to add votes: %w", err)
	}

	return nil
}

func createUsers(ctx context.Context, pool *pgxpool.Pool, count int) ([]int, error) {
	userIDs := make([]int, 0, count)

	usernames := []string{
		"alice_wonder", "bob_builder", "charlie_brown", "diana_prince",
		"evan_almighty", "fiona_apple", "george_curious", "hannah_montana",
		"isaac_newton", "julia_child",
	}

	for i := 0; i < count && i < len(usernames); i++ {
		var userID int
		err := pool.QueryRow(ctx, `
            INSERT INTO users (username, password_hash, created_at, updated_at)
            VALUES ($1, $2, NOW(), NOW())
            ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
            RETURNING user_id
        `, usernames[i], "$2a$10$dummyhashedpassword").Scan(&userID)

		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
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
			return nil, err
		}
		topicIDs = append(topicIDs, topicID)
	}

	return topicIDs, nil
}

func createPosts(ctx context.Context, pool *pgxpool.Pool, topicID int, userIDs []int, count int) ([]int, error) {
	postIDs := make([]int, 0, count)

	postTitles := []string{
		"Getting Started Guide",
		"My Experience So Far",
		"Tips and Tricks",
		"Common Mistakes to Avoid",
		"Advanced Techniques",
		"Beginner's Questions",
		"Discussion: What do you think?",
		"Tutorial: Step by Step",
		"Review and Analysis",
		"Comparison: A vs B",
		"Latest Updates",
		"Personal Story",
		"Help Needed!",
		"Amazing Discovery",
		"Controversial Opinion",
		"Meta Discussion",
		"Resource Collection",
		"Quick Question",
		"Detailed Breakdown",
		"Future Predictions",
	}

	postContents := []string{
		"This is a comprehensive guide covering all the basics you need to know. I've spent months researching this topic and wanted to share my findings with the community.",
		"After trying this for several weeks, here are my thoughts and observations. Overall, I'm quite impressed with the results.",
		"Here are some useful tips that I wish I knew when I started. These will save you a lot of time and frustration.",
		"Learn from my mistakes! Here are the most common pitfalls and how to avoid them.",
		"For those who are already familiar with the basics, here are some advanced strategies to take your skills to the next level.",
		"I'm new to this and have some questions. Any advice would be appreciated!",
		"I'd love to hear everyone's thoughts on this topic. What's your take?",
		"Follow these steps carefully and you'll get great results. Let me know if you have any questions!",
		"Here's my detailed analysis after extensive testing and comparison with alternatives.",
		"I've compared these two options side by side. Here's what I found.",
		"Just saw the latest update and wanted to share the highlights with everyone.",
		"Here's my personal journey and what I learned along the way.",
		"I'm stuck on this problem and could really use some help from the community.",
		"I just discovered something incredible that I have to share!",
		"This might be unpopular, but here's what I really think about this situation.",
		"Let's discuss how we can improve this community and make it better for everyone.",
		"I've compiled a list of the best resources I've found. Hope this helps!",
		"Quick question for the experts here - what's the best approach for this?",
		"Let me break this down into simple terms so everyone can understand.",
		"Based on current trends, here's what I think will happen next.",
	}

	for i := 0; i < count; i++ {
		title := postTitles[i%len(postTitles)]
		content := postContents[i%len(postContents)]
		userID := userIDs[rand.Intn(len(userIDs))]

		// Random creation time (within topic's lifetime)
		hoursAgo := rand.Intn(24 * 7) // Within last week
		createdAt := time.Now().Add(-time.Duration(hoursAgo) * time.Hour)

		var postID int
		err := pool.QueryRow(ctx, `
            INSERT INTO posts (topic_id, title, content, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $5)
            RETURNING post_id
        `, topicID, title, content, userID, createdAt).Scan(&postID)

		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func createComments(ctx context.Context, pool *pgxpool.Pool, postID int, userIDs []int, count int) (int, error) {
	comments := []string{
		"Great post! This really helped me understand the topic better.",
		"Thanks for sharing! I've been looking for something like this.",
		"Interesting perspective. I never thought about it that way.",
		"I disagree with some points, but overall well written.",
		"Can you elaborate more on this part?",
		"This is exactly what I needed. Thank you!",
		"Have you considered this alternative approach?",
		"Well researched and presented. Keep up the good work!",
		"I tried this and it worked perfectly. Thanks!",
		"Could you provide more details about the implementation?",
		"This doesn't work for me. Am I doing something wrong?",
		"Excellent breakdown of a complex topic.",
		"I have a follow-up question about this.",
		"This is similar to what I experienced. Good to know I'm not alone!",
		"Very informative. Learned a lot from this post.",
		"I'm not sure I agree, but I respect your opinion.",
		"This should be in the FAQ. Very helpful!",
		"Can you recommend any additional resources on this?",
		"I had the same problem and here's what worked for me...",
		"Thanks for taking the time to write this up!",
		"This is gold. Saving for future reference.",
		"Mind blown! Never knew about this.",
		"Simple and effective. Love it.",
		"This needs more upvotes. Underrated post!",
		"Following this thread. Very interested to see where this goes.",
	}

	for i := 0; i < count; i++ {
		content := comments[i%len(comments)]
		userID := userIDs[rand.Intn(len(userIDs))]

		// Random creation time
		minutesAgo := rand.Intn(60 * 24) // Within last day
		createdAt := time.Now().Add(-time.Duration(minutesAgo) * time.Minute)

		_, err := pool.Exec(ctx, `
            INSERT INTO comments (post_id, content, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $4)
        `, postID, content, userID, createdAt)

		if err != nil {
			return 0, err
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

	// Add random votes to posts
	for _, postID := range postIDs {
		// 30% of users vote on each post
		for _, userID := range userIDs {
			if rand.Float32() < 0.3 {
				voteType := 1
				if rand.Float32() < 0.2 { // 20% downvotes
					voteType = -1
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

	// Add random votes to comments
	for _, commentID := range commentIDs {
		// 20% of users vote on each comment
		for _, userID := range userIDs {
			if rand.Float32() < 0.2 {
				voteType := 1
				if rand.Float32() < 0.15 { // 15% downvotes
					voteType = -1
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
