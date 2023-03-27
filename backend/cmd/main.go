package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/signal"

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
		GetTranscription(req.GetData())
	}
}

func GetTranscription(wavChank []byte) {
	authToken := "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"

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", "chank.wav")
	if err != nil {
		fmt.Println("Failed to create form file:", err)
		return
	}
	part.Write(wavChank)
	writer.WriteField("model", "whisper-1")
	writer.Close()

	url := "https://api.openai.com/v1/audio/transcriptions"
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	// レスポンスを出力する
	fmt.Println(string(responseBytes))
}
