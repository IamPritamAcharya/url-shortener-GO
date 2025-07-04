package services

import (
	"errors"
	"math/rand"
	"time"
	"url-shortener/database"
	"url-shortener/models"
)

func generateCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	code := make([]rune, n)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func CreateShortURL(original string) (string, error) {
	code := generateCode(6)
	_, err := database.DB.Exec("INSERT INTO urls (original_url, short_code) VALUES ($1, $2)", original, code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func GetOriginalURL(code string) (string, error) {
	var url models.URL
	row := database.DB.QueryRow("SELECT id, original_url, short_code FROM urls WHERE short_code=$1", code)
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode)
	if err != nil {
		return "", errors.New("URL not found")
	}
	return url.OriginalURL, nil
}
