package usecase

import (
	"fmt"
	"net/http"

	"github.com/killinsun/go-meeting-transcriptor/backend/infrastructure"
)

func GetTranscription(wavChank []byte) {
	client := &http.Client{}
	whisper := infrastructure.NewWhisperTranscriptionService(client)

	transcription, err := whisper.Transcribe(wavChank)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Transcription: %v", transcription)
}
