package handlers

import (
	"net/http"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	shortCode, err := services.CreateShortURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + shortCode})
}

func Redirect(c *gin.Context) {
	code := c.Param("code")
	originalURL, err := services.GetOriginalURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
