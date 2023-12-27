package messagesStorageDirectory

import (
	"errors"
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/logger"
	"github.com/mailhedgehog/smtpMessage"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type directoryMessagesRepo struct {
	context *storageContext
	room    contracts.Room
}

func (repo *directoryMessagesRepo) Store(message *smtpMessage.SmtpMessage) (smtpMessage.MessageID, error) {
	b, err := io.ReadAll(message.ToReader())
	if err != nil {
		return "", err
	}
	if repo.context.perRoomLimit > 0 && repo.context.perRoomLimit <= repo.Count() {
		repo.context.storage.RoomsRepo().Delete(repo.room)
	}

	path := filepath.Join(repo.context.roomDirectory(repo.room), string(message.ID))
	err = os.WriteFile(path, b, 0660)

	logManager().Debug(fmt.Sprintf("New message saved at %s", path))

	return message.ID, err
}

func (repo *directoryMessagesRepo) List(query contracts.SearchQuery, offset, limit int) ([]smtpMessage.SmtpMessage, int, error) {
	if offset < 0 || limit < 0 {
		return nil, 0, errors.New("offset and limit should be >= 0")
	}

	dir, err := os.Open(repo.context.roomDirectory(repo.room))
	if err != nil {
		return nil, 0, err
	}
	defer dir.Close()

	unfilteredN, err := dir.Readdir(0)
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(unfilteredN, func(i, j int) bool {
		return unfilteredN[i].ModTime().After(unfilteredN[j].ModTime())
	})

	var n []os.FileInfo

	if len(query) > 0 {
	filtrationLoop:
		for i := range unfilteredN {
			msg, err := repo.Load(smtpMessage.MessageID(unfilteredN[i].Name()))
			if err != nil {
				continue
			}
			for criteria, queryValue := range query {
				queryValue = strings.ToLower(queryValue)
				switch criteria {
				case contracts.SearchParamTo:
					for _, t := range msg.To {
						if strings.Contains(strings.ToLower(t.Address()), queryValue) {
							n = append(n, unfilteredN[i])
							continue filtrationLoop
						}
					}
				case contracts.SearchParamFrom:
					if strings.Contains(strings.ToLower(msg.From.Address()), queryValue) {
						n = append(n, unfilteredN[i])
						continue filtrationLoop
					}
				case contracts.SearchParamContent:
					if strings.Contains(strings.ToLower(msg.GetOrigin()), queryValue) {
						n = append(n, unfilteredN[i])
						continue filtrationLoop
					}
				}
			}
		}
	} else {
		n = unfilteredN
	}

	messages, err := repo.parseList(n, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return messages, len(n), nil
}

func (repo *directoryMessagesRepo) Count() int {
	dir, err := os.Open(repo.context.roomDirectory(repo.room))
	logger.PanicIfError(err)
	defer dir.Close()
	n, _ := dir.Readdirnames(0)
	return len(n)
}

func (repo *directoryMessagesRepo) Delete(messageId smtpMessage.MessageID) error {
	return os.Remove(filepath.Join(repo.context.roomDirectory(repo.room), string(messageId)))
}

func (repo *directoryMessagesRepo) Load(messageId smtpMessage.MessageID) (*smtpMessage.SmtpMessage, error) {
	b, err := os.ReadFile(filepath.Join(repo.context.roomDirectory(repo.room), string(messageId)))
	if err != nil {
		return nil, err
	}

	m := smtpMessage.FromString(string(b), messageId)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (repo *directoryMessagesRepo) parseList(n []os.FileInfo, offset, limit int) ([]smtpMessage.SmtpMessage, error) {
	messages := make([]smtpMessage.SmtpMessage, 0)

	if offset >= len(n) {
		return messages, nil
	}

	endIndex := len(n)
	if offset+limit < len(n) {
		endIndex = offset + limit
	}
	n = n[offset:endIndex]

	for _, fileinfo := range n {
		b, err := os.ReadFile(filepath.Join(repo.context.roomDirectory(repo.room), fileinfo.Name()))
		if err != nil {
			logManager().Error(err.Error())
			continue
		}
		msg := smtpMessage.FromString(string(b), smtpMessage.MessageID(fileinfo.Name()))

		messages = append(messages, *msg)
	}

	logManager().Debug(fmt.Sprintf("Found %d messages", len(messages)))

	return messages, nil
}
