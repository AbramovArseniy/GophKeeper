package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
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
	err = dataBase.Migrate()
	if err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}
	return dataBase, nil
}

func (d *DataBase) Migrate() error {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../internal/server/utils/storage/database/migrations",
		d.dba, driver)
	if err != nil {
		return fmt.Errorf("could not create migration: %w", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
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

func (d *DataBase) FindUser(login string) (*types.User, error) {
	var user types.User
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
	err = row.Scan(&user.ID, &user.Login, &user.HashPassword)
	if err != nil {
		return nil, ErrInvalidUser
	}

	return &user, nil
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
				return types.User{}, storage.ErrUserExists
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
	if errors.Is(err, sql.ErrNoRows) {
		return types.User{}, nil
	}
	return user, err
}
