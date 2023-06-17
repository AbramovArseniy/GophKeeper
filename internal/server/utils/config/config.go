package config

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

type Config struct {
	Address         string `json:"address"`
	DatabaseAddress string `json:"database_dsn"`
	Database        *sql.DB
	JWTSecret       string `json:"jwt_secret"`
	SecretKey       string
}

const defaultAddress = "localhost:8080"

// SetServerParams sets server config
func SetServerParams() (cfg Config) {
	var (
		flagAddress    string
		flagDataBase   string
		flagConfigFile string
		flagJWTSecret  string
		flagSecretKey  string
		cfgFile        string
	)
	flag.StringVar(&flagAddress, "a", defaultAddress, "server_address")
	flag.StringVar(&flagDataBase, "d", "", "db_address")
	flag.StringVar(&flagConfigFile, "c", "", "config_as_json")
	flag.StringVar(&flagJWTSecret, "js", "", "jwt_secret_key")
	flag.StringVar(&flagSecretKey, "k", "", "secret_key_to_enc")
	flag.Parse()
	var exists bool
	if cfgFile, exists = os.LookupEnv("CONFIG"); !exists {
		cfgFile = flagConfigFile
	}
	if cfgFile != "" {
		file, err := os.Open(cfgFile)
		if err != nil {
			log.Println("error while opening config file:", err)
		}
		cfgJSON, err := io.ReadAll(file)
		if err != nil {
			log.Println("error while reading from config file:", err)
		}
		err = json.Unmarshal(cfgJSON, &cfg)
		if err != nil {
			log.Println("error while unmarshalling config json:", err)
		}
	}
	cfg.JWTSecret, exists = os.LookupEnv("JWT_SECRET")
	if !exists {
		cfg.JWTSecret = flagJWTSecret
	}
	cfg.Address, exists = os.LookupEnv("ADDRESS")
	if !exists {
		cfg.Address = flagAddress
	}
	cfg.SecretKey, exists = os.LookupEnv("SECRET")
	if !exists {
		cfg.SecretKey = flagSecretKey
	}
	cfg.DatabaseAddress, exists = os.LookupEnv("DATABASE_DSN")
	if !exists {
		cfg.DatabaseAddress = flagDataBase
	}
	log.Println(cfg.JWTSecret, flagJWTSecret)
	return cfg
}
