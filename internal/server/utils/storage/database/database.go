package database

import "database/sql"

type Database struct{}

func NewDatabase(db *sql.DB) *Database {
	return &Database{}
}

func Migrate(db *sql.DB, dbAddress string) error {
	return nil
}
