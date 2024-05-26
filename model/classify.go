package model

import (
	"context"
	"errors"
	"fmt"
	"llm-primitives/llm"
	"os"
	"strings"
)

func Classify(ctx context.Context, instruction string, text string, choices []string) (int, error) {
	apiKey := os.Getenv(openaiApiKeyName)
	if apiKey == "" {
		return 0, errors.New("please set OPENAI_API_KEY as an environment variable")
	}
	models := llm.AllModels(apiKey)
	defaultChatModel := models.DefaultChatModel
	choicesDisplay, decodeMap := displayChoices(choices)
	inputText := fmt.Sprintf("Instruction:\n%s\n\nText:\n%s\n\nChoices:\n%s\n\nValid JSON:", instruction, text, choicesDisplay)
	messages := []*llm.Message{
		{
			Role:    llm.MessageRoleSystem,
			Content: "Classify the following text with the provided instruction and choices. To classify, provide the key of the choice:\n{\"classification\": string}\n\nFor example, if the correct choice is 'Z. description of choice Z', then provide 'Z' as the classification as valid JSON:\n{\"classification\": \"Z\"}",
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
	if score, ok := res.Obj["classification"].(string); !ok {
		return 0, errors.New("failed to parse classification from response")
	} else if index, ok := decodeMap[score]; !ok {
		return 0, errors.New("invalid classification")
	} else {
		return index, nil
	}
}

func BinaryClassify(ctx context.Context, instruction string, text string) (bool, error) {
	res, err := Classify(ctx, instruction, text, []string{"true", "false"})
	if err != nil {
		return false, err
	}
	return res == 0, nil
}

func displayChoices(choices []string) (string, map[string]int) {
	choicesDisplays := []string{}
	decodeMap := map[string]int{}
	for i, choice := range choices {
		label := indexToAlpha(i)
		choicesDisplays = append(choicesDisplays, fmt.Sprintf("%s. %s", label, choice))
		decodeMap[label] = i
	}
	return strings.Join(choicesDisplays, "\n"), decodeMap
}

func indexToAlpha(index int) string {
	alpha := ""
	for index >= 0 {
		alpha = fmt.Sprint('A'+(index%26)) + alpha
		index = index/26 - 1
	}
	return alpha
}
