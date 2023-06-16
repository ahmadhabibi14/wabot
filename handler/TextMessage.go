package handler

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func TextMsg(client *whatsmeow.Client, v *events.Message, to types.JID, msg string) {
	client.SendMessage(context.Background(), to, &waProto.Message{
		Conversation: proto.String(msg),
	})
}
