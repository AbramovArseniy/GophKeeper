package storage

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type User struct {
	Id           int
	Login        string
	PasswordHash string
}

type Storage interface {
	SaveData(encryptedData []byte, metadata InfoMeta) error
	GetData(metadata InfoMeta) (Info, error)
}

type UserStorage interface {
	FindUser(login string) (*User, error)
}
