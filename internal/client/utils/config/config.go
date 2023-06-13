package config

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

type Config struct {
	ServerAddr string `json:"address"`
}

const defaultAddress = "localhost:8080"

func SetClientParams() (cfg Config) {
	var (
		flagAddress    string
		flagConfigFile string
		cfgFile        string
		exists         bool
	)
	flag.StringVar(&flagAddress, "a", defaultAddress, "server_address")
	flag.StringVar(&flagConfigFile, "c", "", "config_as_json")
	flag.Parse()
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
	cfg.ServerAddr, exists = os.LookupEnv("ADDRESS")
	if !exists {
		cfg.ServerAddr = flagAddress
	}
	return cfg
}
