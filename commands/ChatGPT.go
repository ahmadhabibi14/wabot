package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/0x9ef/openai-go"
	"github.com/ahmadhabibi14/wabot/models"
)

func ChatGPT() string {
	var msg = models.Message{}

	e := openai.New(os.Getenv("OPENAI_API_KEY"))

	msg.Mu.Lock()
	var prompt = fmt.Sprintf("%s", msg.MsgReceive)
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

	msg.Mu.Unlock()
	return resp.Choices[0].Text
}
