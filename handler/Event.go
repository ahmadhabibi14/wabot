package handler

import (
	"strings"

	"github.com/ahmadhabibi14/wabot/commands"
	"github.com/ahmadhabibi14/wabot/models"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Event(client *whatsmeow.Client) func(evt interface{}) {
	textMsg := commands.Messages
	msgReceive := models.Message{}

	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message.GetConversation()
			// img := v.Message.ImageMessage
			if !v.Info.IsGroup {
				if msg != "" {
					for key, value := range textMsg {
						// Check if there is a message contains like key from TextMsg
						if strings.Contains(msg, key) {
							msgReceive.MsgReceive = msg
							TextMsg(client, v, v.Info.Sender, value)
						}
					}
				}
			}
		}
	}
}
