package services

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"url-shortener/database"
	"url-shortener/models"

	"github.com/lib/pq"
)

const (
	MaxRetries       = 5
	CodeLength       = 6
	MaxCustomCodeLen = 50
	MinCustomCodeLen = 3
)

var (
	ErrURLNotFound        = errors.New("URL not found")
	ErrInvalidURL         = errors.New("invalid URL format")
	ErrCodeAlreadyExists  = errors.New("code already exists")
	ErrInvalidCustomCode  = errors.New("invalid custom code format")
	ErrCodeGenerationFail = errors.New("failed to generate unique code")
	ErrDatabaseError      = errors.New("database error")
)

func generateSecureCode(n int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrInvalidURL
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

func validateCustomCode(code string) error {
	if len(code) < MinCustomCodeLen || len(code) > MaxCustomCodeLen {
		return fmt.Errorf("%w: length must be between %d and %d characters",
			ErrInvalidCustomCode, MinCustomCodeLen, MaxCustomCodeLen)
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, code)
	if err != nil {
		return fmt.Errorf("%w: regex error", ErrInvalidCustomCode)
	}

	if !matched {
		return fmt.Errorf("%w: only alphanumeric characters, hyphens, and underscores allowed", ErrInvalidCustomCode)
	}

	reservedWords := []string{"api", "admin", "www", "app", "mail", "ftp", "localhost"}
	for _, word := range reservedWords {
		if strings.EqualFold(code, word) {
			return fmt.Errorf("%w: '%s' is a reserved word", ErrInvalidCustomCode, code)
		}
	}

	return nil
}

func normalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

func codeExists(code string) (bool, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = $1", code).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return count > 0, nil
}

func GetOriginalURL(code string) (string, error) {
	if code == "" {
		return "", ErrURLNotFound
	}

	var url models.URL
	query := "SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1"

	row := database.DB.QueryRow(query, code)
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrURLNotFound
		}
		return "", fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	go func() {
		_, err := database.DB.Exec("UPDATE urls SET last_accessed = NOW(), click_count = click_count + 1 WHERE id = $1", url.ID)
		if err != nil {
			fmt.Printf("Warning: Failed to update click count for URL %d: %v\n", url.ID, err)
		}
	}()

	return url.OriginalURL, nil
}

func CreateShortURL(originalURL string) (string, error) {

	if err := validateURL(originalURL); err != nil {
		return "", err
	}


	normalizedURL := normalizeURL(originalURL)

	var existingCode string
	err := database.DB.QueryRow("SELECT short_code FROM urls WHERE original_url = $1", normalizedURL).Scan(&existingCode)
	if err == nil {
		return existingCode, nil
	} else if err != sql.ErrNoRows {
		return "", fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	var code string
	for i := 0; i < MaxRetries; i++ {
		generatedCode, err := generateSecureCode(CodeLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate code: %w", err)
		}

		exists, err := codeExists(generatedCode)
		if err != nil {
			return "", err
		}

		if !exists {
			code = generatedCode
			break
		}
	}

	if code == "" {
		return "", ErrCodeGenerationFail
	}


	query := `INSERT INTO urls (original_url, short_code, created_at) VALUES ($1, $2, NOW())`
	_, err = database.DB.Exec(query, normalizedURL, code)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { 
				return "", ErrCodeAlreadyExists
			}
		}
		return "", fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return code, nil
}

func CreateShortURLWithCustomCode(originalURL, customCode string) (string, error) {
	if err := validateURL(originalURL); err != nil {
		return "", err
	}

	if err := validateCustomCode(customCode); err != nil {
		return "", err
	}

	normalizedURL := normalizeURL(originalURL)

	exists, err := codeExists(customCode)
	if err != nil {
		return "", err
	}

	if exists {
		return "", ErrCodeAlreadyExists
	}

	var existingCode string
	err = database.DB.QueryRow("SELECT short_code FROM urls WHERE original_url = $1", normalizedURL).Scan(&existingCode)
	if err == nil {
		return "", fmt.Errorf("URL already exists with code: %s", existingCode)
	} else if err != sql.ErrNoRows {
		return "", fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	query := `INSERT INTO urls (original_url, short_code, created_at) VALUES ($1, $2, NOW())`
	_, err = database.DB.Exec(query, normalizedURL, customCode)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { 
				return "", ErrCodeAlreadyExists
			}
		}
		return "", fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return customCode, nil
}

func GetURLStats(code string) (*models.URLStats, error) {
	if code == "" {
		return nil, ErrURLNotFound
	}

	var stats models.URLStats
	query := `SELECT short_code, original_url, created_at, 
			  COALESCE(click_count, 0) as click_count,
			  last_accessed
			  FROM urls WHERE short_code = $1`

	row := database.DB.QueryRow(query, code)
	err := row.Scan(&stats.ShortCode, &stats.OriginalURL, &stats.CreatedAt,
		&stats.ClickCount, &stats.LastAccessed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrURLNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return &stats, nil
}

func DeleteURL(code string) error {
	if code == "" {
		return ErrURLNotFound
	}

	result, err := database.DB.Exec("DELETE FROM urls WHERE short_code = $1", code)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrURLNotFound
	}

	return nil
}
