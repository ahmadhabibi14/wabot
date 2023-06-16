package commands

var Message = map[string]string{
	"/help":   Help(),
	"/chtgpt": ChatGPT(),
}

// var Text = map[string]func()string{
// 	"/mdjrney": MidJourney(),
// 	"/chtgpt":  ChatGPT(),
// 	"/dall-e":  DallE(),
// 	"/menu":    Menu(),
// 	"/sticker": Sticker(),
// }
