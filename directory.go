package messagesStorageDirectory

import (
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/logger"
	"os"
	"path/filepath"
)

var configuredLogger *logger.Logger

func logManager() *logger.Logger {
	if configuredLogger == nil {
		configuredLogger = logger.CreateLogger("messagesStorageDirectory")
	}
	return configuredLogger
}

type StorageConfiguration struct {
	Path string `yaml:"path"`
}

type storageContext struct {
	path         string
	perRoomLimit int
	storage      *Directory
}

func (context *storageContext) roomDirectory(room contracts.Room) string {
	if len(room) <= 0 {
		room = "_default"
	}
	path := filepath.Join(context.path, string(room))
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, 0770)
		logger.PanicIfError(err)
	}

	return path
}

// Directory store messages in local directory
type Directory struct {
	context *storageContext
}

func CreateDirectoryStorage(config *StorageConfiguration, storageConfig *contracts.MessagesStorageConfiguration) *Directory {
	path := config.Path
	if len(path) <= 0 {
		dir, err := os.MkdirTemp("", "mailhedgehog_")
		logger.PanicIfError(err)
		path = dir
	}
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, 0770)
		logger.PanicIfError(err)
	}
	logManager().Debug(fmt.Sprintf("Mail storage directory path is '%s'", path))

	storage := &Directory{
		context: &storageContext{
			path:         path,
			perRoomLimit: storageConfig.PerRoomLimit,
		}}

	storage.context.storage = storage

	return storage
}

func (directory *Directory) RoomsRepo() contracts.RoomsRepo {
	return &directoryRoomsRepo{
		context: directory.context,
	}
}

func (directory *Directory) MessagesRepo(room contracts.Room) contracts.MessagesRepo {
	return &directoryMessagesRepo{
		context: directory.context,
		room:    room,
	}
}
