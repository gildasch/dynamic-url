package movies

import (
	"fmt"
	"image"
	"time"

	"github.com/gildasch/dynamic-url/movies/ffmpeg"
	"github.com/pkg/errors"
)

type Local struct {
	name          string
	video         string
	sub           string
	captions      Captions
	width, height int

	duration time.Duration
}

type Captions interface {
	At(t time.Duration) string
	Between(start, end time.Duration) []Caption
}

func NewLocal(name, video, sub string, captions Captions, width, height int) (*Local, error) {
	d, err := ffmpeg.Duration(video)
	if err != nil {
		return nil, errors.Wrapf(err, "could not inspect movie file %q", video)
	}

	return &Local{
		name:     name,
		video:    video,
		sub:      sub,
		captions: captions,
		width:    width,
		height:   height,
		duration: d,
	}, nil
}

func (l *Local) Name() string {
	return l.name
}

func (l *Local) Duration() time.Duration {
	return l.duration
}

func (l *Local) Frame(at time.Duration) image.Image {
	i, err := ffmpeg.Capture(l.video, at, l.width, l.height)
	if err != nil {
		fmt.Println("unexpected error:", err)
	}
	return i
}

func (l *Local) Frames(at time.Duration, n, framesPerSecond int) []image.Image {
	is, err := ffmpeg.GIFCaptures(l.video, at, l.width, l.height, n, framesPerSecond)
	if err != nil {
		fmt.Println("unexpected error:", err)
	}

	var out []image.Image
	for _, i := range is {
		out = append(out, i)
	}

	return out
}

func (l *Local) WebM(at time.Duration, n, framesPerSecond int) ([]byte, error) {
	return ffmpeg.WebM(l.video, l.sub, at, l.width, l.height, n, framesPerSecond)
}

func (l *Local) Caption(at time.Duration) string {
	return l.captions.At(at)
}

func (l *Local) CaptionBetween(start, end time.Duration) []Caption {
	return l.captions.Between(start, end)
}
