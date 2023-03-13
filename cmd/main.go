package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahmadhabibi14/wabot/handler"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/probandula/figlet4go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontName = "larry3d"
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen, figlet4go.ColorYellow, figlet4go.ColorCyan,
		figlet4go.ColorRed, figlet4go.ColorMagenta,
	}
	ascii.LoadFont("/fonts/larry3d.flf")
	renderStr, _ := ascii.RenderOpts("Habi-BOT", options)
	// Print Banner
	fmt.Println(renderStr)
	// set browser
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()
	store.DeviceProps.Os = proto.String("Habi-BOT")
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
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	eventHandler := handler.EventHandler(client)
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
