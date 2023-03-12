package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/probandula/figlet4go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	// "github.com/sashabaranov/go-openai"
	openai "github.com/sashabaranov/go-openai"
)

var client *whatsmeow.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render("Habi-BOT")
	// set browser
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()
	store.DeviceProps.Os = proto.String("Habi-BOT")
	// Print Banner
	fmt.Println(renderStr)

}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsFromMe {
			if v.Message.GetConversation() != "" {
				// fmt.Println("Received a message!", v.Message.GetConversation())
				msg := v.Message.GetConversation()
				if strings.Contains(msg, "/ai") {
					openAIclient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
					resp, err := openAIclient.CreateChatCompletion(
						context.Background(),
						openai.ChatCompletionRequest{
							Model: openai.GPT3Dot5Turbo,
							Messages: []openai.ChatCompletionMessage{
								{
									Role:    openai.ChatMessageRoleUser,
									Content: msg,
								},
							},
						},
					)
					if err != nil {
						fmt.Printf("Chat Completion Error: %v\n", err)
						return
					}
					// Send Message
					client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
						Conversation: proto.String(strings.TrimSpace(resp.Choices[0].Message.Content)),
					})

					fmt.Println(strings.TrimSpace(resp.Choices[0].Message.Content))
				}
			}
		}
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Panic(err)
		// panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// LOGIN
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render QR Code
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				log.Println("Please scan to your whatsapp")
			} else {
				fmt.Println("Success Login !!!")
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
