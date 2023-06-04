package storage

import (
	"errors"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
)

var (
	ErrInvalidData  = errors.New("data is invalid")
	ErrDataNotFound = errors.New("data is not found")
)

type Storage interface {
	SaveData(encryptedData []byte, metadata InfoMeta) error
	GetData(metadata InfoMeta) (Info, error)
	Close()
}

type UserStorage interface {
	RegisterNewUser(login string, password string) (types.User, error)
	GetUserData(login string) (types.User, error)
	CheckUserData(login, hash string) bool
}
