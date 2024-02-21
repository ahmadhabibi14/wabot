package handlers

import (
	"context"
	"strings"
)

func Help(ctx context.Context, in string) string {
	const msg string = `
Ya, Ada yang bisa saya bantu ?
Akun ini adalah bot 🤖
	`
	return strings.TrimSpace(msg)
}
