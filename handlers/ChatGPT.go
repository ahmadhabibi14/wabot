package handlers

import (
	"context"
	"os"
	"strings"

	"github.com/0x9ef/openai-go"
)

func ChatGPT(ctx context.Context, in string) string {
	in = strings.Replace(in, "/chtgpt", "", 1)

	e := openai.New(os.Getenv("OPENAI_API_KEY"))
	resp, err := e.Completion(ctx, &openai.CompletionOptions{
		Model:  openai.ModelGPT3TextDavinci003,
		Prompt: []string{in},
	})
	if err != nil {
		return err.Error()
	}

	return resp.Choices[0].Text
}
