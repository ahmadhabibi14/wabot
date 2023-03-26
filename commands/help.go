package commands

import "strings"

func Help() string {
	const msg string = `
Ya, Ada yang bisa saya bantu ?
	
Akun ini adalah bot 🤖
	`
	return strings.TrimSpace(msg)
}
