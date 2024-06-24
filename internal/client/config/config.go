package config

type ClientConfig struct {
	Env           string `yaml:"env"`           // окружение (local, dev, prod)
	ServerAddress string `yaml:"serverAddress"` // адрес и порт сервера
	SecretKey     string `yaml:"secretKey"`     // ключ шифрования передаваемых данных
}

func NewClientConfig() (*ClientConfig, error) {
	cfg := &ClientConfig{}

	return cfg, nil
}
