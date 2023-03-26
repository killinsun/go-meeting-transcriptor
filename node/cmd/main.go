package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	pcm "github.com/killinsun/go-meeting-transcriptor/node/domain/recorder"
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

	go func() {
		for {
			filePath, ok := <-filePathCh
			if !ok {
				break
			}
			fmt.Printf("Recorded file: %s\n", filePath)
		}
	}()

	<-sig
	wait.Wait()

	time.Sleep(1 * time.Second)
	fmt.Println("Streaming finished.")
}
