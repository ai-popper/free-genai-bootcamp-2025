package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
	"github.com/lang-portal/backend/internal/models"
)

// GetStudyActivityByID handles GET /api/study_activities/:id
func GetStudyActivityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := database.GetStudyActivityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// GetActivityStudySessions handles GET /api/study_activities/:id/study_sessions
func GetActivityStudySessions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	sessions, total, err := database.GetActivityStudySessions(id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := models.PaginatedResponse{
		Items: sessions,
		Pagination: models.Pagination{
			CurrentPage:  page,
			ItemsPerPage: limit,
			TotalPages:   totalPages,
			TotalItems:   total,
		},
	}

	c.JSON(http.StatusOK, response)
}

// CreateStudySession handles POST /api/study_activities
func CreateStudySession(c *gin.Context) {
	var req struct {
		GroupID         int `json:"group_id" binding:"required"`
		StudyActivityID int `json:"study_activity_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	sessionID, err := database.CreateStudySession(req.GroupID, req.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.CreateSessionResponse{
		ID:      sessionID,
		GroupID: req.GroupID,
	}

	c.JSON(http.StatusOK, response)
}
