package entity

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/logger"
)

type EntityRepo interface {
	CreateEntity(ctx context.Context, entity EntityModel) (int32, error)
	GetEntity(ctx context.Context, id int32) (EntityModel, error)
	GetBinaryFilenameByEntityID(ctx context.Context, entityID int32) (string, error)
	SetChunkCountForCryptoBinary(ctx context.Context, entityID int32, chunkCount int32) error
	GetEntityListByType(ctx context.Context, etype string, userID int32) (map[int32][]string, error)
}

type FieldRepo interface {
	IsFieldType(ctx context.Context, id int32, ftype string) (bool, error)
}

type Property struct {
	ID       int32
	EntityID int32
	FieldID  int32
	Value    string
}

type Metainfo struct {
	ID       int32
	EntityID int32
	Title    string
	Value    string
}

type EntityModel struct {
	ID       int32
	UserID   int32
	Etype    string
	Props    []Property
	Metainfo []Metainfo
}

type Entity struct {
	repoEntity EntityRepo
	repoField  FieldRepo
}

// BinaryFileProperty Данные в поле свойства бинарной сущности содержат JSON в формате:
// {"servername": "имя файла на сервере (полный путь), "clientname": "только имя файла, под которым его грузили с клиента"}
type BinaryFileProperty struct {
	Servername string `json:"servername"`
	Clientname string `json:"clientname"`
	Chunkcount int32  `json:"chunkcount"`
}

func NewEntity(repoEntity EntityRepo, repoField FieldRepo) (*Entity, error) {
	e := &Entity{
		repoEntity: repoEntity,
		repoField:  repoField,
	}

	return e, nil
}

func (e *Entity) AddEntity(ctx context.Context, entity EntityModel) (int32, error) {

	// если среди добавляемых свойств есть ftype=path (означает что данные должны быть сохранены в файле)
	// создаем файл
	// ../filebank/<etype>/<код_пользователя>/<случайная_строка_как_имя_директории>/<chunk_index+эта_же_случайная_строка_как_имя_файла> и запоминаем путь к файлу как свойство
	for i, val := range entity.Props {
		isType, _ := e.repoField.IsFieldType(ctx, val.FieldID, constants.FieldTypePath)
		if !isType {
			continue
		}

		path, _ := os.Executable()
		useridStr := fmt.Sprintf("%v", entity.UserID)

		c := 10
		b := make([]byte, c)
		rand.Read(b)
		randName := hex.EncodeToString(b)

		fileBankDir := filepath.Dir(path) + "/" + constants.FileBankDir + "/" + entity.Etype + "/" + useridStr + "/" + randName
		err := os.MkdirAll(fileBankDir, os.ModePerm)
		if err != nil {
			return 0, err
		}
		fileBankPath := fileBankDir + "/" + randName
		f, err := os.Create(fileBankPath)
		if err != nil {
			return 0, err
		}
		f.Close()

		p := BinaryFileProperty{
			Servername: fileBankPath,
			Clientname: val.Value,
		}
		propval, _ := json.Marshal(p)

		entity.Props[i].Value = string(propval)
	}

	id, err := e.repoEntity.CreateEntity(ctx, entity)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *Entity) Entity(ctx context.Context, id int32) (*EntityModel, error) {

	ent, err := e.repoEntity.GetEntity(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (e *Entity) UploadBinary(stream pb.Keeper_UploadBinaryServer) (int32, error) {

	var uploadSize int32
	var f *os.File
	defer f.Close()

	filename := ""

	for {
		req, err := stream.Recv()

		// получаем из базы путь к файлу для сохранения
		if filename == "" {
			p, err := e.repoEntity.GetBinaryFilenameByEntityID(context.Background(), req.EntityId)
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}

			binprop := &BinaryFileProperty{}
			err = json.Unmarshal([]byte(p), binprop)
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}

			filename = binprop.Servername
			f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}
		}

		if err == io.EOF {
			err = f.Close()
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}

			err = stream.SendAndClose(&pb.UploadBinResponse{
				Size:  uploadSize,
				Error: "",
			})
			if err != nil {
				return uploadSize, err
			}
			return uploadSize, nil
		}
		if err != nil {
			return uploadSize, status.Error(codes.Internal, err.Error())
		}

		_, err = f.Write(req.GetChunkData())
		if err != nil {
			return uploadSize, status.Error(codes.Internal, err.Error())
		}

		uploadSize = uploadSize + int32(len(req.ChunkData))
	}

}

func (e *Entity) DownloadBinary(entityID int32, stream pb.Keeper_DownloadBinaryServer) error {
	ctx := context.Background()

	filedata, err := e.repoEntity.GetBinaryFilenameByEntityID(ctx, entityID)
	if err != nil {
		return err
	}

	fd := &BinaryFileProperty{}
	err = json.Unmarshal([]byte(filedata), fd)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(fd.Servername)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	f, err := os.Open(fd.Servername)
	if err != nil {
		return err
	}
	defer f.Close()

	var totalBytesStreamed int64

	for totalBytesStreamed < fileSize {
		chunk := make([]byte, constants.ChunkSize)
		bytesRead, err := f.Read(chunk)
		if err == io.EOF {
			logger.Log().Info(fmt.Sprintf("download complete: %v", fd.Servername))
			break
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&pb.DownloadBinResponse{ChunkData: chunk[:bytesRead]}); err != nil {
			return err
		}
		totalBytesStreamed += int64(bytesRead)
	}

	return nil
}

/************************************ Зашифрованные бинарные фрагменты  *************************************/

func (e *Entity) UploadCryptoBinary(stream pb.Keeper_UploadCryptoBinaryServer) (int32, error) {

	var uploadSize int32
	var entityID int32 = 0

	dirBase, fileBase := "", ""
	var index int32 = 0

	for {
		req, err := stream.Recv()

		// получаем из базы путь к файлу для сохранения
		if fileBase == "" {
			p, err := e.repoEntity.GetBinaryFilenameByEntityID(context.Background(), req.EntityId)
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}
			entityID = req.EntityId

			binprop := &BinaryFileProperty{}
			err = json.Unmarshal([]byte(p), binprop)
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}

			dirBase = path.Dir(binprop.Servername)
			fileBase = path.Base(binprop.Servername)

		}

		// у каждого фрагмента отдельный файл с префиксом-индексом перед именем файла
		index++
		indexStr := fmt.Sprintf("%06d", index)
		filename := dirBase + "/" + indexStr + "_" + fileBase
		f, errC := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if errC != nil {
			return 0, status.Error(codes.Internal, errC.Error())
		}

		if err == io.EOF {
			_ = f.Close()
			err = os.Remove(filename)
			if err != nil {
				return 0, status.Error(codes.Internal, err.Error())
			}
			index--

			err = stream.SendAndClose(&pb.UploadBinResponse{
				Size:  uploadSize,
				Error: "",
			})
			if err != nil {
				return uploadSize, err
			}

			// успешное завершение
			// сохраняем кол-во файлов-фрагментов в свойство сущности
			err = e.repoEntity.SetChunkCountForCryptoBinary(context.Background(), entityID, index)
			if err != nil {
				return uploadSize, err
			}

			return uploadSize, nil
		}
		if err != nil {
			return uploadSize, status.Error(codes.Internal, err.Error())
		}

		_, err = f.Write(req.GetChunkData())
		if err != nil {
			return uploadSize, status.Error(codes.Internal, err.Error())
		}

		uploadSize = uploadSize + int32(len(req.ChunkData))

		f.Close()
	}

}

func (e *Entity) DownloadCryptoBinary(entityID int32, stream pb.Keeper_DownloadCryptoBinaryServer) error {
	ctx := context.Background()

	filedata, err := e.repoEntity.GetBinaryFilenameByEntityID(ctx, entityID)
	if err != nil {
		return err
	}

	fd := &BinaryFileProperty{}
	err = json.Unmarshal([]byte(filedata), fd)
	if err != nil {
		return err
	}

	filesDir := path.Dir(fd.Servername)
	fileBase := path.Base(fd.Servername)
	fileCount := fd.Chunkcount

	for index := int32(1); index <= fileCount; index++ {

		indexStr := fmt.Sprintf("%06d", index)
		filename := filesDir + "/" + indexStr + "_" + fileBase
		chunk, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.DownloadBinResponse{ChunkData: chunk}); err != nil {
			return err
		}

	}

	return nil
}

// EntityList Получение списка сущностей указанного типа для конкретного пользователя
// Простая карта с кодом сущности и названием(составляется из метаданных)
func (e *Entity) EntityList(ctx context.Context, etype string, userID int32) (map[int32]string, error) {

	list, err := e.repoEntity.GetEntityListByType(ctx, etype, userID)

	data := make(map[int32]string, len(list))
	for key, val := range list {
		a := make(map[string]string)
		for _, v := range val {
			parts := strings.Split(v, ":")
			a[parts[0]] = parts[1]
		}
		meta, _ := json.Marshal(a)
		data[key] = string(meta)
	}

	return data, err
}
