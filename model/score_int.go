package model

import (
	"context"
	"errors"
	"fmt"
	"llm-primitives/llm"
	"os"
)

func ScoreInt(ctx context.Context, instruction string, text string, minBound int, maxBound int) (int, error) {
	apiKey := os.Getenv(openaiApiKeyName)
	if apiKey == "" {
		return 0, errors.New("please set OPENAI_API_KEY as an environment variable")
	}
	models := llm.AllModels(apiKey)
	defaultChatModel := models.DefaultChatModel
	inputText := fmt.Sprintf("Instruction:\n%s\n\nText:\n%s\n\nRange:\n[%d, %d]\n\nValid JSON:", instruction, text, minBound, maxBound)
	messages := []*llm.Message{
		{
			Role:    llm.MessageRoleSystem,
			Content: "Score the following text with the provided instruction and range as an integer value as valid JSON:\n{\"score\": int}",
		},
		{
			Role:    llm.MessageRoleUser,
			Content: inputText,
		},
	}
	res, err := defaultChatModel.Message(ctx, messages, &llm.MessageOptions{
		Temperature: 0.0,
		ForceJson:   true,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to generate a message from the llm: %w", err)
	}
	if score, ok := res.Obj["score"].(int); !ok {
		return 0, errors.New("failed to parse score from response")
	} else {
		return score, nil
	}
}
