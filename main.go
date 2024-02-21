package main

import (
	"os"

	"github.com/ahmadhabibi14/wabot/configs"
	"github.com/ahmadhabibi14/wabot/services"
	"github.com/ahmadhabibi14/wabot/utils"
	_ "github.com/mattn/go-sqlite3"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	botName := os.Getenv("BOT_NAME")
	utils.PrintBanner(botName)
	configs.LoadEnv()
	logLevel := os.Getenv("LOG_LEVEL")
	dbLog := waLog.Stdout("Database", logLevel, true)
	store := services.NewStore("sqlite3", "file:session.db?_foreign_keys=on", dbLog)
	botDevice := store.GetDevice()
	botLog := waLog.Stdout("Client", logLevel, true)
	bot := services.NewBot(botName, botLog, botDevice)
	bot.Start()
}
