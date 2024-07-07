// Package config формирует серверную конфигурацию путем объединения параметров,
// полученных из командной строки и из конфигурационного файла /cmd/server/config.yaml
package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/dnsoftware/gophkeeper/logger"
)

type ServerConfig struct {
	Env                string `yaml:"env"`                // окружение (local, dev, prod)
	ServerAddress      string `yaml:"serverAddress"`      // адрес и порт на которых работает gRPC сервер
	DatabaseDSN        string `yaml:"databaseDSN"`        // параметры доступа к базе данных Postgresql
	SertificateKeyPath string `yaml:"sertificateKeyPath"` // путь к файлу сертификата
	PrivateKeyPath     string `yaml:"privateKeyPath"`     // путь к файлу с приватным ключом
}

func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flagCfg := &ServerConfig{}

	var configFile string
	flag.StringVar(&configFile, "c", "", "yaml server config file path")
	flag.StringVar(&flagCfg.Env, "e", "local", "environment (local, dev, prod)")
	flag.StringVar(&flagCfg.ServerAddress, "a", "", "server address")
	flag.StringVar(&flagCfg.DatabaseDSN, "d", "", "database DSN")
	flag.StringVar(&flagCfg.SertificateKeyPath, "s", "", "path to SSL sertificate key file")
	flag.StringVar(&flagCfg.PrivateKeyPath, "p", "", "path to SSL private key file")
	flag.Parse()

	if configFile != "" {
		rawCfg, err := os.ReadFile(configFile)
		if err != nil {
			logger.Log().Fatal(fmt.Sprintf("config file not found: %s", err))
		}

		err = yaml.Unmarshal(rawCfg, &cfg)
		if err != nil {
			logger.Log().Fatal(fmt.Sprintf("failed parsing config file: %s", err))
		}

	}

	// configs consolidate
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = flagCfg.ServerAddress
	}
	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = flagCfg.DatabaseDSN
	}
	if cfg.SertificateKeyPath == "" {
		cfg.SertificateKeyPath = flagCfg.SertificateKeyPath
	}
	if cfg.PrivateKeyPath == "" {
		cfg.PrivateKeyPath = flagCfg.PrivateKeyPath
	}

	return cfg, nil
}
