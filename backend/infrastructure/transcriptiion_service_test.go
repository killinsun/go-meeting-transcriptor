package service

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
)

func TestTranscribe(t *testing.T) {
	// Mock HTTP Client Response
	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {

			if req.Header.Get("Authorization") == "Bearer " || req.Header.Get("Authorization") == "" {
				return nil, errors.New("Authorization header is empty")
			}

			response := &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"text": "Hello World!"}`))),
			}
			return response, nil
		},
	}
	t.Run("Should return Transcription struct", func(t *testing.T) {
		want := model.Transcription{Text: "Hello World!"}

		service := WhisperTranscriptionService{
			Client: mockHTTPClient,
		}

		got, _ := service.Transcribe([]byte("test"))

		if want.Text != got.Text {
			t.Errorf("\n got %v, \nwant %v", got, want)
		}
	})
}

type MockWhisperTranscriptionService struct{}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (mh *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return mh.DoFunc(req)
}
