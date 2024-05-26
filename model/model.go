package model

const openaiApiKeyName = "OPENAI_API_KEY"

type ObjectResult[T any] struct {
	Obj T
}
