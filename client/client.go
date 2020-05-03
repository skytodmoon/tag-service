package main

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	"github.com/go-programming-tour-book/tag-service/global"

	"github.com/go-programming-tour-book/tag-service/internal/middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/go-programming-tour-book/tag-service/pkg/tracer"

	pb "github.com/go-programming-tour-book/tag-service/proto"
	"google.golang.org/grpc"
)

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

func main() {
	ctx := context.Background()
	newCtx := metadata.AppendToOutgoingContext(ctx, "eddycjy", "Go编程之旅")
	clientConn, err := GetClientConn(newCtx, "localhost:8005", []grpc.DialOption{grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryContextTimeout(),
			middleware.ClientTracing(),
		),
	)})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer clientConn.Close()
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	//newCtx := tagServiceClient.WithOrgCode(ctx, "Go编程之旅")
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}

func setupTracer() error {
	var err error
	jaegerTracer, _, err := tracer.NewJaegerTracer("article-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
