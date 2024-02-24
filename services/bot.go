package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ahmadhabibi14/wabot/handlers"
	"github.com/ahmadhabibi14/wabot/utils"
	"github.com/disintegration/imaging"
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

var CMD_ImgToImg = map[string]func(img *waProto.ImageMessage, ctx context.Context, client *whatsmeow.Client, v *events.Message, to types.JID){
	`/buriq`: generateSticker,
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

		log.Println("IMAGE WIDTH:", img.GetWidth())
		log.Println("IMAGE HEIGHT:", img.GetHeight())
		if img.GetWidth() == img.GetHeight() {
			log.Println("IMAGE RATIO SAME")
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
				Height:        proto.Uint32(100),
				Width:         proto.Uint32(100),
			},
		})
		if err != nil {
			log.Println("ERROR send message:", err)
		}
	}
}

func sendImgBack(img *waProto.ImageMessage, ctx context.Context, client *whatsmeow.Client, v *events.Message, to types.JID) {
	switch img.GetCaption() {
	case `/sendback`:
		_, err := client.SendMessage(ctx, to, &waProto.Message{
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
	case `/resize`:
		data, err := client.Download(img)
		sendMessageIfError(err, client, ctx, to, "error download file")

		path, err := utils.SaveImage(v.Info.ID, data)
		sendMessageIfError(err, client, ctx, to, "cannot save image")

		src, err := imaging.Open(path)
		sendMessageIfError(err, client, ctx, to, "cannot open file")

		if img.GetHeight() != img.GetWidth() {
			src = imaging.CropAnchor(src, 100, 100, imaging.Center)
		}

		src = imaging.Resize(src, 100, 100, imaging.Lanczos)
		sendMessageIfError(err, client, ctx, to, "cannot save file")

		stc, err := os.ReadFile(path)
		sendMessageIfError(err, client, ctx, to, "failed to open file")

		uploaded, err := client.Upload(ctx, stc, whatsmeow.MediaImage)
		sendMessageIfError(err, client, ctx, to, "failed to upload file")

		_, err = client.SendMessage(ctx, to, &waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				Url:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(http.DetectContentType(stc)),
				FileEncSha256: uploaded.FileEncSHA256,
				FileSha256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(stc))),
				ContextInfo: &waProto.ContextInfo{
					StanzaId:      &v.Info.ID,
					Participant:   proto.String(v.Info.Sender.String()),
					QuotedMessage: v.Message,
				},
				Caption: proto.String(`Gambar buriq siyapp`),
			},
		})
		if err != nil {
			log.Println("ERROR send message:", err)
		}
	}
}

func sendMessageIfError(err error, client *whatsmeow.Client, ctx context.Context, to types.JID, res string) {
	if err != nil {
		messageText(client, ctx, to, res)
	}
}
