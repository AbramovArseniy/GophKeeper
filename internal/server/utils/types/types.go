package types

import (
	"errors"
	"net/http"
)

type Authorization interface {
	GenerateToken(user User) (string, error)
	RegisterUser(userdata UserData) (User, error)
	LoginUser(userdata UserData) (User, error)
	GetUserID(r *http.Request) int
	GetUserLogin(r *http.Request) string
	CheckData(u UserData) error
}

type UserDB interface {
	RegisterNewUser(login string, password string) (User, error)
	GetUserData(login string) (User, error)
}

type User struct {
	Login        string
	HashPassword string
	ID           int
}
type UserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrUserExists   = errors.New("such user already exist in DB")
	ErrScanData     = errors.New("error while scan user ID")
	ErrInvalidData  = errors.New("error user data is invalid")
	ErrHashGenerate = errors.New("error can't generate hash")
	ErrKeyNotFound  = errors.New("error user ID not found")
	ErrAlarm        = errors.New("error tx.BeginTx alarm")
	ErrAlarm2       = errors.New("error tx.PrepareContext alarm")
)
