package llm

import (
	"context"
)

type ChatModelID string
type EmbeddingModelID string

const (
	ChatModelGPT35Turbo ChatModelID = "gpt-3.5-turbo"
	ChatModelGPT4Turbo  ChatModelID = "gpt-4-turbo"
	ChatModelGPT4O      ChatModelID = "gpt-4o"
	// TODO: update embedding model
	EmbeddingModelAda EmbeddingModelID = "text-embedding-ada-002"
)

type Models struct {
	DefaultChatModel            ChatModel
	DefaultLongContextChatModel ChatModel
	DefaultLightChatModel       ChatModel
	DefaultCheapChatModel       ChatModel
	DefaultEmbeddingModel       EmbeddingModel
	ChatModels                  map[ChatModelID]ChatModel
	EmbeddingModels             map[EmbeddingModelID]EmbeddingModel
}

func AllModels(api_key string) *Models {
	return &Models{
		DefaultChatModel:            NewOpenAIChatModel(ChatModelGPT4O, api_key),
		DefaultLongContextChatModel: NewOpenAIChatModel(ChatModelGPT4O, api_key),
		DefaultLightChatModel:       NewOpenAIChatModel(ChatModelGPT4O, api_key),
		DefaultCheapChatModel:       NewOpenAIChatModel(ChatModelGPT35Turbo, api_key),
		DefaultEmbeddingModel:       NewOpenAIEmbeddingModel(EmbeddingModelAda, api_key),
		ChatModels: map[ChatModelID]ChatModel{
			ChatModelGPT35Turbo: NewOpenAIChatModel(ChatModelGPT35Turbo, api_key),
			ChatModelGPT4Turbo:  NewOpenAIChatModel(ChatModelGPT4Turbo, api_key),
			ChatModelGPT4O:      NewOpenAIChatModel(ChatModelGPT4O, api_key),
		},
		EmbeddingModels: map[EmbeddingModelID]EmbeddingModel{
			EmbeddingModelAda: NewOpenAIEmbeddingModel(EmbeddingModelAda, api_key),
		},
	}
}

type Message struct {
	Role    MessageRole    `json:"role"`
	Content string         `json:"content"`
	Obj     map[string]any `json:"obj,omitempty"`
}

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type MessageOptions struct {
	Temperature   float64  `json:"temperature"`
	MaxTokens     int      `json:"max_tokens"`
	StopSequences []string `json:"stop_sequences"`

	// for OpenAI models
	ForceJson bool `json:"force_json"`
}

type ChatModel interface {
	Message(ctx context.Context, messages []*Message, options *MessageOptions) (*Message, error)
	ContextLength() int
}

type EmbeddingModel interface {
	Embedding(ctx context.Context, texts []string) ([][]float32, error)
}
