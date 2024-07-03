package handlers

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// несовпадающие пароли - должно выдавать ошибку
	_, err = client.Registration("username", "userpass", "userpassBad")

	require.Error(t, err)
	require.Equal(t, constants.ErrPasswordsNotMatch, err.Error())

	// корректные данные, должно возвратить токен
	token, err := client.Registration("username", "userpass", "userpass")

	require.NoError(t, err)
	require.NotEmpty(t, token)

}

// TestLogin вход на сервер и получение токена
func TestLogin(t *testing.T) {

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// вход с неправильным паролем
	token, err := client.Login(login, "badpass")
	require.Error(t, err)
	require.Equal(t, token, "")
	require.Equal(t, err.Error(), constants.ErrBadPassword)

	// вход с неправильным логином
	token, err = client.Login("bad", password)
	require.Error(t, err)
	require.Equal(t, token, "")
	require.Equal(t, err.Error(), constants.ErrNoSuchUser)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err = client.Login(login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, client.GetToken())

}

// TestEntityCodes получение справочника кодов сущностей с сервера
func TestEntityCodes(t *testing.T) {

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	_, err = client.Login(login, password)
	require.NoError(t, err)

	entCodes, err := client.EntityCodes()
	require.NoError(t, err)
	require.NoError(t, err)
	require.Greater(t, len(entCodes), 0)

}

// TestFields получение набора характеристик полей сущности с сервера
func TestFields(t *testing.T) {

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	_, err = client.Login(login, password)
	require.NoError(t, err)

	fields, err := client.Fields("card")
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
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err := client.Login(login, password)
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

	idEnt, err := client.AddEntity(entreq)
	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// Получение добавленной банковской карты
	ent, err := client.Entity(idEnt)
	_ = ent
	require.NoError(t, err)
	require.Equal(t, "1111222233334444", ent.Props[0].Value)
	require.Equal(t, "Владелец карты", ent.Metainfo[0].Title)
	require.Equal(t, "Василий Пупкин", ent.Metainfo[0].Value)

	// Редактирование
	for pkey, pval := range ent.Props {
		switch pval.FieldId {
		case 3:
			ent.Props[pkey].Value = "5555222233334444"
		case 4:
			ent.Props[pkey].Value = "01/26"
		case 5:
			ent.Props[pkey].Value = "321"
		}
	}

	var metainfo2 []*domain.Metainfo
	metainfo2 = append(metainfo2,
		&domain.Metainfo{
			Title: "Владелец карты 2",
			Value: "Василий Пупкин 2",
		},
		&domain.Metainfo{
			Title: "Банк 2",
			Value: "Суслик Инвест 2",
		})
	ent.Metainfo = metainfo2
	idEnt2, err := client.SaveEntity(*ent)
	require.NoError(t, err)
	require.Equal(t, ent.Id, idEnt2)

	// Проверка отредактированного
	ent2, err := client.Entity(idEnt2)
	for _, pval := range ent2.Props {
		switch pval.FieldId {
		case 3:
			assert.Equal(t, "5555222233334444", pval.Value)
		case 4:
			assert.Equal(t, "01/26", pval.Value)
		case 5:
			assert.Equal(t, "321", pval.Value)
		}
	}

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

	idEnt, err = client.AddEntity(entreq)

	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// теперь после заведения записи на сервере загружаем бинарник на сервер
	size, err := client.UploadBinary(idEnt, uploadFile)
	require.NoError(t, err)
	require.Greater(t, size, int32(0))

	// получаем данные бинарной сущности и загружаем бинарник с сервера
	entBin, err := client.Entity(idEnt)
	require.NoError(t, err)

	fd := &entity.BinaryFileProperty{}
	err = json.Unmarshal([]byte(entBin.Props[0].Value), fd)
	require.NoError(t, err)

	downloadFile, err := client.DownloadBinary(idEnt, fd.Clientname)
	require.NoError(t, err)
	require.NotEmpty(t, downloadFile)

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
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err := client.Login(login, password)
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

	idEnt, err := client.AddEntity(entreq)

	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// теперь после заведения записи на сервере загружаем бинарник на сервер
	size, err := client.UploadCryptoBinary(idEnt, uploadFile)
	require.NoError(t, err)
	require.Greater(t, size, int32(0))

	// получаем данные бинарной сущности и загружаем бинарник с сервера
	entBin, err := client.Entity(idEnt)
	require.NoError(t, err)

	fd := &entity.BinaryFileProperty{}
	err = json.Unmarshal([]byte(entBin.Props[0].Value), fd)
	require.NoError(t, err)

	downloadFile, err := client.DownloadCryptoBinary(idEnt, fd.Clientname)
	require.NoError(t, err)
	require.NotEmpty(t, downloadFile)
}

// Добавление и редактирование бинарных данных
func TestEditBinaryEntity(t *testing.T) {
	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err := client.Login(login, password)
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

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

	idEnt, err := client.AddEntity(entreq)

	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	// теперь после заведения записи на сервере загружаем бинарник на сервер
	size, err := client.UploadCryptoBinary(idEnt, uploadFile)
	require.NoError(t, err)
	require.Greater(t, size, int32(0))

	// получаем данные бинарной сущности и загружаем бинарник с сервера
	entBin, err := client.Entity(idEnt)
	require.NoError(t, err)

	fd := &entity.BinaryFileProperty{}
	err = json.Unmarshal([]byte(entBin.Props[0].Value), fd)
	require.NoError(t, err)

	downloadFile, err := client.DownloadCryptoBinary(idEnt, fd.Clientname)
	require.NoError(t, err)
	require.NotEmpty(t, downloadFile)

	// Редактирование
	uploadFile2 := parts[0] + "cmd/client/testbinary/gopher2.png"
	onlyFilename2 := filepath.Base(uploadFile2)
	entBin.Props[0].Value = onlyFilename2

	var metainfo2 []*domain.Metainfo
	metainfo2 = append(metainfo2,
		&domain.Metainfo{
			Title: "Картинка 2",
			Value: "Гофер 2",
		})
	entBin.Metainfo = metainfo2
	idEnt2, err := client.SaveEntity(*entBin)
	require.NoError(t, err)
	require.Equal(t, entBin.Id, idEnt2)

	// после редактирования записи на сервере загружаем бинарник на сервер
	size, err = client.UploadCryptoBinary(idEnt2, uploadFile2)
	require.NoError(t, err)
	require.Greater(t, size, int32(0))

	// получаем данные бинарной сущности и загружаем бинарник с сервера
	entBin, err = client.Entity(idEnt2)
	require.NoError(t, err)

	fd = &entity.BinaryFileProperty{}
	err = json.Unmarshal([]byte(entBin.Props[0].Value), fd)
	require.NoError(t, err)

	downloadFile, err = client.DownloadCryptoBinary(idEnt2, fd.Clientname)
	require.NoError(t, err)
	require.NotEmpty(t, downloadFile)

}

// TestEntityList получение списка сущностей пользователя определенного типа
func TestEntityList(t *testing.T) {

	ctx := context.Background()

	client, conn, err := setupFull(cfg, cfgClient)
	require.NoError(t, err)
	defer conn.Close()

	// регистрируем пользователя
	login := "username"
	password := "userpass"
	_, err = client.Registration(login, password, password)
	require.NoError(t, err)

	// логин с корректными данными, должны получить токен доступа и выставить токен в клиенте
	token, err := client.Login(login, password)
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// добавление логина/пароля
	var props []*domain.Property
	var metainfo []*domain.Metainfo

	props = append(props,
		&domain.Property{
			FieldId: 1,
			Value:   "login",
		},
		&domain.Property{
			FieldId: 2,
			Value:   "password",
		})

	metainfo = append(metainfo,
		&domain.Metainfo{
			Title: "Сайт",
			Value: "gopher.ru",
		},
		&domain.Metainfo{
			Title: "Куда",
			Value: "В админку например",
		})

	entreq := domain.Entity{
		Id:       0,
		Etype:    constants.LogopasEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	idEnt, err := client.AddEntity(entreq)
	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	props = nil
	metainfo = nil

	props = append(props,
		&domain.Property{
			FieldId: 1,
			Value:   "login2",
		},
		&domain.Property{
			FieldId: 2,
			Value:   "password2",
		})

	metainfo = append(metainfo,
		&domain.Metainfo{
			Title: "Сайт",
			Value: "elephant.ru",
		},
		&domain.Metainfo{
			Title: "Куда",
			Value: "В админку",
		})

	entreq = domain.Entity{
		Id:       0,
		Etype:    constants.LogopasEntity,
		Props:    props,
		Metainfo: metainfo,
	}

	idEnt, err = client.AddEntity(entreq)
	require.NoError(t, err)
	require.Greater(t, idEnt, int32(0))

	ent, err := client.Entity(idEnt)
	_ = ent

	// Получаем список
	list, err := client.EntityList("logopas")
	require.NoError(t, err)
	require.Equal(t, 2, len(list))
}
