package main

import (
	"context"
	"fmt"
	"llm-primitives/model"
)

func main() {
	ctx := context.Background()
	res, err := model.BinaryClassify(ctx, "Determine if the sentiment of the text is positive", "I love this product")
	if err != nil {
		panic(fmt.Errorf("cannot classify sentiment, %w", err))
	}
	fmt.Println(res)
}
