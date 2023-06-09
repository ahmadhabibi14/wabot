package commands

import (
	"context"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func ChatGPT(msg string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 5,
		Prompt:    msg,
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return "Completion error"
	}
	return strings.TrimSpace(resp.Choices[0].Text)
}
