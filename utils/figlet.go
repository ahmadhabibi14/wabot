package utils

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/probandula/figlet4go"
)

func PrintBanner(name string) {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontName = "larry3d"
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen, figlet4go.ColorYellow, figlet4go.ColorCyan,
		figlet4go.ColorRed, figlet4go.ColorMagenta,
	}
	ascii.LoadFont("/larry3d.flf")
	renderStr, _ := ascii.RenderOpts(name, options)
	fmt.Println(renderStr)
}
