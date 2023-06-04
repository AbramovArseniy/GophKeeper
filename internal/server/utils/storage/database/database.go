package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrUserExists            = errors.New("such user already exist in DB")
	ErrScanData              = errors.New("error while scan user ID")
	ErrKeyNotFound           = errors.New("error user ID not found")
	ErrInvalidData           = errors.New("error user data is invalid")
	ErrDataNotFound          = errors.New("error user data is invalid")
	selectDataStmt    string = `SELECT data from keeper WHERE type=$1 AND login=$2 AND name=$3`
	checkUserDatastmt string = `SELECT EXISTS(SELECT login, password_hash FROM users WHERE login = $1 AND password_hash = $2)`
	selectUserStmt    string = `SELECT id, login, password_hash FROM users WHERE login = $1`
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
		return fmt.Errorf("error while inserting row into database: %w", err)
	}

	return nil
}

func (d *DataBase) GetData(metadata storage.InfoMeta) (storage.Info, error) {
	var data []byte
	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		return nil, storage.ErrInvalidData
	}
	defer tx.Rollback()

	selectData, err := tx.PrepareContext(d.ctx, selectDataStmt)
	if err != nil {
		return nil, ErrInvalidData
	}
	defer selectData.Close()

	row := selectData.QueryRowContext(d.ctx, metadata.Type, metadata.Login, metadata.Name)
	err = row.Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrDataNotFound
	}
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

func (d *DataBase) Close() {
	d.db.Close()
}

func (d *DataBase) CheckUserData(login, hash string) bool {
	var exist bool
	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		log.Printf("error while creating tx %s", err)
		return false
	}

	defer tx.Rollback()

	checkUserDatastmt, err := tx.PrepareContext(d.ctx, checkUserDatastmt)
	if err != nil {
		log.Printf("error while creating stmt %s", err)
		return false
	}

	defer func() {
		if err := checkUserDatastmt.Close(); err != nil {
			log.Println("Error when close:", err)
		}
	}()

	row := checkUserDatastmt.QueryRowContext(d.ctx, login, hash)
	err = row.Scan(&exist)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return exist
	}
	if err != nil {
		log.Println(err)
		return exist
	}

	return exist
}

func (d *DataBase) RegisterNewUser(login string, password string) (types.User, error) {
	user := types.User{
		Login:        login,
		HashPassword: password,
	}
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) returning id`
	row := d.db.QueryRowContext(context.Background(), query, login, password)
	if err := row.Scan(&user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrKeyNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return types.User{}, ErrUserExists
			}
		}
		return types.User{}, ErrScanData
	}

	return user, nil
}

func (d *DataBase) GetUserData(login string) (types.User, error) {
	var user types.User

	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		return user, err
	}

	defer tx.Rollback()

	selectUserStmt, err := tx.PrepareContext(d.ctx, selectUserStmt)
	if err != nil {
		return user, err
	}

	defer func() {
		if err := selectUserStmt.Close(); err != nil {
			log.Println("Error when close:", err)
		}
	}()

	row := selectUserStmt.QueryRow(login)
	err = row.Scan(&user.ID, &user.Login, &user.HashPassword)
	if errors.Is(err, pgx.ErrNoRows) {
		return types.User{}, nil
	}
	return user, err
}
