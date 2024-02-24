package handlers

import (
	"context"
	"os"

	"github.com/0x9ef/openai-go"
)

func ChatGPT(ctx context.Context, in string) string {
	e := openai.New(os.Getenv("OPENAI_API_KEY"))
	resp, err := e.Completion(ctx, &openai.CompletionOptions{
		Model:  openai.ModelGPT3TextDavinci003,
		Prompt: []string{in},
	})
	if err != nil {
		return ``
	}

	return resp.Choices[0].Text
}
