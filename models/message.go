package models

import (
	"sync"
)

type Message struct {
	Mu         sync.Mutex
	MsgReceive string
}
