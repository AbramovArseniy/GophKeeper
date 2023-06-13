package storage

import (
	"errors"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	ErrUserExists   = errors.New("such user already exist in DB")
	ErrInvalidData  = errors.New("error data is invalid")
	ErrDataNotFound = errors.New("error data not found")
)

type Storage interface {
	SaveData(encryptedData []byte, metadata InfoMeta) error
	GetData(metadata InfoMeta) ([]byte, error)
}

type UserStorage interface {
	FindUser(login string) (*types.User, error)
	RegisterNewUser(login string, password string) (types.User, error)
	GetUserData(login string) (types.User, error)
}
