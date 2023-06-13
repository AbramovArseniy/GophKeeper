package types

import (
	"context"
	"errors"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
)

var ErrExitCLI = errors.New("Exit")

type ClientAction interface {
	SaveData(ctx context.Context, req storage.Info, infoType storage.InfoType) error
	GetData(ctx context.Context, req GetRequest) (storage.Info, error)
	Register(ctx context.Context, req AuthRequest) error
	Login(ctx context.Context, req AuthRequest) error
	Connect(address string) error
}

type GetRequest struct {
	Name string           `json:"name"`
	Type storage.InfoType `json:"type"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
