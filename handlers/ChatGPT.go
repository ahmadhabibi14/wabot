package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/0x9ef/openai-go"
)

func ChatGPT(ctx context.Context, in string) string {
	e := openai.New(os.Getenv("OPENAI_API_KEY"))

	var prompt = in
	resp, err := e.Completion(context.Background(), &openai.CompletionOptions{
		Model:  openai.ModelGPT3TextDavinci003,
		Prompt: []string{prompt},
	})

	if err != nil {
		errorMSg := fmt.Sprintf("Error: %s", err.Error())
		return errorMSg
	}
	if b, err := json.MarshalIndent(resp, "", "  "); err != nil {
		log.Println(err)
	} else {
		log.Println(string(b))
	}

	return resp.Choices[0].Text
}
