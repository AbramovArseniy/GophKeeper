package handlers

import (
	"log"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/config"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage/database"
	"github.com/go-chi/chi/v5"
)

const contentTypeJSON = "application/json"

// MetricServer has HTTP server info
type Server struct {
	Addr    string
	Storage storage.Storage
}

// NewServer creates new MetricServer
func NewServer(cfg config.Config) *Server {
	var (
		storage storage.Storage
	)
	if cfg.Database == nil {
		log.Fatal("No database connected")
	} else {
		storage, _ = database.NewDatabase(cfg.Database)
	}
	return &Server{
		Addr:    cfg.Address,
		Storage: storage,
	}
}

func (s *Server) Route() chi.Router {
	router := chi.NewRouter()
	return router
}
