package services

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// SummarizeService handles summarizing content
type SummarizeService struct {
	openAIClient *openai.Client
}

// NewSummarizeService creates a new summarize service
func NewSummarizeService() (*SummarizeService, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)
	return &SummarizeService{
		openAIClient: client,
	}, nil
}

// Summarize summarizes the given text
func (s *SummarizeService) Summarize(ctx context.Context, text string) (string, error) {
	resp, err := s.openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that summarizes text concisely.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Please summarize the following text in a few sentences:\n\n%s", text),
			},
		},
		MaxTokens: 150,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create summary: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no summary data returned")
	}

	return resp.Choices[0].Message.Content, nil
}
