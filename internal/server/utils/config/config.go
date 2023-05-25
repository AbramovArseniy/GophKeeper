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
}

const defaultAddress = "localhost"

// SetServerParams sets server config
func SetServerParams() (cfg Config) {
	var (
		flagAddress    string
		flagDataBase   string
		flagConfigFile string
		cfgFile        string
	)
	flag.StringVar(&flagAddress, "a", defaultAddress, "server_address")
	flag.StringVar(&flagDataBase, "d", "", "db_address")
	flag.StringVar(&flagConfigFile, "c", "", "config_as_json")
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
	cfg.Address, exists = os.LookupEnv("ADDRESS")
	if !exists {
		cfg.Address = flagAddress
	}
	cfg.DatabaseAddress, exists = os.LookupEnv("DATABASE_DSN")
	if !exists {
		cfg.DatabaseAddress = flagDataBase
	}
	return cfg
}
