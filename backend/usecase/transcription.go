package usecase

import (
	"fmt"
	"net/http"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/repository"
	"github.com/killinsun/go-meeting-transcriptor/backend/infrastructure"
)

type TranscriptionService struct {
	repo       repository.ITranscriptionRepository
	provider   infrastructure.ITranscriptionProvider
	httpclient *http.Client
}

func NewTranscriptionService(repo repository.ITranscriptionRepository, provider infrastructure.ITranscriptionProvider) *TranscriptionService {
	return &TranscriptionService{
		repo:     repo,
		provider: provider,
	}
}

func (t *TranscriptionService) GetTranscription(wavChank []byte) {
	transcription, err := t.provider.Transcribe(wavChank)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Transcription: %v", transcription)
}
