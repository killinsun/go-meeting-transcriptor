package infrastructure

import "github.com/killinsun/go-meeting-transcriptor/backend/domain/model"

type ITranscriptionProvider interface {
	Transcribe(wavData []byte) (model.Transcription, error)
}
