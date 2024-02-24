package handlers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeminiAI(ctx context.Context, in string) string {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Println(`GeminiAI`, err)
		return ``
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(in))
	if err != nil {
		log.Println(`GeminiAI`, err)
		return ``
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
}
