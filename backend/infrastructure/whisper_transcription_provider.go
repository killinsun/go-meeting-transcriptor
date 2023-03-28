package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/killinsun/go-meeting-transcriptor/backend/domain/model"
)

type WhisperResponse struct {
	Text string `json:"text"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type WhisperTranscriptionProvider struct {
	Client HTTPClient
}

func NewWhisperTranscriptionProvider(client HTTPClient) *WhisperTranscriptionProvider {
	return &WhisperTranscriptionProvider{
		Client: client,
	}
}

func (wt *WhisperTranscriptionProvider) Transcribe(wavData []byte) (model.Transcription, error) {
	transcription := model.Transcription{
		Text:      "",
		Timestamp: time.Now(),
	}

	err := godotenv.Load("../.env")
	authToken := os.Getenv("OPENAI_API_KEY")

	requestBody, contentType, err := wt.buildRequestBody(wavData)

	url := "https://api.openai.com/v1/audio/transcriptions"
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return transcription, err
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Authorization", "Bearer "+authToken)

	response, err := wt.Client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return transcription, err
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return transcription, err
	}

	var whisperResponse WhisperResponse
	if err := json.Unmarshal(responseBytes, &whisperResponse); err != nil {
		log.Fatal(err)
	}

	transcription.Text = whisperResponse.Text
	return transcription, nil
}

func (wt *WhisperTranscriptionProvider) buildRequestBody(wavData []byte) (requestBody bytes.Buffer, contentType string, err error) {
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", "chank.wav")
	if err != nil {
		fmt.Println("Failed to create form file:", err)
		return requestBody, "", err
	}
	part.Write(wavData)
	writer.WriteField("model", "whisper-1")
	writer.WriteField("response_format", "json")
	writer.Close()

	return requestBody, writer.FormDataContentType(), nil
}
