package pcm

import (
	"testing"
)

func TestDetectSilence(t *testing.T) {
	t.Run("All array items are 0", func(t *testing.T) {
		input := []int16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		interval := 3
		pr := NewPCMRecorder(interval)

		got := pr.detectSilence(input)
		want := true

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Some voices are streamed", func(t *testing.T) {
		input := []int16{0, 0, 0, 120, 120, 44, 66, 10, -12, 0, 0, 0, 0, 0, 0, 0}

		interval := 3
		pr := NewPCMRecorder(interval)

		got := pr.detectSilence(input)
		want := false

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestDetectSpeechStopped(t *testing.T) {
	t.Run("Should return true when speech is stopped", func(t *testing.T) {
		interval := 3
		pr := NewPCMRecorder(interval)
		want := true

		contents := make([]int16, 64)
		// silece after some speech should be recognized as 'stop'
		for i := 0; i < 10; i++ {
			contents[i] = 1
		}
		pr.BufferedContents = contents
		pr.silentCount = 51

		got := pr.detectSpeechStopped()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Should return false when speech continue", func(t *testing.T) {
		interval := 3
		pr := NewPCMRecorder(interval)
		want := false

		contents := make([]int16, 64)
		for i := 0; i < len(contents); i++ {
			contents[i] = 1
		}
		pr.BufferedContents = contents
		pr.silentCount = 0

		got := pr.detectSpeechStopped()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Should return false when silence was very short", func(t *testing.T) {
		interval := 3
		pr := NewPCMRecorder(interval)
		want := false

		contents := make([]int16, 64)
		for i := 0; i < len(contents)-10; i++ {
			contents[i] = 1
		}
		pr.BufferedContents = contents
		pr.silentCount = 10

		got := pr.detectSpeechStopped()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestDetectSpeechExceededLimitation(t *testing.T) {
	t.Run("Should return true when speech duration is over an interval", func(t *testing.T) {
		interval := 3
		pr := NewPCMRecorder(interval)
		want := true

		pr.BufferedContents = make([]int16, 44100*pr.Interval)
		got := pr.detectSpeechExceededLimitation()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Should return false when speech duration is not over an interval", func(t *testing.T) {
		interval := 3
		pr := NewPCMRecorder(interval)
		want := false

		pr.BufferedContents = make([]int16, 44100*pr.Interval-1)
		got := pr.detectSpeechExceededLimitation()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
