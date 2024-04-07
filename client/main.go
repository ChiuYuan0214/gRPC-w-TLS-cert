package main

import (
	"context"
	"log"

	pb "practice-grpc/client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// package "credentials"
//
// type PerRPCCredentials interface {
// 	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
// 	RequireTransportSecurity() bool
// }

type ClientTokenAuth struct{}

func (cta ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "adamchiu",
		"appKey": "123456",
	}, nil
}

func (cta ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {

	// 實際prod時，domain name會從browser動態抓取
	creds, _ := credentials.NewClientTLSFromFile("/Users/adam/Desktop/practice-grpc/key/test.pem", "*.adamchiu.com")
	// insecure.NewCredentials() = 不使用TLS
	// conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(creds))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(new(ClientTokenAuth))}
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSayHelloClient(conn)
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "Adam"})
	log.Println(resp.GetResponseMsg())
}
