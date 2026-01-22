// backend/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/adzzfarr/gossip-with-go/backend/internal/api"
	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
)

// Run database migrations
func runMigrations() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL not set, skipping migration")
		return nil
	}

	if !strings.Contains(databaseURL, "sslmode=") {
		if strings.Contains(databaseURL, "?") {
			databaseURL += "&sslmode=require"
		} else {
			databaseURL += "?sslmode=require"
		}
	}

	// Migration instance
	m, err := migrate.New(
		"file://./migrations",
		databaseURL,
	)

	if err != nil {
		return err
	}

	defer m.Close()

	// Run migrations up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply.")
	} else {
		log.Println("Database migrations applied successfully.")
	}

	return nil
}

func main() {
	// Run migrations
	if err := runMigrations(); err != nil {
		log.Printf("Failed to run migrations: %v", err)
		log.Println("Continuing anyway; check if database is already migrated")
	}

	// Initialise database
	dbPool, err := data.OpenDB()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer dbPool.Close()

	log.Println("Database connection pool successfully created.")

	// Initialise Layers
	repo := data.NewRepository(dbPool)

	// Seed Handler
	seedHandler := api.NewSeedHandler(dbPool)

	// Topics
	topicService := service.NewTopicService(repo)
	topicHandler := api.NewTopicHandler(topicService)

	// Users
	userService := service.NewUserService(repo)
	userHandler := api.NewUserHandler(userService)

	// Posts
	postService := service.NewPostService(repo)
	postHandler := api.NewPostHandler(postService)

	// Comments
	commentService := service.NewCommentService(repo)
	commentHandler := api.NewCommentHandler(commentService)

	// Votes
	voteService := service.NewVoteService(repo)
	voteHandler := api.NewVoteHandler(voteService, postService, commentService)

	// JWT (Replace "secret-key" with a secure key from env variables in production)
	jwtService := service.NewJWTService("secret-key", 24*time.Hour) // 24 hours expiry

	// Login
	loginService := service.NewLoginService(repo)
	loginHandler := api.NewLoginHandler(loginService, jwtService)

	// Initialise Gin router
	router := gin.Default()

	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost:") {
				return true
			}

			if strings.HasSuffix(origin, ".vercel.app") {
				return true
			}

			if origin == "https://gossip-with-go-seven.vercel.app/" {
				return true
			}

			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health Check Endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Seed endpoint
	router.POST("/seed", seedHandler.SeedDatabase)

	// Register API Routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().Unix(),
				"version":   "1.0.0",
			})
		})

		// Public Routes - No Auth Required
		v1.POST("/register", userHandler.RegisterUser)
		v1.POST("/login", loginHandler.LoginUser)

		v1.GET("/topics", topicHandler.GetAllTopics)
		v1.GET("/topics/:topicID", topicHandler.GetTopicByID)

		// Optional Auth Routes - Guests can view, Authenticated users get extra info
		optional := v1.Group("")
		optional.Use(api.OptionalAuthMiddleware(jwtService))
		{
			optional.GET("/topics/:topicID/posts", postHandler.GetPostsByTopicID)
			optional.GET("/posts/:postID", postHandler.GetPostByID)
			optional.GET("/posts/:postID/comments", commentHandler.GetCommentsByPostID)
		}

		// Protected Routes (Auth Required)
		protected := v1.Group("")
		protected.Use(api.AuthMiddleware(jwtService))
		{
			// Topics
			protected.POST("/topics", topicHandler.CreateTopic)
			protected.PUT("/topics/:topicID", topicHandler.UpdateTopic)
			protected.DELETE("/topics/:topicID", topicHandler.DeleteTopic)

			// Posts
			protected.POST("/topics/:topicID/posts", postHandler.CreatePost)
			protected.PUT("/posts/:postID", postHandler.UpdatePost)
			protected.DELETE("/posts/:postID", postHandler.DeletePost)

			// Comments
			protected.POST("/posts/:postID/comments", commentHandler.CreateComment)
			protected.PUT("/comments/:commentID", commentHandler.UpdateComment)
			protected.DELETE("/comments/:commentID", commentHandler.DeleteComment)

			// Votes
			protected.POST("/posts/:postID/vote", voteHandler.VoteOnPost)
			protected.DELETE("/posts/:postID/vote", voteHandler.RemoveVoteFromPost)
			protected.POST("/comments/:commentID/vote", voteHandler.VoteOnComment)
			protected.DELETE("/comments/:commentID/vote", voteHandler.RemoveVoteFromComment)

			// User Profiles
			protected.GET("/users/:id", userHandler.GetUserByID)
			protected.GET("/users/:id/posts", userHandler.GetUserPosts)
			protected.GET("/users/:id/comments", userHandler.GetUserComments)
		}
	}

	// Run Server
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
