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

	width := img.GetWidth()
	height := img.GetHeight()
	ratio := 70
	if height < uint32(ratio) {
		ratio = int(height)
	}
	if width < uint32(ratio) {
		ratio = int(width)
	}

	if width == height {
		width = uint32(ratio)
		height = uint32(ratio)
	} else {
		if width > height {
			newRatio := ratio / int(height)
			width = width * uint32(newRatio)
			height = uint32(ratio)
		} else {
			newRatio := ratio / int(width)
			height = height * uint32(newRatio)
			width = uint32(ratio)
		}
	}

	src = imaging.Resize(src, int(width), int(height), imaging.Lanczos)
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
