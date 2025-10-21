package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
	"github.com/lang-portal/backend/internal/models"
)

// GetGroups handles GET /api/groups
func GetGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	groups, total, err := database.GetGroups(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + limit - 1) / limit

	response := models.PaginatedResponse{
		Items: groups,
		Pagination: models.Pagination{
			CurrentPage:  page,
			ItemsPerPage: limit,
			TotalPages:   totalPages,
			TotalItems:   total,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetGroupByID handles GET /api/groups/:id
func GetGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := database.GetGroupByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords handles GET /api/groups/:id/words
func GetGroupWords(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	words, total, err := database.GetGroupWords(id, page, limit)
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

// GetGroupStudySessions handles GET /api/groups/:id/study_sessions
func GetGroupStudySessions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	sessions, total, err := database.GetGroupStudySessions(id, page, limit)
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
