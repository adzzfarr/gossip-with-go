// backend/cmd/server/main.go
package main

import (
	"log"
	"net/http"

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
	defer dbPool.Close() // Close connection pool when done

	log.Println("Database connection pool successfully created.")

	// Initialise Layers (Dependency Injection)
	// Flow: main -> Repository -> Service -> Handler
	repo := data.NewRepository(dbPool)

	// Topics
	topicService := service.NewTopicService(repo)
	topicHandler := api.NewTopicHandler(topicService)

	// Users
	userService := service.NewUserService(repo)
	userHandler := api.NewUserHandler(userService)

	// Initialise Gin router
	router := gin.Default()

	// Health Check Endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Register API Routes
	v1 := router.Group("/api/v1")
	{
		// Topic Routes
		v1.GET("/topics", topicHandler.GetAllTopics)

		// User Routes
		v1.POST("/users", userHandler.RegisterUser)
	}

	// Run Server
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
