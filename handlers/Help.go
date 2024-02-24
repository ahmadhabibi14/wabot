package handlers

import (
	"os"
	"strings"
)

func Help() string {
	botName := os.Getenv("BOT_NAME")
	msg := `*` + botName + `* - WhatsApp Bot

Hey, can I help you ?
This account is a bot 🤖`

	return strings.TrimSpace(msg)
}
