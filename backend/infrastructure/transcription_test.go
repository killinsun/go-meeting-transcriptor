package infrastructure

import (
	"fmt"
	"testing"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func TestSave(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	initializeRedisStorage(r)

	t.Run("Should return Transcription struct", func(t *testing.T) {
		ctx := context.Background()
		meetingId := "meeting001"
		key := fmt.Sprintf("Transcription::%v", meetingId)
		repository := NewRedisTranscriptionRepository(meetingId)

		want := model.Transcription{Text: "Hello World!!"}
		got, _, _ := r.SScan(ctx, key, 0, "", 100).Result()
		if 0 < len(got) {
			t.Errorf("\n got %v, \nwant %v", got, want)
		}

		err := repository.Save(ctx, want)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		got, _, _ = r.SScan(ctx, key, 0, "", 100).Result()
		if want.Text != got[1] {
			t.Errorf("\n got %v, \nwant %v", got[1], want.Text)
		}
	})

	t.Run("supports multiple transcriptions with the same meeting id", func(t *testing.T) {
		ctx := context.Background()
		meetingId := "meeting002"
		key := fmt.Sprintf("Transcription::%v", meetingId)
		repository := NewRedisTranscriptionRepository(meetingId)

		want := [...]model.Transcription{
			{Text: "This is the first text"},
			{Text: "This is the second text"},
		}

		repository.Save(ctx, want[0])
		repository.Save(ctx, want[1])

		got, _, _ := r.SScan(ctx, key, 0, "", 100).Result()
		for i := 1; i < len(want)+1; i++ {
			if got[i] != want[i-1].Text {
				t.Errorf("\n got %v, \nwant %v", got[i], want[i-1].Text)
			}
		}
	})
}

func initializeRedisStorage(r *redis.Client) {
	ctx := context.Background()

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.Scan(ctx, 0, "Transcription::*", 100).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			r.Del(ctx, key)
		}

		if cursor == 0 {
			break
		}
	}

}
