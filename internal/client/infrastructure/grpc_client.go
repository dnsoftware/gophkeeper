package infrastructure

import (
	pb "github.com/dnsoftware/gophkeeper/internal/proto"
)

type GRPCClient struct {
	cl pb.KeeperServer
}

func NewGRPCClient() {

	//creds, err := clientTLSCreds()
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//// перехватчики клиента
	//var opts []grpc.DialOption
	//opts = append(opts, grpc.WithUnaryInterceptor(addTokenInterceptor))
	//
	//conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(creds))
	//if err != nil {
	//	return nil, nil, err
	//}
	//client := pb.NewKeeperClient(conn)
	//
	//
	//client := pb.NewKeeperClient()
}
