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

	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			token = values[0]
		}
	}
	userID := utils.GetUserID(token)

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

func (g *GRPCServer) Entity(ctx context.Context, in *pb.EntityRequest) (*pb.EntityResponse, error) {

	ent, err := g.svs.EntityService.Entity(ctx, in.Id)

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
