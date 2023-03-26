package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	transcriptorpb "github.com/killinsun/go-meeting-transcriptor/backend/pkg/grpc"
)

func main() {
	port := 8080

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	transcriptorpb.RegisterGreetingServiceServer(s, NewGRPCServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server on port %d", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Stopping server...")
	s.GracefulStop()
}

type gRPCServer struct {
	transcriptorpb.UnimplementedGreetingServiceServer
}

func NewGRPCServer() *gRPCServer {
	return &gRPCServer{}
}

func (g *gRPCServer) Hello(ctx context.Context, req *transcriptorpb.HelloRequest) (*transcriptorpb.HelloResponse, error) {
	return &transcriptorpb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}