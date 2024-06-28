package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/dnsoftware/gophkeeper/internal/client/domain"
	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/utils"
	"github.com/dnsoftware/gophkeeper/logger"
)

type GRPCSender struct {
	pb.KeeperClient
	token     string
	password  string
	SecretKey string
}

// NewGRPCSender обмен данными с сервером
func NewGRPCSender(serverAddress string, secretKey string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*GRPCSender, *grpc.ClientConn, error) {

	kc := &GRPCSender{
		SecretKey: secretKey,
	}

	// перехватчики
	excludeMethods := map[string]bool{}
	authInterceptor := NewAuthInterceptor(kc, excludeMethods)

	// методы, данные в которых надо шифровать
	validOutCryptMethods := map[string]bool{constants.MethodAddEntity: true, constants.MethodEntity: true}
	dataOutInterceptor := NewDataOutInterceptor(kc, secretKey, validOutCryptMethods)

	opts = append(opts,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(authInterceptor.TokenInterceptor(), dataOutInterceptor.DataOutputInterceptor()))

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
func (t *GRPCSender) Registration(login string, password string, password2 string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

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

	t.token = res.Token
	t.password = password

	return res.Token, nil
}

func (t *GRPCSender) Login(login string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

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

	return lr.Token, nil
}

func (t *GRPCSender) EntityCodes() ([]*domain.EntityCode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

	var opts []grpc.CallOption

	in := &pb.EntityCodesRequest{}

	ec, err := t.KeeperClient.EntityCodes(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	var entcodes = make([]*domain.EntityCode, 0, len(ec.EntityCodes))
	for _, val := range ec.EntityCodes {
		entcodes = append(entcodes, &domain.EntityCode{
			Etype: val.Etype,
			Name:  val.Name,
		})
	}

	return entcodes, nil
}

func (t *GRPCSender) Fields(etype string) ([]*domain.Field, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

	var opts []grpc.CallOption
	in := &pb.FieldsRequest{Etype: etype}

	resp, err := t.KeeperClient.Fields(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	var fd = make([]*domain.Field, 0, len(resp.Fields))
	for _, val := range resp.Fields {
		fd = append(fd, &domain.Field{
			Id:               val.Id,
			Etype:            val.Etype,
			Name:             val.Name,
			Ftype:            val.Ftype,
			ValidateRules:    val.ValidateRules,
			ValidateMessages: val.ValidateMessages,
		})
	}

	return fd, nil
}

func (t *GRPCSender) AddEntity(ae domain.Entity) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

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

func (t *GRPCSender) UploadBinary(entityId int32, file string) (int32, error) {

	stream, err := t.KeeperClient.UploadBinary(context.Background())
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
func (t *GRPCSender) DownloadBinary(entityId int32, fileName string) (string, error) {

	stream, err := t.KeeperClient.DownloadBinary(context.Background(), &pb.DownloadBinRequest{EntityId: entityId})
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

func (t *GRPCSender) UploadCryptoBinary(entityId int32, file string) (int32, error) {
	stream, err := t.KeeperClient.UploadCryptoBinary(context.Background())
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
		cryptoKey := utils.SymmPassCreate(t.password, t.SecretKey)
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

func (t *GRPCSender) DownloadCryptoBinary(entityId int32, fileName string) (string, error) {
	cryptoKey := utils.SymmPassCreate(t.password, t.SecretKey)

	stream, err := t.KeeperClient.DownloadCryptoBinary(context.Background(), &pb.DownloadBinRequest{EntityId: entityId})
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

func (t *GRPCSender) Entity(id int32) (*domain.Entity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

	resp, err := t.KeeperClient.Entity(ctx, &pb.EntityRequest{Id: id})
	if err != nil {
		return nil, err
	}

	userID := utils.GetUserID(t.token)

	props := make([]*domain.Property, 0, len(resp.Props))
	meta := make([]*domain.Metainfo, 0, len(resp.Metainfo))

	for _, val := range resp.Props {
		props = append(props, &domain.Property{
			EntityId: val.EntityId,
			FieldId:  val.FieldId,
			Value:    val.Value,
		})
	}

	for _, val := range resp.Metainfo {
		meta = append(meta, &domain.Metainfo{
			EntityId: val.EntityId,
			Title:    val.Title,
			Value:    val.Value,
		})
	}

	ent := &domain.Entity{
		Id:       id,
		UserID:   int32(userID),
		Etype:    resp.Etype,
		Props:    props,
		Metainfo: meta,
	}

	return ent, nil
}

func (t *GRPCSender) GetToken() string {
	return t.token
}

func (t *GRPCSender) GetPassword() string {
	return t.password
}

func (t *GRPCSender) Start() {

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

func (t *GRPCSender) EntityList(etype string) (map[int32]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBContextTimeout)
	defer cancel()

	var opts []grpc.CallOption
	in := &pb.EntityListRequest{Etype: etype}

	resp, err := t.KeeperClient.EntityList(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	cryptoKey := utils.SymmPassCreate(t.password, t.SecretKey)

	list := make(map[int32]string, len(resp.List))
	for key, val := range resp.List {
		m := make(map[string]string)
		err := json.Unmarshal([]byte(val), &m)
		if err != nil {
			return nil, err
		}

		str := ""
		for ek, ev := range m {
			k := utils.Decrypt(ek, cryptoKey)
			v := utils.Decrypt(ev, cryptoKey)
			str = str + k + ":" + v + ". "
		}

		list[key] = str
	}

	return list, nil
}
