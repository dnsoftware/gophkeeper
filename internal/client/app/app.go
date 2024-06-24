package app

import (
	"github.com/dnsoftware/gophkeeper/internal/client/config"
	"github.com/dnsoftware/gophkeeper/logger"
)

func ClientRun() {

	cfg, err := config.NewClientConfig()
	_ = cfg
	if err != nil {
		logger.Log().Fatal(err.Error())
	}
	logger.Log().Info("Client starting...")

	//client, err := domain.NewKeeperClient(cfg.ServerAddress, cfg.SecretKey, nil, nil)
	//if err != nil {
	//	logger.Log().Fatal(err.Error())
	//}
	//
	//client.Start()
}
