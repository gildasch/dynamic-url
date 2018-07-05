package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConf(t *testing.T) {
	c, err := parseConf(`
movies:
  lca:
    movie: ici.mkv
    subtitles: la.srt
`)

	require.NoError(t, err)
	assert.Equal(t, conf{
		Movies: map[string]movieConf{
			"lca": movieConf{
				Movie:     "ici.mkv",
				Subtitles: "la.srt",
			},
		},
	}, c)
}
