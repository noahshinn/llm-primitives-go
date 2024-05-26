package model

import (
	"context"
	"errors"
	"fmt"
	"llm-primitives/llm"
	"os"
)

func GenerateText(ctx context.Context, instruction string, text string) (string, error) {
	apiKey := os.Getenv(openaiApiKeyName)
	if apiKey == "" {
		return "", errors.New("please set OPENAI_API_KEY as an environment variable")
	}
	models := llm.AllModels(apiKey)
	defaultChatModel := models.DefaultChatModel
	// TODO: refine messages
	res, err := defaultChatModel.Message(ctx, []*llm.Message{
		{
			Role:    llm.MessageRoleSystem,
			Content: instruction,
		},
		{
			Role:    llm.MessageRoleUser,
			Content: text,
		},
	}, &llm.MessageOptions{
		Temperature: 0.0,
		ForceJson:   false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate a message from the llm: %w", err)
	}
	return res.Content, nil
}
