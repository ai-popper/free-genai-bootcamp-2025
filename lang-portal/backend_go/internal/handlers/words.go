package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
	"github.com/lang-portal/backend/internal/models"
)

// GetWords handles GET /api/words
func GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit := 100
	words, total, err := database.GetWords(page, limit)
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

// GetWordByID handles GET /api/words/:id
func GetWordByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := database.GetWordByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	c.JSON(http.StatusOK, word)
}
