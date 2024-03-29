package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/killinsun/go-meeting-transcriptor/backend/infrastructure"
	transcriptorpb "github.com/killinsun/go-meeting-transcriptor/backend/pkg/grpc"
	"github.com/killinsun/go-meeting-transcriptor/backend/usecase"
)

func main() {
	port := 8080

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	transcriptorpb.RegisterTranscriptorServiceServer(s, NewTranscriptionServer())

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

type transcriptionServer struct {
	transcriptorpb.UnimplementedTranscriptorServiceServer
}

func NewTranscriptionServer() *transcriptionServer {
	return &transcriptionServer{}
}

func (t *transcriptionServer) StreamWav(stream transcriptorpb.TranscriptorService_StreamWavServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&transcriptorpb.WavResponse{
				Done: true,
			})
		}
		if err != nil {
			return err
		}
		ioutil.WriteFile("test.wav", req.GetData(), 0644)
		client := http.Client{}
		whisper := infrastructure.NewWhisperTranscriptionProvider(&client)
		repo := infrastructure.NewRedisTranscriptionRepository("test")
		transcriptionService := usecase.NewTranscriptionService(repo, whisper)
		transcriptionService.GetTranscription(req.GetData())
	}
}
