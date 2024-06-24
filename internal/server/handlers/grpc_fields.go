package handlers

import (
	"context"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
)

// Fields запрос данных для добавления новой сущности (набор полей и их характеристик)
func (g *GRPCServer) Fields(ctx context.Context, in *pb.FieldsRequest) (*pb.FieldsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DBContextTimeout)
	defer cancel()

	fields, err := g.svs.FieldService.Fields(ctx, in.Etype)
	if err != nil {
		return nil, err
	}

	var f = make([]*pb.Field, 0, len(fields))
	for _, val := range fields {
		f = append(f, &pb.Field{
			Id:               val.ID,
			Name:             val.Name,
			Ftype:            val.Ftype,
			ValidateRules:    val.ValidateRules,
			ValidateMessages: val.ValidateMessages,
		})
	}

	return &pb.FieldsResponse{
		Fields: f,
	}, nil
}
