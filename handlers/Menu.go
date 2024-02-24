package handlers

import (
	"os"
)

func Menu() string {
	botName := os.Getenv("BOT_NAME")
	msg := `*` + botName + `* - WhatsApp Bot

Here are command lists

*== General ==*

*/help* - Help
*/menu* - Main menu

*== Command to Text ==*

*/gemini <text>*
Generate text with Gemini AI

*/gpt <text>*
Generate text with ChatGPT

*== Image to Image ==*

*/buriq*
Convert image to low resolution 
`

	return msg
}
