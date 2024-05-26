package llm

func ApproxNumTokensFast(text string) int {
	return len(text) / 4
}
