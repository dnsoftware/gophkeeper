package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/dnsoftware/gophkeeper/internal/server/config"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity_code"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/field"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/user"
	"github.com/dnsoftware/gophkeeper/internal/server/handlers"
	"github.com/dnsoftware/gophkeeper/internal/storage/postgresql"
	"github.com/dnsoftware/gophkeeper/logger"
)

func ServerRun() error {
	cfg, err := config.NewServerConfig()
	if err != nil {
		logger.Log().Fatal("NewServerConfig: " + err.Error())
	}
	logger.Log().Info("Server starting...")

	// миграции
	path, _ := os.Getwd()
	m, err := migrate.New("file://"+path+"/migrations", cfg.DatabaseDSN)
	if err != nil {
		logger.Log().Error("migrate.New: " + err.Error())
		return err
	} else {
		err = m.Up()
		if err != nil && err.Error() != "no change" {
			logger.Log().Error("migrate.New else: " + err.Error())
			return err
		}
	}

	// grpc server
	listen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		logger.Log().Error("net.Listen: " + err.Error())
		return err
	}

	repository, err := postgresql.NewPostgresqlStorage(cfg.DatabaseDSN)
	if err != nil {
		logger.Log().Error("NewPostgresqlStorage: " + err.Error())
		return err
	}

	userService, err := user.NewUser(repository)
	if err != nil {
		logger.Log().Error("user.NewUser: " + err.Error())
		return err
	}

	entityCodeService, err := entity_code.NewEntityCode(repository)
	if err != nil {
		logger.Log().Error("NewEntityCode: " + err.Error())
		return err
	}
	fieldService, _ := field.NewField(repository)
	entityService, _ := entity.NewEntity(repository, repository)

	grpcServer, err := handlers.NewGRPCServer(handlers.Services{userService, entityCodeService, fieldService, entityService}, cfg.SertificateKeyPath, cfg.PrivateKeyPath)
	if err != nil {
		logger.Log().Fatal(err.Error())
	}
	fmt.Println("Сервер gRPC начал работу")

	go func() {
		if err = grpcServer.Serve(listen); err != nil {
			logger.Log().Fatal(err.Error())
		}
	}()

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		// корректное завершение работы gRPC сервера
		grpcServer.GracefulStop()
		fmt.Println("grpc server shutdown gracefully")

		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы
	fmt.Println("Everything shutdown gracefully")

	return nil // нормальное завершение
}
