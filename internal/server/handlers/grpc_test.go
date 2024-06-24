package handlers

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	//"google.golang.org/grpc/test/bufconn"

	configclient "github.com/dnsoftware/gophkeeper/internal/client/config"
	"github.com/dnsoftware/gophkeeper/internal/client/domain"
	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/server/config"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
)

const bufSize = 1024 * 64

var listen *bufconn.Listener
var cfg = config.ServerConfig{
	ServerAddress:      "localhost:9090",
	DatabaseDSN:        "",
	SertificateKeyPath: "",
	PrivateKeyPath:     "",
}

var cfgClient = configclient.ClientConfig{
	Env:           "local",
	ServerAddress: "localhost:9090",
	SecretKey:     "secret",
}

// тестирование отклика сервера
func TestPing(t *testing.T) {
	setupLight(cfg)

	ctx := context.Background()
	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial : %v", err)
	}
	defer conn.Close()
	client := pb.NewKeeperClient(conn)

	resp, err := client.Ping(ctx, &pb.PingRequest{
		Message: "ping",
	})

	require.Equal(t, "pong", resp.Message)
	require.NoError(t, err)

}

// тестирование TLS соединения
func TestTLSCreds(t *testing.T) {

	certDir := getTestsCertDir()

	setupLight(config.ServerConfig{
		ServerAddress:      cfg.ServerAddress,
		SertificateKeyPath: certDir + "/server.crt",
		PrivateKeyPath:     certDir + "/server.key",
	})
	ctx := context.Background()

	creds, err := clientTLSCreds()
	if err != nil {
		t.Fatalf("Failed creds : %v", err)
	}

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("Failed to dial : %v", err)
	}
	defer conn.Close()
	client := pb.NewKeeperClient(conn)

	resp, err := client.Ping(ctx, &pb.PingRequest{
		Message: "ping",
	})
	require.NoError(t, err)

	require.Equal(t, "pong", resp.Message)
	require.NoError(t, err)

}

// неправильный сертификат дает ошибку при обращении к сервреру
func TestBadTLSCreds(t *testing.T) {

	certDir := getTestsCertDir()

	setupLight(config.ServerConfig{
		ServerAddress:      cfg.ServerAddress,
		SertificateKeyPath: certDir + "/badserver.crt",
		PrivateKeyPath:     certDir + "/server.key",
	})
	ctx := context.Background()

	creds, err := clientTLSCreds()
	if err != nil {
		t.Fatalf("Failed creds : %v", err)
	}

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("Failed to dial : %v", err)
	}
	defer conn.Close()
	client := pb.NewKeeperClient(conn)

	_, err = client.Ping(ctx, &pb.PingRequest{
		Message: "ping",
	})

	require.Error(t, err)

}

func TestRegistration(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// несовпадающие пароли - должно выдавать ошибку
	_, err = client.Registration(ctx, "username", "userpass", "userpassBad")

	require.Error(t, err)
	require.Equal(t, constants.ErrPasswordsNotMatch, err.Error())

	// корректные данные, должно возвратить токен
	token, err := client.Registration(ctx, "username", "userpass", "userpass")

	require.NoError(t, err)
	require.NotEmpty(t, token)

}

// TestLogin вход на сервер и получение токена
func TestLogin(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(ctx, login, password, password)
	require.NoError(t, err)

	// вход с неправильным паролем
	token, err := client.Login(ctx, login, "badpass")
	require.Error(t, err)
	require.Equal(t, token, "")
	require.Equal(t, err.Error(), constants.ErrBadPassword)

	// вход с неправильным логином
	token, err = client.Login(ctx, "bad", password)
	require.Error(t, err)
	require.Equal(t, token, "")
	require.Equal(t, err.Error(), constants.ErrNoSuchUser)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err = client.Login(ctx, login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, client.GetToken())

}

// TestEntityCodes получение справочника кодов сущностей с сервера
func TestEntityCodes(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	entCodes, err := client.EntityCodes(ctx)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Greater(t, len(entCodes), 0)

}

// TestFields получение набора характеристик полей сущности с сервера
func TestFields(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	fields, err := client.Fields(ctx, "card")
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, len(fields), 3)

}

// TestAddEntity добавление сущностей, позитивные сценарии
func TestAddEntity(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(ctx, login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err := client.Login(ctx, login, password)
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// добавление банковской карты
	var props []*domain.Property
	var metainfo []*domain.Metainfo

	props = append(props,
		&domain.Property{
			FieldId: 3,
			Value:   "1111222233334444",
		},
		&domain.Property{
			FieldId: 4,
			Value:   "12/25",
		},
		&domain.Property{
			FieldId: 5,
			Value:   "123",
		})

	metainfo = append(metainfo,
		&domain.Metainfo{
			Title: "Владелец карты",
			Value: "Василий Пупкин",
		},
		&domain.Metainfo{
			Title: "Банк",
			Value: "Суслик Инвест",
		})

	entreq := domain.Entity{
		Id:       0,
		Etype:    constants.CardEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	idEnt, err := client.AddEntity(ctx, entreq)
	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// Получение добавленной банковской карты
	req := &pb.EntityRequest{
		Id: idEnt,
	}
	ent, err := client.Entity(ctx, req)
	_ = ent
	require.NoError(t, err)
	require.Equal(t, "1111222233334444", ent.Props[0].Value)
	require.Equal(t, "Владелец карты", ent.Metainfo[0].Title)
	require.Equal(t, "Василий Пупкин", ent.Metainfo[0].Value)

	// добавление произвольных бинарных данных
	props = nil
	metainfo = nil

	p, _ := os.Getwd()
	parts := strings.Split(p, "internal")
	uploadFile := parts[0] + "cmd/client/testbinary/gopher.jpg"
	onlyFilename := filepath.Base(uploadFile)

	props = append(props,
		&domain.Property{
			FieldId: 7,
			Value:   onlyFilename,
		})

	metainfo = append(metainfo,
		&domain.Metainfo{
			Title: "Название картинки",
			Value: "Суслик в естественной среде обитания",
		})

	entreq = domain.Entity{
		Id:       0,
		Etype:    constants.BinaryEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	idEnt, err = client.AddEntity(ctx, entreq)

	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// теперь после заведения записи на сервере загружаем бинарник на сервер
	size, err := client.UploadBinary(ctx, idEnt, uploadFile)
	require.NoError(t, err)
	require.Greater(t, size, int32(0))

	// получаем данные бинарной сущности и загружаем бинарник с сервера
	reqBin := &pb.EntityRequest{
		Id: idEnt,
	}
	entBin, err := client.Entity(ctx, reqBin)
	require.NoError(t, err)

	fd := &entity.BinaryFileProperty{}
	err = json.Unmarshal([]byte(entBin.Props[0].Value), fd)
	require.NoError(t, err)

	downloadFile, err := client.DownloadBinary(ctx, idEnt, fd.Clientname)
	require.NoError(t, err)
	require.NotEmpty(t, downloadFile)

	// добавление логина/пароля

	// добавление произвольных текстовых данных
}

// TestAddCryptoBinary добавление сущностей с зашифрованными бинарными данными
func TestAddCryptoBinary(t *testing.T) {
	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(ctx, login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err := client.Login(ctx, login, password)
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// добавление произвольных бинарных данных
	var props []*domain.Property
	var metainfo []*domain.Metainfo

	p, _ := os.Getwd()
	parts := strings.Split(p, "internal")
	uploadFile := parts[0] + "cmd/client/testbinary/gopher.jpg"
	onlyFilename := filepath.Base(uploadFile)

	props = append(props,
		&domain.Property{
			FieldId: 7,
			Value:   onlyFilename,
		})

	metainfo = append(metainfo,
		&domain.Metainfo{
			Title: "Название картинки",
			Value: "Суслик в естественной среде обитания",
		})

	entreq := domain.Entity{
		Id:       0,
		Etype:    constants.BinaryEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	idEnt, err := client.AddEntity(ctx, entreq)

	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// теперь после заведения записи на сервере загружаем бинарник на сервер
	size, err := client.UploadCryptoBinary(ctx, idEnt, uploadFile)
	require.NoError(t, err)
	require.Greater(t, size, int32(0))

	// получаем данные бинарной сущности и загружаем бинарник с сервера
	reqBin := &pb.EntityRequest{
		Id: idEnt,
	}
	entBin, err := client.Entity(ctx, reqBin)
	require.NoError(t, err)

	fd := &entity.BinaryFileProperty{}
	err = json.Unmarshal([]byte(entBin.Props[0].Value), fd)
	require.NoError(t, err)

	downloadFile, err := client.DownloadCryptoBinary(ctx, idEnt, fd.Clientname)
	require.NoError(t, err)
	require.NotEmpty(t, downloadFile)
}
