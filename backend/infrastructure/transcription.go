package infrastructure

import (
	"fmt"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type RedisTranscriptionRepository struct {
	redisClient    *redis.Client
	redisKeyPrefix string
}

func NewRedisTranscriptionRepository(meetingId string) *RedisTranscriptionRepository {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &RedisTranscriptionRepository{
		redisKeyPrefix: fmt.Sprintf("Transcription:%v", meetingId),
		redisClient:    redisClient,
	}
}

func (t *RedisTranscriptionRepository) Read(ctx context.Context, id string) (transcription model.Transcription, err error) {
	key := fmt.Sprintf("%v:%v", t.redisKeyPrefix, id)
	result, err := t.redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return model.Transcription{}, err
	}
	transcription.Text = result

	return transcription, nil
}

func (t *RedisTranscriptionRepository) Save(ctx context.Context, transcription model.Transcription) (err error) {
	lastId, err := t.redisClient.Incr(ctx, t.redisKeyPrefix).Result()
	if err != nil {
		fmt.Println(err)
		return err
	}

	key := t.redisKeyPrefix + ":" + fmt.Sprintf("%v", lastId)
	text := transcription.Text
	err = t.redisClient.Set(ctx, key, text, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
