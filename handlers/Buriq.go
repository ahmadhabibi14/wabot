package handlers

import (
	"errors"
	"os"

	"github.com/ahmadhabibi14/wabot/utils"
	"github.com/disintegration/imaging"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func Buriq(client *whatsmeow.Client, img *waProto.ImageMessage, id string) ([]byte, string, error) {
	data, err := client.Download(img)
	if err != nil {
		return nil, ``, errors.New("error download file")
	}

	path, err := utils.SaveImage(id, data)
	if err != nil {
		return nil, ``, errors.New("cannot save image")
	}
	defer func() {
		os.Remove(path)
	}()

	src, err := imaging.Open(path)
	if err != nil {
		return nil, ``, errors.New("cannot open file")
	}

	if img.GetHeight() != img.GetWidth() {
		width := img.GetWidth()
		height := img.GetHeight()
		if width > height {
			width = height
		} else {
			height = width
		}
		src = imaging.CropAnchor(src, int(width), int(height), imaging.Center)
	}

	src = imaging.Resize(src, 100, 100, imaging.Lanczos)
	if err != nil {
		return nil, ``, errors.New("cannot resize image")
	}

	err = imaging.Save(src, path)
	if err != nil {
		return nil, ``, errors.New("cannot save file")
	}

	stc, err := os.ReadFile(path)
	if err != nil {
		return nil, ``, errors.New("failed to open file")
	}

	return stc, `Gambar buriq siyapp 😼😼`, nil
}
