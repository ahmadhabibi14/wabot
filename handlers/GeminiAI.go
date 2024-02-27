package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeminiAI(ctx context.Context, in string) string {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return fmt.Sprintf("cannot initialize Gemini AI: %s", err.Error())
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(in))
	if err != nil {
		return fmt.Sprintf("error generate text: %s", err.Error())
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
}
