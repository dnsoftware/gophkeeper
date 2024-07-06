package handlers

import (
	"context"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
)

func (g *GRPCServer) EntityCodes(ctx context.Context, _ *pb.EntityCodesRequest) (*pb.EntityCodesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DBContextTimeout)
	defer cancel()

	codes, err := g.svs.EntityCodeService.EntityCodes(ctx)
	if err != nil {
		return nil, err
	}

	var ec []*pb.EntityCode
	for etype, name := range codes {
		ec = append(ec, &pb.EntityCode{
			Etype: etype,
			Name:  name,
		})
	}

	return &pb.EntityCodesResponse{EntityCodes: ec}, nil
}
