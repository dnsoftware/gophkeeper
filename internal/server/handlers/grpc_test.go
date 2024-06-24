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
	res, err := client.Registration(ctx, &pb.RegisterRequest{
		Login:          "username",
		Password:       "userpass",
		RepeatPassword: "userpassBad",
	})

	require.NoError(t, err)
	require.Equal(t, constants.ErrPasswordsNotMatch, res.Error)

	// корректные данные, должно возвратить токен
	res, err = client.Registration(ctx, &pb.RegisterRequest{
		Login:          "username",
		Password:       "userpass",
		RepeatPassword: "userpass",
	})

	require.NoError(t, err)
	require.Equal(t, res.Error, "")
	require.NotEmpty(t, res.Token)

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
	_, err = client.Registration(ctx, &pb.RegisterRequest{
		Login:          login,
		Password:       password,
		RepeatPassword: password,
	})
	require.NoError(t, err)

	// вход с неправильным паролем
	logResp, err := client.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: "badpass",
	})
	require.NoError(t, err)
	require.Equal(t, logResp.Token, "")
	require.Equal(t, logResp.Error, constants.ErrBadPassword)

	// вход с неправильным логином
	logResp, err = client.Login(ctx, &pb.LoginRequest{
		Login:    "bad",
		Password: password,
	})
	require.NoError(t, err)
	require.Equal(t, logResp.Token, "")
	require.Equal(t, logResp.Error, constants.ErrNoSuchUser)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	logResp, err = client.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, logResp.Token)
	require.Empty(t, logResp.Error)
	require.NotEmpty(t, client.GetToken())

}

// TestEntityCodes получение справочника кодов сущностей с сервера
func TestEntityCodes(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	resp, err := client.EntityCodes(ctx, &pb.EntityCodesRequest{Token: client.GetToken()})
	require.NoError(t, err)
	require.NoError(t, err)
	require.Greater(t, len(resp.EntityCodes), 0)

}

// TestFields получение набора характеристик полей сущности с сервера
func TestFields(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	resp, err := client.Fields(ctx, &pb.FieldsRequest{Etype: "card"})
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, len(resp.Fields), 3)

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
	_, err = client.Registration(ctx, &pb.RegisterRequest{
		Login:          login,
		Password:       password,
		RepeatPassword: password,
	})
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	logResp, err := client.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: password,
	})
	md := metadata.New(map[string]string{"token": logResp.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// добавление банковской карты
	var props []*pb.Property
	var metainfo []*pb.Metainfo

	props = append(props,
		&pb.Property{
			FieldId: 3,
			Value:   "1111222233334444",
		},
		&pb.Property{
			FieldId: 4,
			Value:   "12/25",
		},
		&pb.Property{
			FieldId: 5,
			Value:   "123",
		})

	metainfo = append(metainfo,
		&pb.Metainfo{
			Title: "Владелец карты",
			Value: "Василий Пупкин",
		},
		&pb.Metainfo{
			Title: "Банк",
			Value: "Суслик Инвест",
		})

	entreq := &pb.AddEntityRequest{
		Id:       0,
		Etype:    constants.CardEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	resp, err := client.AddEntity(ctx, entreq)
	idEnt := resp.Id
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
		&pb.Property{
			FieldId: 7,
			Value:   onlyFilename,
		})

	metainfo = append(metainfo,
		&pb.Metainfo{
			Title: "Название картинки",
			Value: "Суслик в естественной среде обитания",
		})

	entreq = &pb.AddEntityRequest{
		Id:       0,
		Etype:    constants.BinaryEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	resp, err = client.AddEntity(ctx, entreq)
	idEnt = resp.Id

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
	_, err = client.Registration(ctx, &pb.RegisterRequest{
		Login:          login,
		Password:       password,
		RepeatPassword: password,
	})
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	logResp, err := client.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: password,
	})
	md := metadata.New(map[string]string{"token": logResp.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// добавление произвольных бинарных данных
	var props []*pb.Property
	var metainfo []*pb.Metainfo

	p, _ := os.Getwd()
	parts := strings.Split(p, "internal")
	uploadFile := parts[0] + "cmd/client/testbinary/gopher.jpg"
	onlyFilename := filepath.Base(uploadFile)

	props = append(props,
		&pb.Property{
			FieldId: 7,
			Value:   onlyFilename,
		})

	metainfo = append(metainfo,
		&pb.Metainfo{
			Title: "Название картинки",
			Value: "Суслик в естественной среде обитания",
		})

	entreq := &pb.AddEntityRequest{
		Id:       0,
		Etype:    constants.BinaryEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	resp, err := client.AddEntity(ctx, entreq)
	idEnt := resp.Id

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
