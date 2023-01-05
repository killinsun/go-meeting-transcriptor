package pcm

type PCMRecorder struct {
	Interval     int
	SilentRatio  float32
	BaseLangCode string
	AltLangCodes []string
}
