package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeminiAI(ctx context.Context, in string) string {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMIN_API_KEY")))
	if err != nil {
		return ``
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(in))
	if err != nil {
		return ``
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
}
