package ffmpeg

import (
	"image"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	d, err := Duration(os.Getenv("FFMPEG_TEST_MOVIE"))
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Hour+11*time.Minute+46*time.Second, d)
}

func BenchmarkDuration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Duration(os.Getenv("FFMPEG_TEST_MOVIE"))
	}
}

func TestDurationNotFound(t *testing.T) {
	_, err := Duration("/does/not/exists.mkv")
	assert.Error(t, err)
}

func TestDurationInvalidFile(t *testing.T) {
	_, err := Duration("main.go")
	assert.Error(t, err)
}

func TestCapture(t *testing.T) {
	capt, err := Capture(os.Getenv("FFMPEG_TEST_MOVIE"), 1*time.Hour, 1024, 576)
	require.NoError(t, err)
	assert.Equal(t, image.Rect(0, 0, 1024, 576), capt.Bounds())
}

func BenchmarkCapture(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Capture(os.Getenv("FFMPEG_TEST_MOVIE"), 1*time.Hour, 1024, 576)
	}
}
