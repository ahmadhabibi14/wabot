package commands

var Messages = map[string]string{
	"/help":   Help(),
	"/chtgpt": ChatGPT(),
	"/gemini": GeminiAI(),
}

// var Text = map[string]func()string{
// 	"/mdjrney": MidJourney(),
// 	"/chtgpt":  ChatGPT(),
// 	"/dall-e":  DallE(),
// 	"/menu":    Menu(),
// 	"/sticker": Sticker(),
// }
