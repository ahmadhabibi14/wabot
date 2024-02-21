package services

import (
	"github.com/ahmadhabibi14/wabot/utils"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Store struct {
	Log     waLog.Logger
	DB      string
	Address string
}

func NewStore(dbName, dbAddress string, dbLog waLog.Logger) *Store {
	return &Store{
		Log:     dbLog,
		DB:      dbName,
		Address: dbAddress,
	}
}

func (s *Store) GetDevice() *store.Device {
	container, err := sqlstore.New(s.DB, s.Address, s.Log)
	utils.PanicIfError(err)

	device, err := container.GetFirstDevice()
	utils.PanicIfError(err)

	return device
}
