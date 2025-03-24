package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/yourusername/pkb/internal/models"
)

// EmbeddingService handles generating and storing embeddings
type EmbeddingService struct {
	openAIClient *openai.Client
}

// NewEmbeddingService creates a new embedding service
func NewEmbeddingService() (*EmbeddingService, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)
	return &EmbeddingService{
		openAIClient: client,
	}, nil
}

// GenerateEmbedding generates an embedding for the given text
func (s *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := s.openAIClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: []string{text},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	return resp.Data[0].Embedding, nil
}

// SerializeEmbedding serializes an embedding to a byte array
func (s *EmbeddingService) SerializeEmbedding(embedding []float32) ([]byte, error) {
	return json.Marshal(embedding)
}

// DeserializeEmbedding deserializes an embedding from a byte array
func (s *EmbeddingService) DeserializeEmbedding(data []byte) ([]float32, error) {
	var embedding []float32
	if err := json.Unmarshal(data, &embedding); err != nil {
		return nil, fmt.Errorf("failed to deserialize embedding: %w", err)
	}
	return embedding, nil
}
