package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
	"github.com/lang-portal/backend/internal/models"
)

// ResetHistory handles POST /api/reset_history
func ResetHistory(c *gin.Context) {
	err := database.ResetHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.SuccessResponse{
		Success: true,
		Message: "Study history reset successfully",
	}

	c.JSON(http.StatusOK, response)
}

// FullReset handles POST /api/full_reset
func FullReset(c *gin.Context) {
	migrationsDir := filepath.Join("db", "migrations")
	err := database.FullReset(migrationsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.SuccessResponse{
		Success: true,
		Message: "Full database reset completed successfully",
	}

	c.JSON(http.StatusOK, response)
}
