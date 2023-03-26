package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	transcriptorpb "github.com/killinsun/go-meeting-transcriptor/backend/pkg/grpc"
	pcm "github.com/killinsun/go-meeting-transcriptor/node/domain/recorder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client transcriptorpb.GreetingServiceClient
)

func main() {
	fmt.Println("Streaming. Press Ctrl + C to stop.")

	baseDir := time.Now().Format("audio_20060102_T150405")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		panic("Could not create a new directory")
	}

	audioSystem := &pcm.PortAudioSystem{}
	pr := pcm.NewPCMRecorder(audioSystem, fmt.Sprintf(baseDir+"/file"), 30)

	pr.GetDeviceInfo()

	conn, err := grpc.Dial(
		"localhost:8080",

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()
	client = transcriptorpb.NewGreetingServiceClient(conn)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	filePathCh := make(chan string)

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		if err := pr.Start(sig, filePathCh, &wait); err != nil {
			log.Fatalf("Error starting PCMRecorder: %v", err)
		}
	}()

	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			filePath, ok := <-filePathCh
			if !ok {
				break
			}

			if err := stream.Send(&transcriptorpb.HelloRequest{Name: filePath}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	<-sig
	wait.Wait()

	time.Sleep(1 * time.Second)

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}

	HelloServerStream("Killinsun")
	fmt.Println("Streaming finished.")
}

func Hello(name string) {
	req := &transcriptorpb.HelloRequest{
		Name: name,
	}

	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream(name string) {
	req := &transcriptorpb.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("All the responses have already received.")
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res)
	}
}
