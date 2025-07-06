package main

import (
	"url-shortener/database"
	"url-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)
	r.GET("/:code", handlers.Redirect)
	r.GET("/stats/:code", handlers.GetStats)

	r.POST("/shorten", handlers.ShortenURL)
	r.POST("/shorten/custom", handlers.ShortenWithCustomCode)

	r.DELETE("/delete/:code", handlers.DeleteURL)

	r.Run(":8080")
}
