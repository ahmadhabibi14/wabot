package utils

import (
	"fmt"
	"os"
)

func SaveImage(fileName string, data []byte) (string, error) {
	rawPath := fmt.Sprintf("tmp/%s.jpg", fileName)
	err := os.WriteFile(rawPath, data, 0600)
	if err != nil {
		return ``, err
	}
	return rawPath, nil
}
