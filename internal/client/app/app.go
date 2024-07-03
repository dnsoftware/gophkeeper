package app

import (
	"os"

	"github.com/chzyer/readline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dnsoftware/gophkeeper/internal/client/config"
	"github.com/dnsoftware/gophkeeper/internal/client/domain"
	"github.com/dnsoftware/gophkeeper/internal/client/infrastructure"
	"github.com/dnsoftware/gophkeeper/logger"
)

func ClientRun() {

	cfg, err := config.NewClientConfig()
	_ = cfg
	if err != nil {
		logger.Log().Fatal(err.Error())
	}
	logger.Log().Info("Client starting...")

	path, _ := os.Getwd()
	certFile := path + "/cert/ca.crt"

	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		logger.Log().Fatal(err.Error())
	}

	var opts []grpc.DialOption
	sender, _, err := infrastructure.NewGRPCSender(cfg.ServerAddress, cfg.SecretKey, creds, opts...)
	if err != nil {
		logger.Log().Fatal(err.Error())
	}

	rl, err := domain.NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: domain.FilterInput,
	})
	if err != nil {
		logger.Log().Fatal(err.Error())
	}

	client, err := domain.NewGophKeepClient(rl, sender)
	if err != nil {
		logger.Log().Fatal(err.Error())
	}

	client.Start()
}
