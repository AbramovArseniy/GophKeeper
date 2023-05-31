package handlers

import (
	"context"
	"encoding/json"
	"io"
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
	context := context.Background()
	var err error
	var storage storage.Storage

	if cfg.Database == nil {
		log.Fatal("No database connected")
	} else {
		storage, err = database.NewDatabase(context, cfg.DatabaseAddress)
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

func (s *Server) PostSaveDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeJSON {
		http.Error(w, "wrong content type", http.StatusBadRequest)
		log.Println("wrong content type:", r.Header.Get("Content-Type"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusInternalServerError)
		log.Println("error while reading request body:", err)
		return
	}
	defer r.Body.Close()
	var meta storage.InfoMeta
	err = json.Unmarshal(body, &meta)
	if err != nil {
		http.Error(w, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return
	}
	data := storage.NewInfo(meta.Type)
	if data == nil {
		http.Error(w, "wrong data type", http.StatusNotImplemented)
		log.Println("wrong data type")
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return
	}
	encData, err := data.MakeBinary()
	if err != nil {
		http.Error(w, "cannot encrypt data", http.StatusInternalServerError)
		log.Println("error while encrypting data:", err)
		return
	}
	s.Storage.SaveData(encData, meta)
	w.WriteHeader(http.StatusOK)
}

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
