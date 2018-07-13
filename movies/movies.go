package movies

import (
	"image"
	"time"
)

type Movie interface {
	Name() string
	Duration() time.Duration
	Frame(at time.Duration) image.Image
	Frames(at time.Duration, n, framesPerSecond int) []image.Image
	Caption(at time.Duration) string
	CaptionBetween(start, end time.Duration) []Caption
}

type Caption struct {
	Text  string
	Start time.Duration
	End   time.Duration
}
