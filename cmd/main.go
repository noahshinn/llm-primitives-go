package main

import (
	"context"
	"fmt"
	"llm-primitives/model"
)

func main() {
	ctx := context.Background()
	type address struct {
		StreetName   string `json:"street_name"`
		StreetNumber int    `json:"street_number"`
	}
	addr, err := model.Parse[address](ctx, "My street is 123 main st")
	if err != nil {
		panic(fmt.Errorf("cannot parse address, %w", err))
	}
	fmt.Println("Street Name:", addr.Obj.StreetName)
	fmt.Println("Street Number:", addr.Obj.StreetNumber)
}
