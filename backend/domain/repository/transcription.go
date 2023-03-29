package repository

import (
	"context"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
)

type ITranscriptionRepository interface {
	Read(ctx context.Context, id string) (transcription model.Transcription, err error)
	Save(ctx context.Context, transcription model.Transcription) error
}
