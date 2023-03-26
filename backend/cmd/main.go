package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

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

func (g *gRPCServer) HelloServerStream(req *transcriptorpb.HelloRequest, stream transcriptorpb.GreetingService_HelloServerStreamServer) error {
	resCount := 5

	for i := 0; i < resCount; i++ {
		if err := stream.Send(&transcriptorpb.HelloResponse{
			Message: fmt.Sprintf("[%d]Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (g *gRPCServer) HelloClientStream(stream transcriptorpb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&transcriptorpb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}
