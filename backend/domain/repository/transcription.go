package repository

import (
	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
)

type ITranscriptionRepository interface {
	Read() (transcription model.Transcription, err error)
	Save(transcription model.Transcription) error
}
