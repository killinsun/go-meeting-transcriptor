package pcm

type PCMRecorder struct {
	Interval         int
	SilentRatio      float32
	BaseLangCode     string
	AltLangCodes     []string
	BufferedContents []int16
	silentCount      int
}

func NewPCMRecorder(interval int) *PCMRecorder {
	var pr = &PCMRecorder{
		Interval: interval,
	}
	return pr
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
