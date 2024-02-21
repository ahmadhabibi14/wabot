package services

import (
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func event(client *whatsmeow.Client) func(evt interface{}) {
	textMsg := commands.Messages
	msgReceive := models.Message{}

	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message.GetConversation()
			if !v.Info.IsGroup {
				if msg != "" {
					for key, value := range textMsg {
						if strings.Contains(msg, key) {
							msgReceive.MsgReceive = msg
							textMsg(client, v, v.Info.Sender, value)
						}
					}
				}
			}
		}
	}
}
