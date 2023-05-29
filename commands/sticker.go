// package commands

// import (
// 	"fmt"
// 	"log"
// 	"go.mau.fi/whatsmeow"
// 	waProto "go.mau.fi/whatsmeow/binary/proto"
// )

// func Sticker(msg string) {
// 	stickerMsg := &waProto.StickerMessage{

// 	}
// }

/*
type StickerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url               *string      `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	FileSha256        []byte       `protobuf:"bytes,2,opt,name=fileSha256" json:"fileSha256,omitempty"`
	FileEncSha256     []byte       `protobuf:"bytes,3,opt,name=fileEncSha256" json:"fileEncSha256,omitempty"`
	MediaKey          []byte       `protobuf:"bytes,4,opt,name=mediaKey" json:"mediaKey,omitempty"`
	Mimetype          *string      `protobuf:"bytes,5,opt,name=mimetype" json:"mimetype,omitempty"`
	Height            *uint32      `protobuf:"varint,6,opt,name=height" json:"height,omitempty"`
	Width             *uint32      `protobuf:"varint,7,opt,name=width" json:"width,omitempty"`
	DirectPath        *string      `protobuf:"bytes,8,opt,name=directPath" json:"directPath,omitempty"`
	FileLength        *uint64      `protobuf:"varint,9,opt,name=fileLength" json:"fileLength,omitempty"`
	MediaKeyTimestamp *int64       `protobuf:"varint,10,opt,name=mediaKeyTimestamp" json:"mediaKeyTimestamp,omitempty"`
	FirstFrameLength  *uint32      `protobuf:"varint,11,opt,name=firstFrameLength" json:"firstFrameLength,omitempty"`
	FirstFrameSidecar []byte       `protobuf:"bytes,12,opt,name=firstFrameSidecar" json:"firstFrameSidecar,omitempty"`
	IsAnimated        *bool        `protobuf:"varint,13,opt,name=isAnimated" json:"isAnimated,omitempty"`
	PngThumbnail      []byte       `protobuf:"bytes,16,opt,name=pngThumbnail" json:"pngThumbnail,omitempty"`
	ContextInfo       *ContextInfo `protobuf:"bytes,17,opt,name=contextInfo" json:"contextInfo,omitempty"`
	StickerSentTs     *int64       `protobuf:"varint,18,opt,name=stickerSentTs" json:"stickerSentTs,omitempty"`
	IsAvatar          *bool        `protobuf:"varint,19,opt,name=isAvatar" json:"isAvatar,omitempty"`
}

func (x *StickerMessage) Reset() {
	*x = StickerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_binary_proto_def_proto_msgTypes[60]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
*/