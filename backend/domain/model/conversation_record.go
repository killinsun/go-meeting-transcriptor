package model

import "time"

type ConversationRecord struct {
	Id         string
	WavData    []byte
	Timestamp  time.Time
	repository IConversationRecordRepository
}

type IConversationRecordRepository interface {
	Read() (wavData []byte, err error)
	Save(id string, wavData []byte, timestamp time.Time) error
}

func NewConversationRecord(repository IConversationRecordRepository) *ConversationRecord {
	return &ConversationRecord{
		repository: repository,
	}
}

func (c *ConversationRecord) Read() (wavData []byte, err error) {
	wavData, err = c.repository.Read()
	return wavData, err
}

func (c *ConversationRecord) Save() error {
	return c.repository.Save(c.Id, c.WavData, c.Timestamp)
}
