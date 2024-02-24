package services

import (
	"context"
	"log"
	"net/http"
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
				log.Println("Please scan to your whatsapp")
			} else {
				log.Println("Success Login !!!")
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

			generalMsg(msg, ctx, client, v, v.Info.Sender)
			commandToText(msg, ctx, client, v, v.Info.Sender)
			imageToImage(ctx, client, v, v.Info.Sender, img)
		}
	}
}

var CMD_General = map[string]func() string{
	`/menu`: handlers.Menu,
	`/help`: handlers.Help,
}

func generalMsg(
	msg string,
	ctx context.Context,
	client *whatsmeow.Client,
	v *events.Message,
	to types.JID,
) {
	for key, value := range CMD_General {
		if msg == key {
			res := value()
			messageText(client, ctx, v.Info.Sender, res)
		}
	}
}

var CMD_CommandToText = map[string]func(
	ctx context.Context,
	in string,
) string{
	`/gpt`:    handlers.ChatGPT,
	`/gemini`: handlers.GeminiAI,
}

func commandToText(
	msg string,
	ctx context.Context,
	client *whatsmeow.Client,
	v *events.Message,
	to types.JID,
) {
	for key, value := range CMD_CommandToText {
		if strings.HasPrefix(msg, key) {
			in := strings.TrimPrefix(msg, key)
			res := value(ctx, in)
			if res != `` {
				messageText(client, ctx, v.Info.Sender, res)
			}
		}
	}
}

var CMD_ImageToImage = map[string]func(
	client *whatsmeow.Client,
	img *waProto.ImageMessage,
	id string,
) ([]byte, string, error){
	`/buriq`: handlers.Buriq,
}

func imageToImage(
	ctx context.Context,
	client *whatsmeow.Client,
	v *events.Message,
	to types.JID,
	img *waProto.ImageMessage,
) {
	imgCaption := img.GetCaption()
	for key, value := range CMD_ImageToImage {
		if strings.Contains(imgCaption, key) {
			data, caption, err := value(client, img, v.Info.ID)
			sendMessageIfError(err, client, ctx, to)
			if err == nil {
				messageImage(ctx, client, v, to, data, caption)
			}
		}
	}
}

func messageText(
	client *whatsmeow.Client,
	ctx context.Context,
	to types.JID,
	msg string,
) {
	_, err := client.SendMessage(ctx, to, &waProto.Message{
		Conversation: proto.String(msg),
	})
	if err != nil {
		log.Println(err)
	}
}

func messageImage(
	ctx context.Context,
	client *whatsmeow.Client,
	v *events.Message,
	to types.JID,
	data []byte,
	caption string,
) {
	uploaded, err := client.Upload(ctx, data, whatsmeow.MediaImage)
	sendMessageIfError(err, client, ctx, to)

	_, err = client.SendMessage(ctx, to, &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &v.Info.ID,
				Participant:   proto.String(v.Info.Sender.String()),
				QuotedMessage: v.Message,
			},
			Caption: proto.String(caption),
		},
	})
	if err != nil {
		log.Println("ERROR send message:", err)
	}
}

func sendMessageIfError(
	err error,
	client *whatsmeow.Client,
	ctx context.Context,
	to types.JID,
) {
	if err != nil {
		messageText(client, ctx, to, err.Error())
	}
}
