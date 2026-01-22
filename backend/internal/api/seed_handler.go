package api

import (
	"net/http"

	"github.com/adzzfarr/gossip-with-go/backend/cmd/seed"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedHandler struct {
	pool *pgxpool.Pool
}

func NewSeedHandler(pool *pgxpool.Pool) *SeedHandler {
	return &SeedHandler{pool: pool}
}

func (h *SeedHandler) SeedDatabase(c *gin.Context) {
	// Add a simple secret key check for security

	if err := seed.Run(h.pool); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to seed database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database seeded successfully",
		"status":  "success",
	})
}
