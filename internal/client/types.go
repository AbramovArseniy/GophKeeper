package client

import (
	"context"
	"errors"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"google.golang.org/grpc/metadata"
)

var ErrExitCLI = errors.New("Exit")

type CommandLine struct {
	action *MDAct
}

type MDAct struct {
	act ClientAction
	md  *metadata.MD
}

// Дописать!
func (mda *MDAct) Connection(address string) error {
	return nil
}

type ClientAction interface {
	SavePassword(ctx context.Context, req storage.InfoLoginPass) error
	GetPassword(ctx context.Context, req GetRequest) (*storage.InfoLoginPass, error)
	SaveCard(ctx context.Context, req storage.InfoCard) error
	GetCard(ctx context.Context, req GetRequest) (*storage.InfoCard, error)
	SaveText(ctx context.Context, req storage.InfoText) error
	GetText(ctx context.Context, req GetRequest) (*storage.InfoText, error)
}

type GetRequest struct {
	Name string
}
