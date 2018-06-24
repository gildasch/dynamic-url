package utils

import (
	"image/jpeg"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithCaption(t *testing.T) {
	f, err := os.Open("source.jpg")
	require.NoError(t, err)

	in, err := jpeg.Decode(f)
	require.NoError(t, err)

	out := WithCaption(in, "Some text to print as a caption. Could be something else")

	fout, err := os.OpenFile("withcaption.jpg", os.O_RDWR|os.O_CREATE, 0644)
	require.NoError(t, err)

	err = jpeg.Encode(fout, out, nil)
	require.NoError(t, err)
}
