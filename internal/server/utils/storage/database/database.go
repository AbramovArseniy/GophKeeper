package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrUserExists         = errors.New("such user already exist in DB")
	ErrScanData           = errors.New("error while scan user ID")
	ErrInvalidData        = errors.New("error user data is invalid")
	ErrKeyNotFound        = errors.New("error user ID not found")
	selectDataStmt string = `SELECT data from keeper WHERE type=$1 AND login=$2 AND name=$3`
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
		log.Printf("Error during create users %s", err)
	}

	_, err = d.db.ExecContext(d.ctx, `CREATE TABLE IF NOT EXISTS keeper(
		login INT NOT NULL,
		data BYTEA NOT NULL,
		type SMALLINT NOT NULL,
		name VARCHAR NOT NULL,
		UNIQUE(login, type, name)
	);`)
	if err != nil {
		log.Printf("Error during create keeper %s", err)
	}
}

func (d *DataBase) SaveData(encryptedData []byte, metadata storage.InfoMeta) error {
	_, err := d.db.ExecContext(d.ctx, `INSERT INTO keeper (data, login, type, name) VALUES ($1, $2, $3, $4)`,
		encryptedData, metadata.Login, metadata.Type, metadata.Name)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBase) GetData(metadata storage.InfoMeta) (storage.Info, error) {
	var data []byte
	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		return nil, ErrInvalidData
	}
	defer tx.Rollback()

	selectData, err := tx.PrepareContext(d.ctx, selectDataStmt)
	if err != nil {
		return nil, ErrInvalidData
	}
	defer selectData.Close()

	row := selectData.QueryRowContext(d.ctx, metadata.Type, metadata.Login, metadata.Name)
	err = row.Scan(&data)
	if err != nil {
		return nil, ErrInvalidData
	}
	info := storage.NewInfo(metadata.Type)
	err = info.DecodeBinary(data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding binary: %w", err)
	}
	return info, nil
}
