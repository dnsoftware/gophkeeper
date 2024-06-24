package domain

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

	"github.com/dnsoftware/gophkeeper/internal/client/infrastructure"
	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/utils"
	"github.com/dnsoftware/gophkeeper/logger"
)

type KeeperClient struct {
	//cfg *config.ClientConfig
	pb.KeeperClient
	token     string
	password  string
	secretKey string
}

func NewKeeperClient(serverAddress string, secretKey string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*KeeperClient, *grpc.ClientConn, error) {

	kc := &KeeperClient{
		secretKey: secretKey,
	}

	// перехватчики
	excludeMethods := map[string]bool{}
	authInterceptor := infrastructure.NewAuthInterceptor(kc, excludeMethods)

	// методы, данные в которых надо шифровать
	validOutCryptMethods := map[string]bool{constants.MethodAddEntity: true, constants.MethodEntity: true}
	dataOutInterceptor := infrastructure.NewDataOutInterceptor(kc, secretKey, validOutCryptMethods)

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

func (t *KeeperClient) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
	lr, err := t.KeeperClient.Login(ctx, in, opts...)
	t.token = lr.Token
	t.password = in.Password

	return lr, err
}

func (t *KeeperClient) AddEntity(ctx context.Context, in *pb.AddEntityRequest, opts ...grpc.CallOption) (*pb.AddEntityResponse, error) {
	resp, err := t.KeeperClient.AddEntity(ctx, in, opts...)

	return resp, err
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
