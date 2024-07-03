package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	pb "github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
	"github.com/dnsoftware/gophkeeper/internal/utils"
	"github.com/dnsoftware/gophkeeper/logger"
)

func (g *GRPCServer) AddEntity(ctx context.Context, in *pb.AddEntityRequest) (*pb.AddEntityResponse, error) {

	userID := getContextUserID(ctx)

	var props = make([]entity.Property, 0, len(in.Props))
	for _, val := range in.Props {
		props = append(props, entity.Property{
			EntityID: val.EntityId,
			FieldID:  val.FieldId,
			Value:    val.Value,
		})
	}

	var metainfo = make([]entity.Metainfo, 0, len(in.Metainfo))
	for _, val := range in.Metainfo {
		metainfo = append(metainfo, entity.Metainfo{
			EntityID: val.EntityId,
			Title:    val.Title,
			Value:    val.Value,
		})
	}

	ent := entity.EntityModel{
		UserID:   int32(userID),
		Etype:    in.Etype,
		Props:    props,
		Metainfo: metainfo,
	}

	id, err := g.svs.EntityService.AddEntity(ctx, ent)
	if err != nil {
		return nil, err
	}

	return &pb.AddEntityResponse{
		Id:    id,
		Error: "",
	}, nil
}

func (g *GRPCServer) SaveEditEntity(ctx context.Context, in *pb.SaveEntityRequest) (*pb.SaveEntityResponse, error) {

	userID := getContextUserID(ctx)

	var props = make([]entity.Property, 0, len(in.Props))
	for _, val := range in.Props {
		props = append(props, entity.Property{
			EntityID: val.EntityId,
			FieldID:  val.FieldId,
			Value:    val.Value,
		})
	}

	var metainfo = make([]entity.Metainfo, 0, len(in.Metainfo))
	for _, val := range in.Metainfo {
		metainfo = append(metainfo, entity.Metainfo{
			EntityID: val.EntityId,
			Title:    val.Title,
			Value:    val.Value,
		})
	}

	ent := entity.EntityModel{
		ID:       in.Id,
		UserID:   int32(userID),
		Etype:    in.Etype,
		Props:    props,
		Metainfo: metainfo,
	}

	err := g.svs.EntityService.SaveEditEntity(ctx, ent)
	if err != nil {
		return nil, err
	}

	return &pb.SaveEntityResponse{
		Id:    ent.ID,
		Error: "",
	}, nil
}

func (g *GRPCServer) Entity(ctx context.Context, in *pb.EntityRequest) (*pb.EntityResponse, error) {

	ent, err := g.svs.EntityService.Entity(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var props = make([]*pb.Property, 0, len(ent.Props))
	for _, val := range ent.Props {
		props = append(props, &pb.Property{
			EntityId: val.EntityID,
			FieldId:  val.FieldID,
			Value:    val.Value,
		})
	}

	var metainfo = make([]*pb.Metainfo, 0, len(ent.Metainfo))
	for _, val := range ent.Metainfo {
		metainfo = append(metainfo, &pb.Metainfo{
			EntityId: val.EntityID,
			Title:    val.Title,
			Value:    val.Value,
		})
	}

	ret := &pb.EntityResponse{
		Id:       ent.ID,
		Etype:    ent.Etype,
		Props:    props,
		Metainfo: metainfo,
	}

	return ret, err
}

func (g *GRPCServer) DeleteEntity(ctx context.Context, in *pb.DeleteEntityRequest) (*pb.DeleteEntityResponse, error) {

	userID := getContextUserID(ctx)

	err := g.svs.EntityService.DeleteEntity(ctx, in.Id, int32(userID))
	if err != nil {
		return nil, err
	}

	return &pb.DeleteEntityResponse{Error: ""}, nil
}

func (g *GRPCServer) UploadBinary(stream pb.Keeper_UploadBinaryServer) error {

	size, err := g.svs.EntityService.UploadBinary(stream)
	if err != nil {
		return err
	}

	logger.Log().Info(fmt.Sprintf("Загружен бинарный файл размером %v", size))

	return nil
}

func (g *GRPCServer) DownloadBinary(in *pb.DownloadBinRequest, stream pb.Keeper_DownloadBinaryServer) error {

	err := g.svs.EntityService.DownloadBinary(in.EntityId, stream)
	if err != nil {
		return err
	}

	return nil
}

func (g *GRPCServer) UploadCryptoBinary(stream pb.Keeper_UploadCryptoBinaryServer) error {

	size, err := g.svs.EntityService.UploadCryptoBinary(stream)
	if err != nil {
		return err
	}

	logger.Log().Info(fmt.Sprintf("Загружен бинарный файл размером %v", size))

	return nil
}

func (g *GRPCServer) DownloadCryptoBinary(in *pb.DownloadBinRequest, stream pb.Keeper_DownloadCryptoBinaryServer) error {

	err := g.svs.EntityService.DownloadCryptoBinary(in.EntityId, stream)
	if err != nil {
		return err
	}

	return nil
}

// EntityList Получение списка сущностей указанного типа для конкретного пользователя
// Простая карта с кодом сущности и названием(составляется из метаданных)
func (g *GRPCServer) EntityList(ctx context.Context, in *pb.EntityListRequest) (*pb.EntityListResponse, error) {
	userID := getContextUserID(ctx)

	list, err := g.svs.EntityService.EntityList(ctx, in.Etype, int32(userID))
	if err != nil {
		return nil, err
	}

	return &pb.EntityListResponse{
		List: list,
	}, nil
}

func getContextUserID(ctx context.Context) int {
	var token string
	var userID int

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			token = values[0]
			userID = utils.GetUserID(token)
		}
	}

	return userID
}
