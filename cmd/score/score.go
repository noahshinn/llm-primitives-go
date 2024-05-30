package main

import (
	"context"
	"fmt"
	"llm-primitives/model"
)

func main() {
	ctx := context.Background()
	res, err := model.ScoreInt(ctx, "Rate the product review from (1) bad to (5) good", "I love this product", 1, 5)
	if err != nil {
		panic(fmt.Errorf("cannot score, %w", err))
	}
	fmt.Println(res)
}
