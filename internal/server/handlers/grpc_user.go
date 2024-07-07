package handlers

import (
	"context"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
)

// Registration регистрация нового пользователя
func (g *GRPCServer) Registration(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, constants.DBContextTimeout)
	defer cancel()

	token, err := g.svs.UserService.Registration(ctx, in.Login, in.Password, in.RepeatPassword)
	if err != nil {
		return &pb.RegisterResponse{
			Token: "",
			Error: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Token: token,
		Error: "",
	}, nil

}

// Login вход пользователя
func (g *GRPCServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DBContextTimeout)
	defer cancel()

	token, err := g.svs.UserService.Login(ctx, in.Login, in.Password)
	if err != nil {
		return &pb.LoginResponse{
			Token: "",
			Error: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Token: token,
		Error: "",
	}, nil
}
