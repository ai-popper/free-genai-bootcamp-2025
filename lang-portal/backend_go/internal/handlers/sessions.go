package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
	"github.com/lang-portal/backend/internal/models"
)

// GetStudySessions handles GET /api/study_sessions
func GetStudySessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	sessions, total, err := database.GetStudySessions(page, limit)
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

// GetStudySessionByID handles GET /api/study_sessions/:id
func GetStudySessionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := database.GetStudySessionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudySessionWords handles GET /api/study_sessions/:id/words
func GetStudySessionWords(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	words, total, err := database.GetStudySessionWords(id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := models.PaginatedResponse{
		Items: words,
		Pagination: models.Pagination{
			CurrentPage:  page,
			ItemsPerPage: limit,
			TotalPages:   totalPages,
			TotalItems:   total,
		},
	}

	c.JSON(http.StatusOK, response)
}

// RecordWordReview handles POST /api/study_sessions/:id/words/:word_id/review
func RecordWordReview(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	wordID, err := strconv.Atoi(c.Param("word_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var req struct {
		Correct bool `json:"correct"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	review, err := database.CreateWordReview(sessionID, wordID, req.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ReviewResponse{
		Success:        true,
		WordID:         review.WordID,
		StudySessionID: review.StudySessionID,
		Correct:        review.Correct,
		CreatedAt:      review.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}
