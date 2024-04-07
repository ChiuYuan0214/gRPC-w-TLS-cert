package main

import (
	"context"
	"errors"
	"net"
	pb "practice-grpc/server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type SayHelloServer struct {
	pb.UnimplementedSayHelloServer
}

func (s *SayHelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("token is missing")
	}

	var appId string
	var appKey string
	if val, ok := md["appid"]; ok {
		appId = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appKey = val[0]
	}
	if appId != "adamchiu" || appKey != "123456" {
		return nil, errors.New("token invalid")
	}

	return &pb.HelloResponse{ResponseMsg: "hello " + req.RequestName}, nil
}

func main() {
	creds, _ := credentials.NewServerTLSFromFile("/Users/adam/Desktop/practice-grpc/key/test.pem", "/Users/adam/Desktop/practice-grpc/key/test.key")
	listener, _ := net.Listen("tcp", ":9090")

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterSayHelloServer(grpcServer, &SayHelloServer{})

	grpcServer.Serve(listener)
}
