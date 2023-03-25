package pcm

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

type AudioSystem interface {
	Initialize() error
	OpenDefaultStream(numInputChannels int, numOutputChannels int, sampleRate float64, framesPerBuffer int, args ...interface{}) (AudioSystemStream, error)
	Terminate() error
}

type AudioSystemStream interface {
	Close() error
	Read() error
	Start() error
	Stop() error
	Time() time.Duration
}

type PortAudioSystem struct{}

func (p *PortAudioSystem) Initialize() error {
	return portaudio.Initialize()
}

func (p *PortAudioSystem) Terminate() error {
	return portaudio.Terminate()
}

func (p *PortAudioSystem) OpenDefaultStream(numInputChannels int, numOutputChannels int, sampleRate float64, framesPerBuffer int, args ...interface{}) (AudioSystemStream, error) {
	stream, err := portaudio.OpenDefaultStream(numInputChannels, numOutputChannels, sampleRate, framesPerBuffer, args...)
	return stream, err
}

type PCMRecorder struct {
	Interval             int
	SilentRatio          float32
	BaseLangCode         string
	AltLangCodes         []string
	BufferedContents     []int16
	recognitionStartTime time.Duration
	silentCount          int
	audioSystem          AudioSystem
}

func NewPCMRecorder(audioSystem AudioSystem, interval int) *PCMRecorder {
	var pr = &PCMRecorder{
		Interval:             interval,
		recognitionStartTime: -1,
		audioSystem:          audioSystem,
	}
	return pr
}

func (pr *PCMRecorder) Start(sig chan os.Signal, filepathCh chan string, wait *sync.WaitGroup) error {
	pr.audioSystem.Initialize()
	defer pr.audioSystem.Terminate()

	var err error
	stream, err := pr.initializeAudioStream()
	if err != nil {
		log.Fatalf("Could not open default stream \n %v", err)
	}
	(*stream).Start()
	defer (*stream).Close()

loop:
	for {
		select {
		case <-sig:
			wait.Done()
			close(filepathCh)
			break loop
		default:
		}

		pr.processAudioInput(*stream, filepathCh)
	}

	return nil
}

func (pr *PCMRecorder) initializeAudioStream() (*AudioSystemStream, error) {
	input := make([]int16, 64)
	stream, err := pr.audioSystem.OpenDefaultStream(1, 0, 44100, len(input), input)
	return &stream, err
}

func (pr *PCMRecorder) processAudioInput(stream AudioSystemStream, filePathCh chan string) {
	input := make([]int16, 64)

	if err := stream.Read(); err != nil {
		log.Fatalf("Could not read stream\n%v", err)
	}

	if !pr.detectSilence(input) {
		pr.record(input, stream.Time())
	} else {
		pr.silentCount++
	}

	if pr.detectSpeechStopped() || pr.detectSpeechExceededLimitation() {
		pr.finalizeRecording(filePathCh)
	}
}

func (pr *PCMRecorder) finalizeRecording(filepathCh chan string) {
	outputFileName := fmt.Sprintf("_%d.wav", int(pr.recognitionStartTime))
	fmt.Println(outputFileName)
	pr.writePCMData(outputFileName, pr.BufferedContents)
	filepathCh <- outputFileName

	pr.BufferedContents = nil
	pr.silentCount = 0
	pr.recognitionStartTime = -1
}

func (pr *PCMRecorder) record(input []int16, startTime time.Duration) {
	pr.silentCount = 0
	if pr.recognitionStartTime == -1 {
		pr.recognitionStartTime = startTime
	}
	pr.BufferedContents = append(pr.BufferedContents, input...)
}

func (pr *PCMRecorder) detectSilence(input []int16) bool {
	silent := true
	for _, bit := range input {
		// TODO: We should support ratio
		if bit != 0 {
			silent = false
			break
		}
	}
	return silent
}

func (pr *PCMRecorder) detectSpeechStopped() bool {
	return len(pr.BufferedContents) > 0 && pr.silentCount > 50
}

func (pr *PCMRecorder) detectSpeechExceededLimitation() bool {
	return len(pr.BufferedContents) >= (44100 * pr.Interval)
}

func (pr *PCMRecorder) writePCMData(outputFileName string, pcmData []int16) {
	if exists(outputFileName) {
		log.Fatalf("The audio file is already exists.")
	}
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Could not create a new file to write \n %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Could not close output file \n %v", err)
		}
	}()

	wav := NewWAVEncoder(file, pcmData)
	wav.Encode()
}

func exists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}
