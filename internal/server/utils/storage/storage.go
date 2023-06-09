package storage

import (
	"errors"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	ErrInvalidData  = errors.New("error data is invalid")
	ErrDataNotFound = errors.New("error data not found")
)

type User struct {
	Id           int
	Login        string
	PasswordHash string
}

type Storage interface {
	SaveData(encryptedData []byte, metadata InfoMeta) error
	GetData(metadata InfoMeta) ([]byte, error)
}

type UserStorage interface {
	FindUser(login string) (*User, error)
	RegisterNewUser(login string, password string) (types.User, error)
	GetUserData(login string) (types.User, error)
}
