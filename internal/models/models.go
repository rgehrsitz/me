package models

import (
	"time"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeNote     ContentType = "note"
	ContentTypeSnippet  ContentType = "snippet"
	ContentTypeBookmark ContentType = "bookmark"
	ContentTypeDocument ContentType = "document"
)

// Content represents a piece of content in the PKB
type Content struct {
	ID        int64       `json:"id"`
	Type      ContentType `json:"type"`
	Title     string      `json:"title"`
	Body      string      `json:"body"`
	SourceURL string      `json:"source_url,omitempty"`
	FilePath  string      `json:"file_path,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Tags      []string    `json:"tags,omitempty"`
	Embedding *Embedding  `json:"embedding,omitempty"`
}

// Embedding represents a vector embedding for a piece of content
type Embedding struct {
	ID         int64     `json:"id"`
	ContentID  int64     `json:"content_id"`
	Vector     []float32 `json:"vector"`
	Model      string    `json:"model"`
	Dimensions int       `json:"dimensions"`
}

// Tag represents a tag for organizing content
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SearchQuery represents a search query
type SearchQuery struct {
	Query     string   `json:"query"`
	Type      string   `json:"type,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
	Semantic  bool     `json:"semantic"`
}

// SearchResult represents a search result
type SearchResult struct {
	Content  Content `json:"content"`
	Score    float64 `json:"score,omitempty"`
	Snippet  string  `json:"snippet,omitempty"`
}
