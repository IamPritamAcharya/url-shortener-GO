package main

import (
	"url-shortener/database"
	"url-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.POST("/shorten", handlers.ShortenURL)
	r.GET("/:code", handlers.Redirect)

	r.Run(":8080")
}
