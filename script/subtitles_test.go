package script

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubtitles(t *testing.T) {
	s, err := NewSubtitles("lca.srt")
	require.NoError(t, err)

	assert.Len(t, s.quotes, 11)

	assert.Equal(t,
		"Voilà enfin le roi de la classe ! L’homme trop bien sapé, Abitbol !",
		s.At(48*time.Second))
	assert.Equal(t,
		"",
		s.At(47*time.Second))
	assert.Equal(t,
		"Le grand play-boy des fonds marins qui fait rêver les ménagères…",
		s.At(1*time.Minute+2*time.Second+408*time.Millisecond))
	assert.Equal(t,
		"",
		s.At(1*time.Minute+2*time.Second+409*time.Millisecond))
	assert.Equal(t,
		"Mais moi je les baise, les ménagères. Pas vrai ?",
		s.At(1*time.Minute+2*time.Second+458*time.Millisecond))
	assert.Equal(t,
		"",
		s.At(0))
	assert.Equal(t,
		"",
		s.At(2*time.Hour+1*time.Minute))
}
