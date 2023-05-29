package commands

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func DallE(msg string) /*string*/ {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()
	// image as base64
	reqBase64 := openai.ImageRequest{
		Prompt:         msg,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
	respBase64, err := client.CreateImage(ctx, reqBase64)
	if err != nil {
		log.Printf("Image creation error :: %v\n", err)
	}
	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		log.Printf("Base64 decode error :: %v\n", err)
	}
	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		log.Printf("PNG decode error :: %v\n", err)
	}
	file, err := os.Create("img/generated/imageGenerated.png")
	if err != nil {
		log.Printf("File creation error :: %v\n", err)
	}
	defer file.Close()
	if err := png.Encode(file, imgData); err != nil {
		log.Printf("PNG encode error %v\n", err)
	}
	// var success string = "The image was saved in img/generated/imageGenerated.png"
	return
}
