package main

import (
	"context"
	"fmt"
	"llm-primitives/model"
)

func main() {
	ctx := context.Background()
	res, err := model.Classify(ctx, "Determine the sentiment of the following text", "I love this product", []string{"positive", "negative", "neutral"})
	if err != nil {
		panic(fmt.Errorf("cannot classify sentiment, %w", err))
	}
	fmt.Println(res)
}
