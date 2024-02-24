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
			img := v.Message.GetImageMessage()
			ctx := context.Background()

			textToText(msg, ctx, client, v, v.Info.Sender)
			generateSticker(img, ctx, client, v, v.Info.Chat)
			sendImgBack(img, ctx, client, v, v.Info.Chat)
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

func generateSticker(img *waProto.ImageMessage, ctx context.Context, client *whatsmeow.Client, v *events.Message, to types.JID) {
	if img.GetCaption() == `/sticker` {
		data, err := client.Download(img)
		if err != nil {
			log.Println("ERROR Download file")
		}
		rawPath := fmt.Sprintf("tmp/%s.jpg", v.Info.ID)
		err = os.WriteFile(rawPath, data, 0600)
		if err != nil {
			log.Println("ERROR cannot write file")
		}

		_, err = client.SendMessage(ctx, to, &waProto.Message{
			StickerMessage: &waProto.StickerMessage{
				Url:           proto.String(img.GetUrl()),
				DirectPath:    proto.String(img.GetDirectPath()),
				MediaKey:      img.GetMediaKey(),
				Mimetype:      proto.String(img.GetMimetype()),
				FileEncSha256: img.GetFileEncSha256(),
				FileSha256:    img.GetFileSha256(),
				FileLength:    proto.Uint64(img.GetFileLength()),
				ContextInfo:   img.GetContextInfo(),
				PngThumbnail:  img.GetThumbnailSha256(),
			},
		})
		if err != nil {
			log.Println("ERROR send message:", err)
		}
	}
}

func sendImgBack(img *waProto.ImageMessage, ctx context.Context, client *whatsmeow.Client, v *events.Message, to types.JID) {
	// if img.GetCaption() == `/sendback` {
	data, err := client.Download(img)
	if err != nil {
		log.Println("ERROR Download file")
	}

	rawPath := fmt.Sprintf("tmp/%s.jpg", v.Info.ID)
	err = os.WriteFile(rawPath, data, 0600)
	if err != nil {
		log.Println("ERROR cannot write file")
	}

	_, err = client.SendMessage(ctx, to, &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Url:           proto.String(img.GetUrl()),
			DirectPath:    proto.String(img.GetDirectPath()),
			MediaKey:      img.GetMediaKey(),
			Mimetype:      proto.String(img.GetMimetype()),
			FileEncSha256: img.GetFileEncSha256(),
			FileSha256:    img.GetFileSha256(),
			FileLength:    proto.Uint64(img.GetFileLength()),
			ContextInfo:   img.GetContextInfo(),
		},
	})
	if err != nil {
		log.Println("ERROR send message:", err)
	}
	// }
}
