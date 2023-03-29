package infrastructure

import (
	"fmt"
	"testing"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func TestRead(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	initializeRedisStorage(r)

	t.Run("Should return Transcription struct", func(t *testing.T) {
		ctx := context.Background()
		meetingId := "meeting001"
		key := fmt.Sprintf("Transcription:%v:1", meetingId)
		want := model.Transcription{Text: "Hello"}
		repository := NewRedisTranscriptionRepository(meetingId)

		err := r.Set(ctx, key, want.Text, 0).Err()
		assertNil(err, t)

		got, err := repository.Read(ctx, "1")
		assertNil(err, t)
		assert(got.Text, want.Text, t)
	})
}

func TestSave(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	initializeRedisStorage(r)

	t.Run("Should add a new transcription", func(t *testing.T) {
		ctx := context.Background()
		meetingId := "meeting001"
		lastId := "1"
		key := fmt.Sprintf("Transcription:%v:%v", meetingId, lastId)
		repository := NewRedisTranscriptionRepository(meetingId)

		want := model.Transcription{Text: "Hello World!!"}
		got, _ := r.Get(ctx, key).Result()
		assert(got, "", t)

		err := repository.Save(ctx, want)
		assertNil(err, t)

		got, _ = r.Get(ctx, key).Result()
		assert(got, want.Text, t)
	})

	t.Run("supports multiple transcriptions with the same meeting id", func(t *testing.T) {
		ctx := context.Background()
		meetingId := "meeting002"
		repository := NewRedisTranscriptionRepository(meetingId)

		want := [...]model.Transcription{
			{Text: "This is the first text"},
			{Text: "This is the second text"},
		}

		repository.Save(ctx, want[0])
		repository.Save(ctx, want[1])

		for i := 1; i < len(want)+1; i++ {
			key := fmt.Sprintf("Transcription:%v:%v", meetingId, i)
			got, _ := r.Get(ctx, key).Result()
			assert(got, want[i-1].Text, t)
		}
	})
}

func initializeRedisStorage(r *redis.Client) {
	ctx := context.Background()

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.Scan(ctx, 0, "Transcription:*", 100).Result()
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

func assert(got, want interface{}, t *testing.T) {
	if got != want {
		t.Errorf("\n got %v, \nwant %v", got, want)
	}
}

func assertNil(got interface{}, t *testing.T) {
	if got != nil {
		t.Errorf("\n got %v", got)
	}
}
