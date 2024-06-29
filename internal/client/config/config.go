package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/dnsoftware/gophkeeper/logger"
)

type ClientConfig struct {
	Env           string `yaml:"env"`           // окружение (local, dev, prod)
	ServerAddress string `yaml:"serverAddress"` // адрес и порт сервера
	SecretKey     string `yaml:"secretKey"`     // ключ шифрования передаваемых данных
}

func NewClientConfig() (*ClientConfig, error) {
	cfg := &ClientConfig{}
	flagCfg := &ClientConfig{}

	var configFile string
	flag.StringVar(&configFile, "c", "", "yaml client config file path")
	flag.StringVar(&flagCfg.Env, "e", "local", "environment (local, dev, prod)")
	flag.StringVar(&flagCfg.ServerAddress, "a", "", "server address")
	flag.StringVar(&flagCfg.SecretKey, "k", "", "secret key for encryption")
	flag.Parse()

	if configFile != "" {
		rawCfg, err := os.ReadFile(configFile)
		if err != nil {
			logger.Log().Fatal(fmt.Sprintf("client config file not found: %s", err))
		}

		err = yaml.Unmarshal(rawCfg, &cfg)
		if err != nil {
			logger.Log().Fatal(fmt.Sprintf("failed parsing client config file: %s", err))
		}

	}

	// configs consolidate
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = flagCfg.ServerAddress
	}
	if cfg.SecretKey == "" {
		cfg.SecretKey = flagCfg.SecretKey
	}

	return cfg, nil
}
