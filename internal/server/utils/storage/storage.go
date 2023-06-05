package storage

import (
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
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
	GetData(metadata InfoMeta) ([]byte, error)
}

type UserStorage interface {
	FindUser(login string) (*User, error)
}

func Migrate(dsn string) error {
	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = "file://internal/server/utils/storage/database/migrations"
	}
	m, err := migrate.New(path, dsn+"&x-migrations-table=migrations")
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
