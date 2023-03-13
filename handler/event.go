package handler

import (
	"log"
	"strings"

	"github.com/ahmadhabibi14/wabot/commands"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(client *whatsmeow.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if !v.Info.IsGroup {
				if v.Message.GetConversation() != "" {
					msg := v.Message.GetConversation()
					if strings.Contains(msg, "/ai") { // ChatGPT Command
						MessageHandler(client, v, v.Info.Sender, commands.ChatGPT(msg))
					} else if v.Message.GetConversation() == "/help" { // Help Command
						MessageHandler(client, v, v.Info.Sender, commands.Help())
					}
				}
			} else if v.Info.IsGroup == true {
				msg := v.Message.GetConversation()
				groupInfo, err := client.GetGroupInfo(v.Info.Chat)
				if err != nil {
					log.Fatalf("Error get group info %v\n", err)
				}
				if strings.Contains(msg, "/ai") { // ChatGPT Command
					MessageHandler(client, v, groupInfo.JID, commands.ChatGPT(msg))
				}
			}
		}
	}
}
