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
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/crypto"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage/database"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const contentTypeJSON = "application/json"

// MetricServer has HTTP server info
type Server struct {
	Addr      string
	Storage   storage.Storage
	jwtSecret string
	Auth      types.Authorization
	SecretKey []byte
}

// NewServer creates new MetricServer
func NewServer(cfg config.Config) *Server {
	context := context.Background()
	var err error
	secret := []byte(cfg.SecretKey)
	db, err := database.NewDatabase(context, cfg.DatabaseAddress)
	if err != nil {
		log.Println("error while creating new database:", err)
	}
	return &Server{
		Addr:      cfg.Address,
		Storage:   db,
		SecretKey: secret,
		jwtSecret: cfg.JWTSecret,
		Auth:      NewAuth(context, db, cfg.JWTSecret),
	}
}

func (s *Server) RegistHandler(c echo.Context) error {
	httpStatus, token, err := services.RegistService(c.Request(), s.Auth)

	c.Response().Header().Set("Authorization", "Bearer "+token)
	c.Response().Writer.WriteHeader(httpStatus)

	return err
}

func (s *Server) AuthHandler(c echo.Context) error {
	httpStatus, token, err := services.AuthService(c.Request(), s.Auth)
	c.Response().Header().Set("Authorization", token)
	c.Response().Writer.WriteHeader(httpStatus)
	return err
}

func (s *Server) PostSaveDataHandler(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != contentTypeJSON {
		http.Error(c.Response().Writer, "wrong content type", http.StatusBadRequest)
		log.Println("wrong content type:", c.Request().Header.Get("Content-Type"))
		return nil
	}
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response().Writer, "cannot read request body", http.StatusInternalServerError)
		log.Println("error while reading request body:", err)
		return nil
	}
	var meta storage.InfoMeta
	err = json.Unmarshal(body, &meta)
	if err != nil {
		http.Error(c.Response().Writer, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return nil
	}
	data := storage.NewInfo(meta.Type)
	if data == nil {
		http.Error(c.Response().Writer, "wrong data type", http.StatusBadRequest)
		log.Println("wrong data type")
		return nil
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(c.Response().Writer, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return nil
	}
	binData, err := data.MakeBinary()
	if err != nil {
		http.Error(c.Response().Writer, "cannot make data binary", http.StatusInternalServerError)
		log.Println("error while making data binary:", err)
		return nil
	}
	encData, err := crypto.Encrypt(binData, s.SecretKey)
	if err != nil {
		http.Error(c.Response().Writer, "cannot encrypt data", http.StatusInternalServerError)
		log.Println("error while encrypting data:", err)
		return nil
	}
	if s.Auth != nil {
		meta.Login = s.Auth.GetUserLogin(c.Request())
	} else {
		log.Println("no jwt auth")
	}
	err = s.Storage.SaveData(encData, meta)
	if errors.Is(err, storage.ErrInvalidData) {
		http.Error(c.Response().Writer, "invalid data", http.StatusBadRequest)
	}
	if err != nil {
		http.Error(c.Response().Writer, "cannot save data to database", http.StatusInternalServerError)
		log.Println("error while  saving data to database:", err)
		return nil
	}
	c.Response().Writer.WriteHeader(http.StatusOK)
	return nil
}

// func (s *Server) GetDataByTypeHandler(c echo.Context) error {

// }

// func (s *Server) GetAllUsersDataHandler(c echo.Context) error {

// }

func (s *Server) GetDataByNameHandler(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != contentTypeJSON {
		http.Error(c.Response().Writer, "wrong content type", http.StatusBadRequest)
		log.Println("wrong content type:", c.Request().Header.Get("Content-Type"))
		return nil
	}
	defer c.Request().Body.Close()
	reqBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response().Writer, "cannot read request body", http.StatusInternalServerError)
		log.Println("error while reading request body:", err)
		return nil
	}
	var meta storage.InfoMeta
	err = json.Unmarshal(reqBody, &meta)
	if err != nil {
		http.Error(c.Response().Writer, "cannot unmarshal request body", http.StatusInternalServerError)
		log.Println("error while unmarshalling request body:", err)
		return nil
	}
	if s.Auth != nil {
		meta.Login = s.Auth.GetUserLogin(c.Request())
	} else {
		log.Println("no jwt auth")
	}
	encData, err := s.Storage.GetData(meta)
	if errors.Is(err, storage.ErrDataNotFound) {
		http.Error(c.Response().Writer, "no data found", http.StatusNotFound)
		return nil
	}
	if errors.Is(err, storage.ErrInvalidData) {
		http.Error(c.Response().Writer, "invalid data", http.StatusBadRequest)
		return nil
	}
	if err != nil {
		http.Error(c.Response().Writer, "cannot get data from database", http.StatusInternalServerError)
		log.Println("error while getting data from database:", err)
		return nil
	}
	binData, err := crypto.Decrypt(encData, s.SecretKey)
	if err != nil {
		http.Error(c.Response().Writer, "cannot decrypt data", http.StatusInternalServerError)
		log.Println("error while decrypting data:", err)
		return nil
	}
	data := storage.NewInfo(meta.Type)
	err = data.DecodeBinary(binData)
	if err != nil {
		log.Println("error while decoding binary:", err)
		http.Error(c.Response().Writer, "cannot decode data binary", http.StatusInternalServerError)
		return nil
	}
	respBody, err := json.MarshalIndent(&data, "  ", "")
	if err != nil {
		http.Error(c.Response().Writer, "cannot marshal response body", http.StatusInternalServerError)
		log.Println("error while marshalling response body:", err)
		return nil
	}
	_, err = c.Response().Writer.Write(respBody)
	if err != nil {
		http.Error(c.Response().Writer, "cannot write response body", http.StatusInternalServerError)
		log.Println("error while writing response body:", err)
		return nil
	}
	c.Response().Writer.WriteHeader(http.StatusOK)
	return nil
}

func (s *Server) Route() *echo.Echo {
	e := echo.New()

	e.POST("/user/auth/register/", s.RegistHandler)
	e.POST("/user/auth/login/", s.AuthHandler)

	logged := e.Group("/user", echojwt.WithConfig(echojwt.Config{SigningKey: []byte(s.jwtSecret)}))
	logged.POST("/add-data/", s.PostSaveDataHandler)
	//logged.POST("/get-data-by-type/", s.GetDataByTypeHandler)
	//logged.GET("/get-users-data/", s.GetAllUsersDataHandler)
	logged.POST("/get-data-by-name/", s.GetDataByNameHandler)

	return e
}
