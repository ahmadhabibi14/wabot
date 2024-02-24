package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ahmadhabibi14/wabot/handlers"
	"github.com/ahmadhabibi14/wabot/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type Bot struct {
	Name   string
	Log    waLog.Logger
	Device *store.Device
}

func NewBot(botName string, botLog waLog.Logger, botDevice *store.Device) *Bot {
	return &Bot{
		Name:   botName,
		Log:    botLog,
		Device: botDevice,
	}
}

func (b *Bot) Start() {
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()
	store.DeviceProps.Os = proto.String(b.Name)

	client := whatsmeow.NewClient(b.Device, b.Log)
	eventHandler := event(client)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		utils.PanicIfError(err)
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("Please scan to your whatsapp")
			} else {
				fmt.Println("Success Login !!!")
			}
		}
	} else {
		err := client.Connect()
		utils.PanicIfError(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
}

func event(client *whatsmeow.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message.GetConversation()
			// img := v.Message.GetImageMessage()
			ctx := context.Background()
			if !v.Info.IsGroup {
				textToText(msg, ctx, client, v, v.Info.Sender)
				// if img.GetCaption() == `/sticker` {
				// 	data, _ := client.Download(img)
				// 	stc := CreateStickerIMG(client, v, data)
				// 	_, err := client.SendMessage(context.Background(), v.Info.Sender, stc)
				// 	if err != nil {
				// 		return
				// 	}
				// }
			}
		}
	}
}

var CMD_TextToText = map[string]func(ctx context.Context, in string) string{
	`/help`:   handlers.Help,
	`/chtgpt`: handlers.ChatGPT,
	`/gemini`: handlers.GeminiAI,
}

func textToText(msg string, ctx context.Context, client *whatsmeow.Client, v *events.Message, to types.JID) {
	for key, value := range CMD_TextToText {
		if strings.Contains(msg, key) {
			res := value(ctx, msg)
			messageText(client, ctx, v.Info.Sender, res)
		}
	}
}

func messageText(client *whatsmeow.Client, ctx context.Context, to types.JID, msg string) {
	_, err := client.SendMessage(ctx, to, &waProto.Message{
		Conversation: proto.String(msg),
	})
	if err != nil {
		log.Println(err)
	}
}
