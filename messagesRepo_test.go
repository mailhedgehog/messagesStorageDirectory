package messagesStorageDirectory

import (
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/gounit"
	"github.com/mailhedgehog/logger"
	"github.com/mailhedgehog/smtpMessage"
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	room := "foo_bar"

	storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room + fmt.Sprint(i))).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)

		_, _ = storage.MessagesRepo(contracts.Room(room + fmt.Sprint(i))).Store(&smtpMessage.SmtpMessage{
			ID: smtpMessage.MessageID(fmt.Sprint(i + 1)),
		})
	}

	// Created 15 directories
	dir, err := os.Open(storage.context.path)
	logger.PanicIfError(err)
	defer dir.Close()
	n, _ := dir.Readdirnames(0)
	(*gounit.T)(t).AssertEqualsInt(15, len(n))

	// In each directory exists 2 messages
	dir, err = os.Open(storage.context.roomDirectory(contracts.Room(room + "2")))
	logger.PanicIfError(err)
	defer dir.Close()
	n, _ = dir.Readdirnames(0)
	(*gounit.T)(t).AssertEqualsInt(2, len(n))
}

func TestStore_ClearIfOverLimit(t *testing.T) {
	room := "foo_bar"

	storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 6})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room)).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	// Created 1 directory
	dir, err := os.Open(storage.context.path)
	logger.PanicIfError(err)
	defer dir.Close()
	n, _ := dir.Readdirnames(0)
	(*gounit.T)(t).AssertEqualsInt(1, len(n))

	// Messages 2 times deleted and then created
	dir, err = os.Open(storage.context.roomDirectory(contracts.Room(room)))
	logger.PanicIfError(err)
	defer dir.Close()
	n, _ = dir.Readdirnames(0)
	(*gounit.T)(t).AssertEqualsInt(3, len(n))
}

func TestCount(t *testing.T) {
	room := "foo_bar"

	storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room)).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	for i := 0; i < 4; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room + "2")).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	(*gounit.T)(t).AssertEqualsInt(15, storage.MessagesRepo(contracts.Room(room)).Count())
	(*gounit.T)(t).AssertEqualsInt(4, storage.MessagesRepo(contracts.Room(room+"2")).Count())
}

func TestDelete(t *testing.T) {
	room := "foo_bar"

	storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room)).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	for i := 0; i < 4; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room + "2")).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	(*gounit.T)(t).AssertEqualsInt(15, storage.MessagesRepo(contracts.Room(room)).Count())
	(*gounit.T)(t).AssertEqualsInt(4, storage.MessagesRepo(contracts.Room(room+"2")).Count())

	(*gounit.T)(t).AssertNotError(storage.MessagesRepo(contracts.Room(room)).Delete("3"))
	(*gounit.T)(t).AssertNotError(storage.MessagesRepo(contracts.Room(room + "2")).Delete("1"))
	(*gounit.T)(t).AssertNotError(storage.MessagesRepo(contracts.Room(room + "2")).Delete("2"))

	(*gounit.T)(t).AssertEqualsInt(14, storage.MessagesRepo(contracts.Room(room)).Count())
	(*gounit.T)(t).AssertEqualsInt(2, storage.MessagesRepo(contracts.Room(room+"2")).Count())
}

func TestLoad(t *testing.T) {
	room := "foo_bar"

	storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(contracts.Room(room)).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	msg, err := storage.MessagesRepo(contracts.Room(room)).Load("3")

	(*gounit.T)(t).AssertNotError(err)

	(*gounit.T)(t).AssertEqualsString("3", string(msg.ID))
}

func TestList(t *testing.T) {
	room := "foo_bar"

	storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}
		from, _ := smtpMessage.MessagePathFromString(fmt.Sprintf("<@foo,@bar,@baz:quix-%d@quib.com> <foo,bar>", i))
		msg.From = from
		to, _ := smtpMessage.MessagePathFromString(fmt.Sprintf("<foo-%d@quib.com>", i))
		msg.To = append(msg.To, to)
		msg.SetOrigin("")

		storedId, err := storage.MessagesRepo(contracts.Room(room)).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	msgs, count, err := storage.MessagesRepo(contracts.Room(room)).List(contracts.SearchQuery{
		contracts.SearchParamFrom: "quix-1",
	}, 0, 4)

	(*gounit.T)(t).AssertNotError(err)
	(*gounit.T)(t).AssertEqualsInt(6, count)
	(*gounit.T)(t).AssertEqualsInt(4, len(msgs))
	(*gounit.T)(t).AssertEqualsString("quix-10@quib.com", msgs[2].From.Address())
	(*gounit.T)(t).AssertEqualsString("foo-10@quib.com", msgs[2].To[0].Address())
}
