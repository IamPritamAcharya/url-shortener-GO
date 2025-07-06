package handlers

import (
	"net/http"
	"os"
	"strings"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type CustomShortenRequest struct {
	URL        string `json:"url" binding:"required"`
	CustomCode string `json:"custom_code" binding:"required"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	Code        string `json:"code"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return strings.TrimSuffix(baseURL, "/")
}


func ShortenURL(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input",
			Message: "URL is required and must be a valid string",
		})
		return
	}

	req.URL = strings.TrimSpace(req.URL)

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input",
			Message: "URL cannot be empty",
		})
		return
	}

	shortCode, err := services.CreateShortURL(req.URL)
	if err != nil {
		switch err {
		case services.ErrInvalidURL:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid URL",
				Message: "Please provide a valid URL format",
			})
		case services.ErrCodeGenerationFail:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Service temporarily unavailable",
				Message: "Failed to generate unique code, please try again",
			})
		case services.ErrDatabaseError:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Failed to create short URL",
			})
		}
		return
	}

	baseURL := getBaseURL()
	response := ShortenResponse{
		ShortURL:    baseURL + "/" + shortCode,
		OriginalURL: req.URL,
		Code:        shortCode,
	}

	c.JSON(http.StatusCreated, response)
}

func Redirect(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))

	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Short code is required",
		})
		return
	}

	originalURL, err := services.GetOriginalURL(code)
	if err != nil {
		switch err {
		case services.ErrURLNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Short URL not found or has expired",
			})
		case services.ErrDatabaseError:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Failed to retrieve URL",
			})
		}
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}


func ShortenWithCustomCode(c *gin.Context) {
	var req CustomShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input",
			Message: "Both URL and custom_code are required",
		})
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	req.CustomCode = strings.TrimSpace(req.CustomCode)

	if req.URL == "" || req.CustomCode == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input",
			Message: "URL and custom code cannot be empty",
		})
		return
	}

	shortCode, err := services.CreateShortURLWithCustomCode(req.URL, req.CustomCode)
	if err != nil {
		switch err {
		case services.ErrInvalidURL:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid URL",
				Message: "Please provide a valid URL format",
			})
		case services.ErrInvalidCustomCode:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid custom code",
				Message: err.Error(),
			})
		case services.ErrCodeAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Code already exists",
				Message: "The custom code '" + req.CustomCode + "' is already in use. Please choose a different code.",
			})
		case services.ErrDatabaseError:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Database operation failed",
			})
		default:
		
			if strings.Contains(err.Error(), "URL already exists with code") {
				c.JSON(http.StatusConflict, ErrorResponse{
					Error:   "URL already shortened",
					Message: err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal server error",
					Message: "Failed to create short URL",
				})
			}
		}
		return
	}

	baseURL := getBaseURL()
	response := ShortenResponse{
		ShortURL:    baseURL + "/" + shortCode,
		OriginalURL: req.URL,
		Code:        shortCode,
	}

	c.JSON(http.StatusCreated, response)
}

func GetStats(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))

	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Short code is required",
		})
		return
	}

	stats, err := services.GetURLStats(code)
	if err != nil {
		switch err {
		case services.ErrURLNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Short URL not found",
			})
		case services.ErrDatabaseError:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Failed to retrieve statistics",
			})
		}
		return
	}

	c.JSON(http.StatusOK, stats)
}

func DeleteURL(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))

	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: "Short code is required",
		})
		return
	}

	err := services.DeleteURL(code)
	if err != nil {
		switch err {
		case services.ErrURLNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Short URL not found",
			})
		case services.ErrDatabaseError:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Database operation failed",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Message: "Failed to delete URL",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Short URL deleted successfully",
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "url-shortener",
	})
}
