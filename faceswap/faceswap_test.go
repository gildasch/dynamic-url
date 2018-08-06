package faceswap

import (
	"io/ioutil"
	"os"
	"testing"

	docker "docker.io/go-docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFaceSwap(t *testing.T) {
	c, err := docker.NewEnvClient()
	require.NoError(t, err)

	w := Wuhuikais{Client: c}

	dir, err := os.Getwd()
	require.NoError(t, err)

	err = w.FaceSwap(dir, "gildas2.png", "Rogelio.png", "out.jpg")
	require.NoError(t, err)

	actual, err := ioutil.ReadFile("out.jpg")
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("expected.jpg")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestBuild(t *testing.T) {
	c, err := docker.NewEnvClient()
	require.NoError(t, err)

	w := Wuhuikais{Client: c}

	err = w.build()
	require.NoError(t, err)
	assert.NotZero(t, w.imageID)
}
