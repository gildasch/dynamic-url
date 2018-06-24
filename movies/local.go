package movies

import (
	"fmt"
	"image"
	"time"

	"github.com/gildasch/dynamic-url/movies/ffmpeg"
	"github.com/gildasch/dynamic-url/script"
	"github.com/pkg/errors"
)

type Local struct {
	name          string
	video         string
	captions      *script.Script
	width, height int

	duration time.Duration
}

func NewLocal(name, video string, captions *script.Script, width, height int) (*Local, error) {
	d, err := ffmpeg.Duration(video)
	if err != nil {
		return nil, errors.Wrapf(err, "could not inspect movie file %q", video)
	}

	return &Local{
		name:     name,
		video:    video,
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

func (l *Local) Frames(at time.Duration, n int) []image.Image {
	is, err := ffmpeg.GIFCaptures(l.video, at, l.width, l.height, n)
	if err != nil {
		fmt.Println("unexpected error:", err)
	}

	var out []image.Image
	for _, i := range is {
		out = append(out, i)
	}

	return out
}

func (l *Local) Caption(at time.Duration) string {
	return l.captions.At(at)
}
