package model

import (
	"reflect"
	"testing"
	"time"
)

func TestConversationRecord_Read(t *testing.T) {
	t.Run("Should call Read function", func(t *testing.T) {
		want := []byte("test")
		repo := &MockConversationRecordRepository{}

		cr := NewConversationRecord(repo)
		got, _ := cr.Read()

		if !reflect.DeepEqual(want, got) || repo.readCalled != 1 {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestConversationRecord_Save(t *testing.T) {
	t.Run("Should call Save function", func(t *testing.T) {
		repo := &MockConversationRecordRepository{}

		cr := NewConversationRecord(repo)
		err := cr.Save()

		if err != nil || repo.saveCalled != 1 {
			t.Errorf("got %v", err)
		}
	})
}

type MockConversationRecordRepository struct {
	saveCalled int
	readCalled int
}

func (m *MockConversationRecordRepository) Read() ([]byte, error) {
	m.readCalled++
	return []byte("test"), nil
}

func (m *MockConversationRecordRepository) Save(id string, wavData []byte, timestamp time.Time) error {
	m.saveCalled++
	return nil
}
