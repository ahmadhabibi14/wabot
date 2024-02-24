package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func CreateStickerIMG(client *whatsmeow.Client, v *events.Message, data []byte) *waProto.Message {

	RawPath := fmt.Sprintf("tmp/%s%s", v.Info.ID, ".jpg")
	ConvertedPath := fmt.Sprintf("tmp/sticker/%s%s", v.Info.ID, ".webp")
	err := os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save image: %v", err)
	}
	exc := exec.Command("cwebp", "-q", "80", RawPath, "-o", ConvertedPath)
	err = exc.Run()
	if err != nil {
		fmt.Println("Failed to Convert Image to WebP")
	}
	createExif := fmt.Sprintf("webpmux -set exif %s %s -o %s", "tmp/exif/raw.exif", ConvertedPath, ConvertedPath)
	cmd := exec.Command("bash", "-c", createExif)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to set webp metadata", err)
	}
	stc, err := os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}
	uploaded, err := client.Upload(context.Background(), stc, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(stc)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &v.Info.ID,
				Participant:   proto.String(v.Info.Sender.String()),
				QuotedMessage: v.Message,
			},
		},
	}
}
