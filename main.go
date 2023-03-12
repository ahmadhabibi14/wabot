package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ahmadhabibi14/wabot/models"
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
					reqBody := models.Request{
						ModelRequest: "text-davinci-003",
						Prompt:       msg,
						Temperature:  1,
						MaxTokens:    100,
					}
					jsonBody, _ := json.Marshal(reqBody)
					req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonBody))
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
					if err != nil {
						panic(err)
					}
					// Send the request
					httpClient := &http.Client{}
					resp, err := httpClient.Do(req)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()

					// Read the response
					jsonData, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						panic(err)
					}

					var data models.TextCompletionResponse
					json.Unmarshal([]byte(jsonData), &data)

					client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
						Conversation: proto.String(strings.TrimSpace(data.Choices[0].Text)),
					})

					fmt.Println(strings.TrimSpace(data.Choices[0].Text))
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
