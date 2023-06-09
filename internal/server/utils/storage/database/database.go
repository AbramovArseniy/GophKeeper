package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrUserExists         = errors.New("such user already exist in DB")
	ErrScanData           = errors.New("error while scan user ID")
	ErrInvalidUser        = errors.New("error user is invalid")
	ErrKeyNotFound        = errors.New("error user ID not found")
	selectDataStmt string = `SELECT data FROM keeper WHERE type=$1 AND login=$2 AND name=$3`
	selectUserStmt string = `SELECT id, login, password_hash FROM users WHERE login=$1`
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
	dataBase := &DataBase{
		db:  db,
		ctx: ctx,
		dba: dba,
	}
	dataBase.Migrate()
	return dataBase, nil
}

func (d *DataBase) Migrate() error {
	// 	_, err := d.db.ExecContext(d.ctx, `CREATE TABLE IF NOT EXISTS users (
	// 		id SERIAL UNIQUE,
	// 		login VARCHAR UNIQUE NOT NULL,
	// 		password_hash VARCHAR NOT NULL
	// 	);`)
	// 	if err != nil {
	// 		log.Printf("Error during create users %s", err)
	// 	}

	// 	_, err = d.db.ExecContext(d.ctx, `CREATE TABLE IF NOT EXISTS keeper(
	// 		login INT NOT NULL,
	// 		data BYTEA NOT NULL,
	// 		type SMALLINT NOT NULL,
	// 		name VARCHAR NOT NULL,
	// 		UNIQUE(login, type, name)
	// 	);`)
	// 	if err != nil {
	// 		log.Printf("Error during create keeper %s", err)
	// 	}
	path := "file://internal/server/utils/storage/database/migrations"
	m, err := migrate.New(path, d.dba+"&x-migrations-table=migrations")
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

func (d *DataBase) SaveData(encryptedData []byte, metadata storage.InfoMeta) error {
	_, err := d.db.ExecContext(d.ctx, `INSERT INTO keeper (data, login, type, name) VALUES ($1, $2, $3, $4)`,
		encryptedData, metadata.Login, metadata.Type, metadata.Name)
	if err != nil {
		return fmt.Errorf("error while inserting row into database: %w", err)
	}

	return nil
}

func (d *DataBase) GetData(metadata storage.InfoMeta) ([]byte, error) {
	var data []byte
	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		return nil, storage.ErrInvalidData
	}
	defer tx.Rollback()

	selectData, err := tx.PrepareContext(d.ctx, selectDataStmt)
	if err != nil {
		return nil, storage.ErrInvalidData
	}
	defer selectData.Close()

	row := selectData.QueryRowContext(d.ctx, metadata.Type, metadata.Login, metadata.Name)
	err = row.Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrDataNotFound
	}
	if err != nil {
		return nil, storage.ErrInvalidData
	}
	return data, nil
}

func (d *DataBase) Close() {
	d.db.Close()
}

func (d *DataBase) FindUser(login string) (*storage.User, error) {
	var user storage.User
	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		return nil, ErrInvalidUser
	}
	defer tx.Rollback()

	selectUser, err := tx.PrepareContext(d.ctx, selectUserStmt)
	if err != nil {
		return nil, ErrInvalidUser
	}
	defer selectUser.Close()

	row := selectUser.QueryRowContext(d.ctx, login)
	err = row.Scan(&user.Id, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, ErrInvalidUser
	}

	return &user, nil
}
