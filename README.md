# MailHedgehog package to store emails to file storage

All emails will be stored in physical files on server. Useful for simple implementation with small amount of emails.

## Usage

```go
storage := CreateDirectoryStorage(&StorageConfiguration{Path: ""}, &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})
msg, err := storage.MessagesRepo(room).Load("ID")
```

## Development

```shell
go mod tidy
go mod verify
go mod vendor
go test --cover
```

## Credits

- [![Think Studio](https://yaroslawww.github.io/images/sponsors/packages/logo-think-studio.png)](https://think.studio/)
