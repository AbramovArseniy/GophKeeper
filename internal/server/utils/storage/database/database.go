package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataBase struct {
	db  *sql.DB
	ctx context.Context
	dba string
}

func NewDatabase(ctx context.Context, dba string) (*DataBase, error) {
	if dba == "" {
		err := fmt.Errorf("there is no DB address")
		return nil, err
	}
	db, err := sql.Open("pgx", dba)
	if err != nil {
		return nil, err
	}
	return &DataBase{
		db:  db,
		ctx: ctx,
		dba: dba,
	}, nil
}

func (d *DataBase) Migrate() {
	_, err := d.db.ExecContext(d.ctx, `CREATE TABLE IF NOT EXISTS users (
		id SERIAL UNIQUE,
		login VARCHAR UNIQUE NOT NULL,
		password_hash VARCHAR NOT NULL
	);`)
	if err != nil {
		log.Printf("error during create users %s", err)
	}

	_, err = d.db.ExecContext(d.ctx, `CREATE TABLE IF NOT EXISTS keeper (
		entry_num VARCHAR(255) PRIMARY KEY,
		entry_login VARCHAR(16),
		entry_pass VARCHAR(255),
		entry_text TEXT,
		entry_binary BYTEA,
		entry_bank BIGINT,
		login VARCHAR(16) NOT NULL,
		date_time TIMESTAMP NOT NULL
	);`)
	if err != nil {
		log.Printf("error during create orders %s", err)
	}
}

// var (
// 	ErrUserExists   = errors.New("such user already exist in DB")
// 	ErrScanData     = errors.New("error while scan user ID")
// 	ErrInvalidData  = errors.New("error user data is invalid")
// 	ErrKeyNotFound  = errors.New("error user ID not found")
// )
