package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/AbramovArseniy/GophKeeper/internal/server/services"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/config"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage/database"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
	"github.com/go-chi/chi/v5"
)

const contentTypeJSON = "application/json"

// MetricServer has HTTP server info
type Server struct {
	Addr    string
	Storage storage.Storage
	Auth    types.Authorization
}

// NewServer creates new MetricServer
func NewServer(cfg config.Config) *Server {
	var (
		storage storage.Storage
		err     error
	)
	if cfg.Database == nil {
		log.Fatal("No database connected")
	} else {
		storage, err = database.NewDatabase(context.Background(), cfg.DatabaseAddress)
		if err != nil {
			log.Println("error while creating new database:", err)
		}
	}
	return &Server{
		Addr:    cfg.Address,
		Storage: storage,
	}
}

func (s *Server) RegistHandler(w http.ResponseWriter, r *http.Request) {
	httpStatus, token, err := services.RegistService(r, s.Auth)
	if err != nil {
		log.Println("error with register service:", err)
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(httpStatus)
}

func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	httpStatus, token, err := services.AuthService(r, s.Auth)
	if err != nil {
		log.Println("error with authentication service:", err)
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(httpStatus)
}

func (s *Server) PostSaveDataHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) GetDataByTypeHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) GetAllUsersDataHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) Route() chi.Router {
	router := chi.NewRouter()
	router.Post("/user/register", s.RegistHandler)
	router.Post("/user/login", s.AuthHandler)
	router.Post("/user/add-data", s.PostSaveDataHandler)
	router.Post("/user/get-data-by-type", s.GetDataByTypeHandler)
	router.Get("/user/get-users-data", s.GetAllUsersDataHandler)
	return router
}
