package messagesStorageDirectory

import (
	"errors"
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/logger"
	"os"
	"sort"
)

type directoryRoomsRepo struct {
	context *storageContext
}

func (repo *directoryRoomsRepo) List(offset, limit int) ([]contracts.Room, error) {
	if offset < 0 || limit < 0 {
		return nil, errors.New("offset and limit should be >= 0")
	}

	dir, err := os.Open(repo.context.path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	n, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	rooms := make([]contracts.Room, 0)

	sort.Slice(n, func(i, j int) bool {
		return n[i].Name() < n[j].Name()
	})

	if offset >= len(n) {
		return rooms, nil
	}

	endIndex := len(n)
	if offset+limit < len(n) {
		endIndex = offset + limit
	}

	n = n[offset:endIndex]

	for _, fileinfo := range n {
		rooms = append(rooms, fileinfo.Name())
	}

	logManager().Debug(fmt.Sprintf("Found %d rooms", len(rooms)))

	return rooms, nil
}

func (repo *directoryRoomsRepo) Count() int {
	dir, err := os.Open(repo.context.path)
	logger.PanicIfError(err)
	defer dir.Close()
	n, _ := dir.Readdirnames(0)
	return len(n)
}

func (repo *directoryRoomsRepo) Delete(room contracts.Room) error {
	err := os.RemoveAll(repo.context.roomDirectory(room))
	if err != nil {
		return err
	}
	return nil
}
