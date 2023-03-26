package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	client transcriptorpb.TranscriptorServiceClient
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
	client = transcriptorpb.NewTranscriptorServiceClient(conn)

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

	wavStream, err := client.StreamWav(context.Background())
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
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			if err := wavStream.Send(&transcriptorpb.WavChunk{Data: b}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	<-sig
	wait.Wait()

	time.Sleep(1 * time.Second)

	res, err := wavStream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetDone())
	}
	fmt.Println("Streaming finished.")
}
