// backend/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/adzzfarr/gossip-with-go/backend/internal/api"
	"github.com/adzzfarr/gossip-with-go/backend/internal/data"
	"github.com/adzzfarr/gossip-with-go/backend/internal/service"
)

func main() {
	// Initialise database
	dbPool, err := data.OpenDB()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer dbPool.Close()

	log.Println("Database connection pool successfully created.")

	// Initialise Layers
	repo := data.NewRepository(dbPool)

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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health Check Endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Register API Routes
	v1 := router.Group("/api/v1")
	{
		// Public Routes (No Auth Required)
		v1.POST("/users", userHandler.RegisterUser)
		v1.POST("/login", loginHandler.LoginUser)

		v1.GET("/topics", topicHandler.GetAllTopics)
		v1.GET("/topics/:topicID", topicHandler.GetTopicByID)

		v1.GET("/topics/:topicID/posts", postHandler.GetPostsByTopicID)
		v1.GET("/posts/:postID", postHandler.GetPostByID)

		v1.GET("/posts/:postID/comments", commentHandler.GetCommentsByPostID)

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
