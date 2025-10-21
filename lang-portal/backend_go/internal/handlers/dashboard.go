package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
)

// GetLastStudySession handles GET /api/dashboard/last_study_session
func GetLastStudySession(c *gin.Context) {
	session, err := database.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No study sessions found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudyProgress handles GET /api/dashboard/study_progress
func GetStudyProgress(c *gin.Context) {
	progress, err := database.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickStats handles GET /api/dashboard/quick_stats
func GetQuickStats(c *gin.Context) {
	stats, err := database.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
