// Package config формирует клиентскую конфигурацию путем объединения параметров,
// полученных из командной строки и из конфигурационного файла /cmd/client/config.yaml
package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/dnsoftware/gophkeeper/logger"
)

// ClientConfig конфигурация клиента
type ClientConfig struct {
	Env           string `yaml:"env"`           // окружение (local, dev, prod)
	ServerAddress string `yaml:"serverAddress"` // адрес и порт сервера
	SecretKey     string `yaml:"secretKey"`     // ключ шифрования передаваемых данных
}

// NewClientConfig создание конфигурационной структуры
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
			logger.Log().Error(fmt.Sprintf("client config file not found: %s", err))
			return nil, err
		}

		err = yaml.Unmarshal(rawCfg, &cfg)
		if err != nil {
			logger.Log().Error(fmt.Sprintf("failed parsing client config file: %s", err))
			return nil, err
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
