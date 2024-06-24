package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/utils"
	"github.com/dnsoftware/gophkeeper/logger"
)

//type ClientSender interface {
//	Registration(login string, password string, password2 string) (string, error)
//}

type KeeperClient struct {
	//cfg *config.ClientConfig
	pb.KeeperClient
	//sender    ClientSender
	token     string
	password  string
	secretKey string
}

type AddEntity struct {
	Id       int32       // ID сущности
	Etype    string      // тип сущности: card, text, logopas, binary и т.д.
	Props    []*Property // массив значений свойств
	Metainfo []*Metainfo // массив значений метаинформации
}

type Property struct {
	EntityId int32  // код сущности
	FieldId  int32  // код описания поля свйоства
	Value    string // значение свойства
}

type Metainfo struct {
	EntityId int32  // код сущности
	Title    string // наименование метаинформации
	Value    string // значение метаинформации
}

type EntityCode struct {
	Etype string
	Name  string
}

type Field struct {
	Id               int32
	Name             string
	Ftype            string
	ValidateRules    string
	ValidateMessages string
}

// NewKeeperClient
func NewKeeperClient(serverAddress string, secretKey string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*KeeperClient, *grpc.ClientConn, error) {

	kc := &KeeperClient{
		secretKey: secretKey,
	}

	// перехватчики
	excludeMethods := map[string]bool{}
	authInterceptor := NewAuthInterceptor(kc, excludeMethods)

	// методы, данные в которых надо шифровать
	validOutCryptMethods := map[string]bool{constants.MethodAddEntity: true, constants.MethodEntity: true}
	dataOutInterceptor := NewDataOutInterceptor(kc, secretKey, validOutCryptMethods)

	opts = append(opts,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(authInterceptor.TokenInterceptor()),
		grpc.WithUnaryInterceptor(dataOutInterceptor.DataOutputInterceptor()))

	conn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		return nil, nil, err
	}

	kc.KeeperClient = pb.NewKeeperClient(conn)

	return kc, conn, nil
}

// Registration регистрация пользователя
// На входе: логин, пароль, повторный пароль
// Возвращает токен авторизации в случае успеха и ошибку
func (t *KeeperClient) Registration(ctx context.Context, login string, password string, password2 string) (string, error) {

	if password != password2 {
		return "", fmt.Errorf("пароли не совпадают")
	}

	res, err := t.KeeperClient.Registration(ctx, &pb.RegisterRequest{
		Login:          login,
		Password:       password,
		RepeatPassword: password2,
	})

	if err != nil {
		return "", err
	}

	if res.Error != "" {
		return "", fmt.Errorf(res.Error)
	}

	return res.Token, nil
}

func (t *KeeperClient) Login(ctx context.Context, login string, password string) (string, error) {

	lr, err := t.KeeperClient.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		return "", err
	}

	if lr.Error != "" {
		return "", fmt.Errorf(lr.Error)
	}

	t.token = lr.Token
	t.password = password

	return lr.Token, err
}

func (t *KeeperClient) EntityCodes(ctx context.Context) ([]*EntityCode, error) {
	var opts []grpc.CallOption

	in := &pb.EntityCodesRequest{Token: t.token}

	ec, err := t.KeeperClient.EntityCodes(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	var entcodes = make([]*EntityCode, 0, len(ec.EntityCodes))
	for _, val := range ec.EntityCodes {
		entcodes = append(entcodes, &EntityCode{
			Etype: val.Etype,
			Name:  val.Name,
		})
	}

	return entcodes, nil
}

func (t *KeeperClient) Fields(ctx context.Context, etype string) ([]*Field, error) {
	var opts []grpc.CallOption

	in := &pb.FieldsRequest{Etype: etype}

	resp, err := t.KeeperClient.Fields(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	var fd = make([]*Field, 0, len(resp.Fields))
	for _, val := range resp.Fields {
		fd = append(fd, &Field{
			Id:               val.Id,
			Name:             val.Name,
			Ftype:            val.Ftype,
			ValidateRules:    val.ValidateRules,
			ValidateMessages: val.ValidateMessages,
		})
	}

	return fd, nil
}

func (t *KeeperClient) AddEntity(ctx context.Context, ae AddEntity) (int32, error) {
	var opts []grpc.CallOption

	var props = make([]*pb.Property, 0, len(ae.Props))
	for _, val := range ae.Props {
		props = append(props, &pb.Property{
			EntityId: val.EntityId,
			FieldId:  val.FieldId,
			Value:    val.Value,
		})
	}

	var metainfo = make([]*pb.Metainfo, 0, len(ae.Metainfo))
	for _, val := range ae.Metainfo {
		metainfo = append(metainfo, &pb.Metainfo{
			EntityId: val.EntityId,
			Title:    val.Title,
			Value:    val.Value,
		})
	}

	in := &pb.AddEntityRequest{
		Id:       ae.Id,
		Etype:    ae.Etype,
		Props:    props,
		Metainfo: metainfo,
	}

	resp, err := t.KeeperClient.AddEntity(ctx, in, opts...)

	if resp.Error != "" {
		return 0, fmt.Errorf(resp.Error)
	}

	if err != nil {
		return 0, err
	}

	return resp.Id, err
}

func (t *KeeperClient) UploadBinary(ctx context.Context, entityId int32, file string) (int32, error) {
	stream, err := t.KeeperClient.UploadBinary(ctx)
	if err != nil {
		return 0, err
	}

	fil, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	// размер фрагментов передачи бинарных данных
	buf := make([]byte, constants.ChunkSize)

	for {
		num, err := fil.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		if err := stream.Send(&pb.UploadBinRequest{EntityId: entityId, ChunkData: buf[:num]}); err != nil {
			return 0, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return 0, err
	}

	if res.Error != "" {
		return 0, errors.New(res.Error)
	}

	return res.Size, nil

}

// DownloadBinary возвращает путь к загруженному файлу
func (t *KeeperClient) DownloadBinary(ctx context.Context, entityId int32, fileName string) (string, error) {

	stream, err := t.KeeperClient.DownloadBinary(ctx, &pb.DownloadBinRequest{EntityId: entityId})
	if err != nil {
		return "", err
	}

	wd, _ := os.Getwd()
	parts := strings.Split(wd, "internal")
	uploadDir := parts[0] + "cmd/client/" + constants.FileStorage
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	uploadFile := uploadDir + "/" + fmt.Sprintf("%v_", time.Now().Unix()) + fileName
	f, err := os.OpenFile(uploadFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var downloaded int64
	var buffer bytes.Buffer

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			err = f.Close()
			if err != nil {
				return "", err
			}
			break
		}
		if err != nil {
			buffer.Reset()
			return "", err
		}

		chunk := res.GetChunkData()
		shardSize := len(chunk)
		downloaded += int64(shardSize)

		_, err = f.Write(chunk)
	}

	return uploadFile, nil
}

/************************************ Шифрование бинарного потока ************************************/

func (t *KeeperClient) UploadCryptoBinary(ctx context.Context, entityId int32, file string) (int32, error) {
	stream, err := t.KeeperClient.UploadCryptoBinary(ctx)
	if err != nil {
		return 0, err
	}

	fil, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	// размер фрагментов передачи бинарных данных
	buf := make([]byte, constants.ChunkSize)

	for {
		num, err := fil.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		// шифруем
		cryptoKey := utils.SymmPassCreate(t.password, t.secretKey)
		crypted := utils.EncryptBinary(buf[:num], cryptoKey)

		if err := stream.Send(&pb.UploadBinRequest{EntityId: entityId, ChunkData: crypted}); err != nil {
			return 0, err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return 0, err
	}

	if res.Error != "" {
		return 0, errors.New(res.Error)
	}

	return res.Size, nil

}

func (t *KeeperClient) DownloadCryptoBinary(ctx context.Context, entityId int32, fileName string) (string, error) {
	cryptoKey := utils.SymmPassCreate(t.password, t.secretKey)

	stream, err := t.KeeperClient.DownloadCryptoBinary(ctx, &pb.DownloadBinRequest{EntityId: entityId})
	if err != nil {
		return "", err
	}

	wd, _ := os.Getwd()
	parts := strings.Split(wd, "internal")
	uploadDir := parts[0] + "cmd/client/" + constants.FileStorage
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	uploadFile := uploadDir + "/" + fmt.Sprintf("%v_", time.Now().Unix()) + fileName
	f, err := os.OpenFile(uploadFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var downloaded int64
	var buffer bytes.Buffer

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			err = f.Close()
			if err != nil {
				return "", err
			}
			break
		}
		if err != nil {
			buffer.Reset()
			return "", err
		}

		// расшифровка
		chunk := utils.DecryptBinary(res.GetChunkData(), cryptoKey)

		shardSize := len(chunk)
		downloaded += int64(shardSize)

		_, err = f.Write(chunk)
	}

	return uploadFile, nil
}

func (t *KeeperClient) GetToken() string {
	return t.token
}

func (t *KeeperClient) GetPassword() string {
	return t.password
}

func (t *KeeperClient) Start() {

	ex, err := os.Executable()
	if err != nil {
		logger.Log().Fatal(err.Error())
	}
	workDir := filepath.Dir(ex)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            workDir + "/logs/" + constants.LogReadline,
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt(">>> ")
			continue
		}
		cmd := strings.Join(cmds, "\n")
		cmds = cmds[:0]
		rl.SetPrompt("> ")
		rl.SaveHistory(cmd)
		println(cmd)
	}

}
