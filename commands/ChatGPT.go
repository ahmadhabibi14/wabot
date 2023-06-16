package commands

import (
	"fmt"

	"github.com/ahmadhabibi14/wabot/models"
)

func ChatGPT() string {
	var msg = models.Message{}
	var chatgpt string = "chatgpt"
	var resp string = fmt.Sprintf("%s : %s", chatgpt, msg)
	return resp
}
