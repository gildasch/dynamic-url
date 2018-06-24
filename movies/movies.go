package movies

import (
	"image"
	"time"
)

type Movie interface {
	Name() string
	Duration() time.Duration
	Frame(at time.Duration) image.Image
	Frames(at time.Duration, n int) []image.Image
	Caption(at time.Duration) string
}
