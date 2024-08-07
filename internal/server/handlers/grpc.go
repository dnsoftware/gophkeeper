// Package handlers обмен данными с клиентом по gRPC протоколу
package handlers

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/field"
)

// UserService интерфейс для работы с регистрацией и аутентификацией/авторизацией
type UserService interface {
	// Registration регистрация нового пользователя. Возвращает токен доступа в случае удачи и ошибку, если что-то пошло не так
	Registration(ctx context.Context, login string, password string, repeatPassword string) (string, error)

	// Login вход пользователя. Возвращает токен доступа в случае удачи и ошибку, если что-то пошло не так
	Login(ctx context.Context, login string, password string) (string, error)
}

// EntityCodeService интерфейс для работы со справочником сущностей
type EntityCodeService interface {
	// EntityCodes запрос списка доступных к добавлению типов сущностей (таблица entity_codes)
	EntityCodes(ctx context.Context) (map[string]string, error)
}

// FieldsService интерфейс для работа с полями свойств сущностей
type FieldService interface {
	// Fields запрос списка характеристик полей сущности
	Fields(ctx context.Context, etype string) ([]field.EntityFields, error)
}

type EntityService interface {
	// AddEntity добавить сущность
	AddEntity(ctx context.Context, entity entity.EntityModel) (int32, error)
	// SaveEditEntity сохранить отредактированную
	SaveEditEntity(ctx context.Context, entity entity.EntityModel) error
	// DeleteEntity удалить сущность
	DeleteEntity(ctx context.Context, id int32, userID int32) error
	// Entity Получить сущность
	Entity(ctx context.Context, id int32) (*entity.EntityModel, error)
	// EntityList Список сущностей определенного типа для пользователя
	EntityList(ctx context.Context, etype string, userID int32) (map[int32]string, error)

	// UploadBinary потоковая загрузка незашифрованного бинарного файла
	UploadBinary(stream pb.Keeper_UploadBinaryServer) (int32, error)
	// DownloadBinary потоковая отдача незашифрованного бинарного файла
	DownloadBinary(entityID int32, stream pb.Keeper_DownloadBinaryServer) error

	// UploadCryptoBinary потоковая загрузка зашифрованного бинарного файла
	UploadCryptoBinary(stream pb.Keeper_UploadCryptoBinaryServer) (int32, error)
	// DownloadCryptoBinary потоковая отдача зашифрованного бинарного файла
	DownloadCryptoBinary(entityID int32, stream pb.Keeper_DownloadCryptoBinaryServer) error
}

// Services сервисы
type Services struct {
	UserService       UserService       // работа с регистрацией и аутентификацией/авторизацией
	EntityCodeService EntityCodeService // работа с данными пользователя (сохранение, получение, изменение)
	FieldService      FieldService      // работа с полями свойств сущностей
	EntityService     EntityService     // работа с сущностью
}

// GRPCServer gRPC сервер
type GRPCServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedKeeperServer

	svs Services // набор сервисов для работы с бизнес логикой

	Server *grpc.Server // пакет обеспечивающий работу gRPC сервера
}

func NewGRPCServer(services Services, certificateKeyPath string, privateKeyPath string) (*grpc.Server, error) {

	server := &GRPCServer{
		svs: services,
	}

	var opts []grpc.ServerOption

	if certificateKeyPath != "" && privateKeyPath != "" {
		creds, err := credentials.NewServerTLSFromFile(certificateKeyPath, privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("could not load TLS keys for gRPC: %s", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(checkUserInterceptor))

	// создаём gRPC-сервер
	server.Server = grpc.NewServer(opts...)

	// регистрируем сервис
	pb.RegisterKeeperServer(server.Server, server)

	return server.Server, nil
}

// Ping проверка связи
func (g *GRPCServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {

	if in.Message == "ping" {
		return &pb.PingResponse{Message: "pong"}, nil
	}

	return nil, errors.New("bad ping")
}
