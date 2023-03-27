package model

import "time"

type Transcription struct {
	Text      string
	Timestamp time.Time
}
