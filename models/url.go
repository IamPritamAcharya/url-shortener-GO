package models

import (
	"time"
)

type URL struct {
	ID           int        `json:"id" db:"id"`
	OriginalURL  string     `json:"original_url" db:"original_url"`
	ShortCode    string     `json:"short_code" db:"short_code"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	LastAccessed *time.Time `json:"last_accessed,omitempty" db:"last_accessed"`
	ClickCount   int        `json:"click_count" db:"click_count"`
}

type URLStats struct {
	ShortCode    string     `json:"short_code"`
	OriginalURL  string     `json:"original_url"`
	CreatedAt    time.Time  `json:"created_at"`
	ClickCount   int        `json:"click_count"`
	LastAccessed *time.Time `json:"last_accessed,omitempty"`
}
