package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"llm-primitives/llm"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/mitchellh/mapstructure"
)

func Parse[T any](ctx context.Context, text string) (ObjectResult[T], error) {
	apiKey := os.Getenv(openaiApiKeyName)
	if apiKey == "" {
		return ObjectResult[T]{}, errors.New("please set OPENAI_API_KEY as an environment variable")
	}
	models := llm.AllModels(apiKey)
	defaultChatModel := models.DefaultChatModel
	jsonSchema, err := structToJsonSchema[T]()
	if err != nil {
		return ObjectResult[T]{}, fmt.Errorf("failed to generate json schema: %w", err)
	}
	jsonSchemaStr := jsonSchemaToString(jsonSchema)
	res, err := defaultChatModel.Message(ctx, []*llm.Message{
		{
			Role:    llm.MessageRoleSystem,
			Content: "Parse the following text into a struct with the provided json schema.",
		},
		{
			Role:    llm.MessageRoleUser,
			Content: fmt.Sprintf("Text: %s\nSchema: %s", text, jsonSchemaStr),
		},
	}, &llm.MessageOptions{
		Temperature: 0.0,
		ForceJson:   true,
	})
	if err != nil {
		return ObjectResult[T]{}, fmt.Errorf("failed to generate a message from the llm: %w", err)
	}
	obj := res.Obj
	return jsonResponseToObject[T](obj)
}

func structToJsonSchema[T any]() (map[string]any, error) {
	var t T
	schema := jsonschema.Reflect(&t)
	b, err := schema.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}
	m := map[string]any{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	return m, nil
}

func jsonResponseToObject[T any](response map[string]any) (ObjectResult[T], error) {
	standardizedResponse := map[string]any{}
	for k, v := range response {
		if !strings.HasPrefix(k, "$") {
			standardizedResponse[k] = v
			ccK := toCamelCase(k)
			standardizedResponse[ccK] = v
			pcK := toPascalCase(k)
			standardizedResponse[pcK] = v
		}
	}
	var obj T
	err := mapstructure.Decode(standardizedResponse, &obj)
	if err != nil {
		return ObjectResult[T]{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return ObjectResult[T]{Obj: obj}, nil
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i > 0 {
			parts[i] = fmt.Sprintf("%s%s", strings.ToUpper(part[:1]), strings.ToLower(part[1:]))
		} else {
			parts[i] = strings.ToLower(part)
		}
	}
	return strings.Join(parts, "")
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = fmt.Sprintf("%s%s", strings.ToUpper(part[:1]), strings.ToLower(part[1:]))
	}
	return strings.Join(parts, "")
}

func jsonSchemaToString(schema map[string]any) string {
	return fmt.Sprintf("%v", schema)
}
