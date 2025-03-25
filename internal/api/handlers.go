package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rgehrsitz/me/internal/db"
	"github.com/rgehrsitz/me/internal/models"
)

// CreateContent handles the creation of new content
func (s *Server) CreateContent(c *gin.Context) {
	var content db.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.db.CreateContent(&content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create content: %v", err)})
		return
	}

	// Generate embedding if text is present
	if content.Body != "" {
		go func() {
			ctx := context.Background()
			embedding, err := s.embeddingService.GenerateEmbedding(ctx, content.Body)
			if err != nil {
				fmt.Printf("Failed to generate embedding: %v\n", err)
				return
			}

			embeddingBytes, err := s.embeddingService.SerializeEmbedding(embedding)
			if err != nil {
				fmt.Printf("Failed to serialize embedding: %v\n", err)
				return
			}

			_, err = s.db.StoreEmbedding(id, embeddingBytes, "openai-ada-002", len(embedding))
			if err != nil {
				fmt.Printf("Failed to store embedding: %v\n", err)
			}
		}()
	}

	// Get the created content with ID
	createdContent, err := s.db.GetContent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Content created but failed to retrieve: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, createdContent)
}

// ListContent handles listing all content
func (s *Server) ListContent(c *gin.Context) {
	contentType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	contents, err := s.db.ListContent(contentType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list content: %v", err)})
		return
	}

	c.JSON(http.StatusOK, contents)
}

// GetContent handles getting a single content item
func (s *Server) GetContent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	content, err := s.db.GetContent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Content not found: %v", err)})
		return
	}

	c.JSON(http.StatusOK, content)
}

// UpdateContent handles updating a content item
func (s *Server) UpdateContent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var content db.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content.ID = id
	err = s.db.UpdateContent(&content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update content: %v", err)})
		return
	}

	// Re-generate embedding if text changed
	if content.Body != "" {
		go func() {
			ctx := context.Background()
			embedding, err := s.embeddingService.GenerateEmbedding(ctx, content.Body)
			if err != nil {
				fmt.Printf("Failed to generate embedding: %v\n", err)
				return
			}

			embeddingBytes, err := s.embeddingService.SerializeEmbedding(embedding)
			if err != nil {
				fmt.Printf("Failed to serialize embedding: %v\n", err)
				return
			}

			_, err = s.db.StoreEmbedding(id, embeddingBytes, "openai-ada-002", len(embedding))
			if err != nil {
				fmt.Printf("Failed to store embedding: %v\n", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content updated successfully"})
}

// DeleteContent handles deleting a content item
func (s *Server) DeleteContent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = s.db.DeleteContent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete content: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content deleted successfully"})
}

// GenerateEmbedding handles generating an embedding for a content item
func (s *Server) GenerateEmbedding(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	content, err := s.db.GetContent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Content not found: %v", err)})
		return
	}

	if content.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content has no text to embed"})
		return
	}

	ctx := c.Request.Context()
	embedding, err := s.embeddingService.GenerateEmbedding(ctx, content.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate embedding: %v", err)})
		return
	}

	embeddingBytes, err := s.embeddingService.SerializeEmbedding(embedding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to serialize embedding: %v", err)})
		return
	}

	embeddingID, err := s.db.StoreEmbedding(id, embeddingBytes, "openai-ada-002", len(embedding))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store embedding: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Embedding generated successfully",
		"embedding_id": embeddingID,
		"dimensions":   len(embedding),
	})
}

// Search handles searching for content
func (s *Server) Search(c *gin.Context) {
	var query models.SearchQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	results, err := s.searchService.Search(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Search failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, results)
}

// SummarizeContent handles summarizing a content item
func (s *Server) SummarizeContent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	content, err := s.db.GetContent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Content not found: %v", err)})
		return
	}

	if content.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content has no text to summarize"})
		return
	}

	ctx := c.Request.Context()
	summary, err := s.summarizeService.Summarize(ctx, content.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate summary: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content_id": id,
		"summary":    summary,
	})
}

// ListTags handles listing all tags
func (s *Server) ListTags(c *gin.Context) {
	rows, err := s.db.Query("SELECT id, name FROM tags ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list tags: %v", err)})
		return
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan tag: %v", err)})
			return
		}
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, tags)
}

// CreateTag handles creating a new tag
func (s *Server) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := s.db.Exec("INSERT INTO tags (name) VALUES (?)", tag.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create tag: %v", err)})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get tag ID: %v", err)})
		return
	}

	tag.ID = id
	c.JSON(http.StatusCreated, tag)
}
