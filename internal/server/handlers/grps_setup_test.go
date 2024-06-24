// Package handlers содержит функции необходимые для предварительной настройки тестового окружения
package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"

	configclient "github.com/dnsoftware/gophkeeper/internal/client/config"
	domainclient "github.com/dnsoftware/gophkeeper/internal/client/infrastructure"
	"github.com/dnsoftware/gophkeeper/internal/server/config"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity_code"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/field"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/user"
	mock_domain "github.com/dnsoftware/gophkeeper/internal/server/mocks"
	"github.com/dnsoftware/gophkeeper/internal/storage/postgresql"
)

// "облегченный" вариант предварительных настроек
func setupLight(cfg config.ServerConfig) error {

	listen = bufconn.Listen(bufSize)
	repoUser := &mock_domain.MockUserStorage{}
	userService, _ := user.NewUser(repoUser)

	repoEntityCodeStorage := &mock_domain.MockEntityCodeStorage{}
	entityCodeService, _ := entity_code.NewEntityCode(repoEntityCodeStorage)

	repoFields := &mock_domain.MockFieldStorage{}
	fieldsService, _ := field.NewField(repoFields)

	repoEntity := &mock_domain.MockEntityRepo{}
	entityService, _ := entity.NewEntity(repoEntity, repoFields)
	server, err := NewGRPCServer(Services{userService, entityCodeService, fieldsService, entityService}, cfg.SertificateKeyPath, cfg.PrivateKeyPath)
	if err != nil {
		return errors.New("Not start GRPC server: " + err.Error())
	}

	go func() {
		if err := server.Serve(listen); err != nil {
			log.Fatalf("Test grpc server exited with error: %v", err)
		}
	}()

	return nil
}

// setupFull полная настройка тестового окружения
// SSL сертификаты, тестовая база Postgresql, gRPC сервер, gRPC клиент
func setupFull(cfg config.ServerConfig, cfgClient configclient.ClientConfig) (*domainclient.GRPCSender, *grpc.ClientConn, error) {

	certDir := getTestsCertDir()

	cfg.SertificateKeyPath = certDir + "/server.crt"
	cfg.PrivateKeyPath = certDir + "/server.key"

	listen = bufconn.Listen(bufSize)

	repository, err := setupDatabase()
	if err != nil {
		return nil, nil, err
	}
	userService, _ := user.NewUser(repository)
	entityCodeService, _ := entity_code.NewEntityCode(repository)
	fieldService, _ := field.NewField(repository)
	entityService, _ := entity.NewEntity(repository, repository)
	server, err := NewGRPCServer(Services{userService, entityCodeService, fieldService, entityService}, cfg.SertificateKeyPath, cfg.PrivateKeyPath)
	if err != nil {
		return nil, nil, errors.New("Not start GRPC server: " + err.Error())
	}

	go func() {
		if err := server.Serve(listen); err != nil {
			log.Fatalf("Test grpc server exited with error: %v", err)
		}
	}()

	//client, conn, err := NewTestClient(cfg)
	creds, err := clientTLSCreds()
	if err != nil {
		return nil, nil, err
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithContextDialer(bufDialer))
	client, conn, err := domainclient.NewKeeperClient(cfg.ServerAddress, cfgClient.SecretKey, creds, opts...)

	return client, conn, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return listen.Dial()
}

func getTestsCertDir() string {
	path, _ := os.Getwd()
	sep := "internal"
	parts := strings.Split(path, sep)

	return parts[0] + sep + "/certs"
}

func clientTLSCreds() (credentials.TransportCredentials, error) {

	certFile := getTestsCertDir() + "/ca.crt"

	return credentials.NewClientTLSFromFile(certFile, "")
}

// setupDatabase настройка чистой тестовой БД с начальными данными
func setupDatabase() (*postgresql.PgStorage, error) {
	ctx := context.Background()

	dbname := "users"
	user := "user"
	password := "password"

	// Start the postgres container
	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	// очистка БД
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	query := "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// конец очистки БД

	pgs, err := postgresql.NewPostgresqlStorage(dsn)
	if err != nil {
		return nil, err
	}

	// заполняем базу через миграции
	// не забываем импортировать
	// _ "github.com/golang-migrate/migrate/v4/source/file"
	// _ "github.com/golang-migrate/migrate/v4/database/postgres"

	path, _ := os.Getwd()
	parts := strings.Split(path, "internal")
	sourceURL := "file://" + parts[0] + "cmd/server/migrations"
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return nil, err
	} else {
		err = m.Up()
		if err != nil && err.Error() != "no change" {
			return nil, err
		}

	}

	return pgs, nil
}
