package handlers

import (
	"context"
	"encoding/json"
	"errors"
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

	storage, err = database.NewDatabase(context, cfg.DatabaseAddress)
	if err != nil {
		log.Println("error while creating new database:", err)
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
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusInternalServerError)
		log.Println("error while reading request body:", err)
		return
	}
	var meta storage.InfoMeta
	err = json.Unmarshal(body, &meta)
	if err != nil {
		http.Error(w, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return
	}
	data := storage.NewInfo(meta.Type)
	if data == nil {
		http.Error(w, "wrong data type", http.StatusBadRequest)
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
	err = s.Storage.SaveData(encData, meta)
	if errors.Is(err, storage.ErrInvalidData) {
		http.Error(w, "invalid data", http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, "cannot save data to database", http.StatusInternalServerError)
		log.Println("error while  saving data to database:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetDataByTypeHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetAllUsersDataHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetDataByNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeJSON {
		http.Error(w, "wrong content type", http.StatusBadRequest)
		log.Println("wrong content type:", r.Header.Get("Content-Type"))
		return
	}
	defer r.Body.Close()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusInternalServerError)
		log.Println("error while reading request body:", err)
		return
	}
	var meta storage.InfoMeta
	err = json.Unmarshal(reqBody, &meta)
	if err != nil {
		http.Error(w, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return
	}
	if s.Auth != nil {
		meta.Login = s.Auth.GetUserLogin(r)
	} else {
		log.Println("no jwt auth")
	}
	info, err := s.Storage.GetData(meta)
	if errors.Is(err, storage.ErrDataNotFound) {
		http.Error(w, "no data found", http.StatusNotFound)
		return
	}
	if errors.Is(err, storage.ErrInvalidData) {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "cannot get data from database", http.StatusInternalServerError)
		log.Println("error while getting data from database:", err)
		return
	}
	respBody, err := json.MarshalIndent(&info, "  ", "")
	if err != nil {
		http.Error(w, "cannot marshal response body", http.StatusInternalServerError)
		log.Println("error while marshalling response body:", err)
		return
	}
	_, err = w.Write(respBody)
	if err != nil {
		http.Error(w, "cannot write response body", http.StatusInternalServerError)
		log.Println("error while writing response body:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Route() chi.Router {
	router := chi.NewRouter()
	router.Post("/user/register/", s.RegistHandler)
	router.Post("/user/login/", s.AuthHandler)
	router.Post("/user/add-data/", s.PostSaveDataHandler)
	router.Post("/user/get-data-by-type/", s.GetDataByTypeHandler)
	router.Get("/user/get-users-data/", s.GetAllUsersDataHandler)
	router.Post("/user/get-data-by-name/", s.GetDataByNameHandler)
	return router
}
