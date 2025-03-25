package services

import (
	"context"

	"fmt"
	"math"
	"strings"

	"github.com/rgehrsitz/me/internal/db"
	"github.com/rgehrsitz/me/internal/models"
)

// SearchService handles searching for content
type SearchService struct {
	db              *db.DB
	embeddingService *EmbeddingService
}

// NewSearchService creates a new search service
func NewSearchService(db *db.DB, embeddingService *EmbeddingService) *SearchService {
	return &SearchService{
		db:              db,
		embeddingService: embeddingService,
	}
}

// Search searches for content based on the given query
func (s *SearchService) Search(ctx context.Context, query models.SearchQuery) ([]models.SearchResult, error) {
	if query.Semantic {
		return s.semanticSearch(ctx, query)
	}
	return s.keywordSearch(query)
}

// keywordSearch performs a keyword-based search
func (s *SearchService) keywordSearch(query models.SearchQuery) ([]models.SearchResult, error) {
	// Build the SQL query
	sqlQuery := `
		SELECT c.id, c.type, c.title, c.body, c.source_url, c.file_path, c.created_at, c.updated_at
		FROM content c
		WHERE (c.title LIKE ? OR c.body LIKE ?)`

	args := []interface{}{
		"%" + query.Query + "%",
		"%" + query.Query + "%",
	}

	// Add type filter if specified
	if query.Type != "" {
		sqlQuery += " AND c.type = ?"
		args = append(args, query.Type)
	}

	// Add tag filters if specified
	if len(query.Tags) > 0 {
		placeholders := strings.Repeat("?,", len(query.Tags)-1) + "?"
		sqlQuery += fmt.Sprintf(` AND c.id IN (
			SELECT ct.content_id
			FROM content_tags ct
			JOIN tags t ON ct.tag_id = t.id
			WHERE t.name IN (%s)
			GROUP BY ct.content_id
			HAVING COUNT(DISTINCT t.name) = ?
		)`, placeholders)

		for _, tag := range query.Tags {
			args = append(args, tag)
		}
		args = append(args, len(query.Tags))
	}

	// Add limit and offset
	if query.Limit <= 0 {
		query.Limit = 10
	}
	sqlQuery += " ORDER BY c.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, query.Limit, query.Offset)

	// Execute the query
	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	// Process the results
	results := []models.SearchResult{}
	for rows.Next() {
		var content models.Content
		var createdAt, updatedAt string

		err := rows.Scan(
			&content.ID,
			&content.Type,
			&content.Title,
			&content.Body,
			&content.SourceURL,
			&content.FilePath,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}

		// Get tags for this content
		content.Tags, err = s.getContentTags(content.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get content tags: %w", err)
		}

		// Create a search result
		result := models.SearchResult{
			Content: content,
			Snippet: extractSnippet(content.Body, query.Query),
		}

		results = append(results, result)
	}

	return results, nil
}

// semanticSearch performs a semantic search using embeddings
func (s *SearchService) semanticSearch(ctx context.Context, query models.SearchQuery) ([]models.SearchResult, error) {
	// Generate embedding for the query
	queryEmbedding, err := s.embeddingService.GenerateEmbedding(ctx, query.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Since SQLite doesn't have built-in vector similarity search,
	// we'll retrieve all embeddings and compute similarity in memory
	// In a production system, you'd want to use a vector database or SQLite extension
	// Build the SQL query to get all embeddings
	sqlQuery := `
		SELECT c.id, c.type, c.title, c.body, c.source_url, c.file_path, c.created_at, c.updated_at, e.embedding
		FROM content c
		JOIN embeddings e ON c.id = e.content_id`

	args := []interface{}{}

	// Add type filter if specified
	if query.Type != "" {
		sqlQuery += " WHERE c.type = ?"
		args = append(args, query.Type)
	}

	// Execute the query
	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute semantic search query: %w", err)
	}
	defer rows.Close()

	// Process the results and compute similarities
	results := []models.SearchResult{}
	for rows.Next() {
		var content models.Content
		var createdAt, updatedAt string
		var embeddingBytes []byte

		err := rows.Scan(
			&content.ID,
			&content.Type,
			&content.Title,
			&content.Body,
			&content.SourceURL,
			&content.FilePath,
			&createdAt,
			&updatedAt,
			&embeddingBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan semantic search result: %w", err)
		}

		// Deserialize the embedding
		contentEmbedding, err := s.embeddingService.DeserializeEmbedding(embeddingBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize embedding: %w", err)
		}

		// Compute cosine similarity
		similarity := cosineSimilarity(queryEmbedding, contentEmbedding)

		// Get tags for this content
		content.Tags, err = s.getContentTags(content.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get content tags: %w", err)
		}

		// Filter by tags if specified
		if len(query.Tags) > 0 && !containsAllTags(content.Tags, query.Tags) {
			continue
		}

		// Create a search result
		result := models.SearchResult{
			Content: content,
			Score:   similarity,
			Snippet: extractSnippet(content.Body, query.Query),
		}

		results = append(results, result)
	}

	// Sort results by similarity score (descending)
	sortResultsByScore(results)

	// Apply limit and offset
	if query.Limit <= 0 {
		query.Limit = 10
	}
	start := query.Offset
	end := query.Offset + query.Limit
	if start >= len(results) {
		return []models.SearchResult{}, nil
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

// getContentTags retrieves the tags for a content item
func (s *SearchService) getContentTags(contentID int64) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT t.name 
		FROM tags t 
		JOIN content_tags ct ON t.id = ct.tag_id 
		WHERE ct.content_id = ?`, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// cosineSimilarity computes the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	normA = math.Sqrt(normA)
	normB = math.Sqrt(normB)

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}

// extractSnippet extracts a snippet from the text containing the query
func extractSnippet(text, query string) string {
	if query == "" || text == "" {
		return ""
	}

	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	index := strings.Index(lowerText, lowerQuery)
	if index == -1 {
		// If the exact query isn't found, just return the beginning of the text
		if len(text) > 150 {
			return text[:150] + "..."
		}
		return text
	}

	// Find a good starting point for the snippet
	start := index - 75
	if start < 0 {
		start = 0
	}

	// Find a good ending point for the snippet
	end := index + len(query) + 75
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]

	// Add ellipsis if we're not at the beginning or end
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}

	return snippet
}

// containsAllTags checks if the content tags contain all the query tags
func containsAllTags(contentTags, queryTags []string) bool {
	tagMap := make(map[string]bool)
	for _, tag := range contentTags {
		tagMap[tag] = true
	}

	for _, tag := range queryTags {
		if !tagMap[tag] {
			return false
		}
	}

	return true
}

// sortResultsByScore sorts search results by score in descending order
func sortResultsByScore(results []models.SearchResult) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}
