package infrastructure

import (
	"testing"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
	"golang.org/x/net/context"
)

func TestSave(t *testing.T) {
	meetingId := "meeting001"
	want := model.Transcription{Text: "Hello World!!"}
	ctx := context.Background()

	repository := NewRedisTranscriptionRepository(meetingId)

	err := repository.Save(ctx, want)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	got, _ := repository.Read(ctx)

	if want.Text != got.Text {
		t.Errorf("\n got %v, \nwant %v", got, want)
	}
}
