package script

import (
	"testing"
	"time"

	"github.com/gildasch/dynamic-url/movies"
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

func TestSubtitlesBetween(t *testing.T) {
	s, err := NewSubtitles("lca.srt")
	require.NoError(t, err)

	actual := s.Between(41*time.Second, 1*time.Minute)

	assert.Equal(t,
		[]movies.Caption{{
			Text:  "V12 au capitaine George Abitbol, on vous demande sur le pont.",
			Start: 35*time.Second + 826*time.Millisecond,
			End:   41*time.Second + 426*time.Millisecond}, {
			Text:  "— Qui ? — Un dénommé José.",
			Start: 41*time.Second + 608*time.Millisecond,
			End:   43*time.Second + 608*time.Millisecond,
		}, {
			Text:  "OK, j’arrive.",
			Start: 43*time.Second + 558*time.Millisecond,
			End:   45*time.Second + 58*time.Millisecond,
		}, {
			Text:  "Voilà enfin le roi de la classe ! L’homme trop bien sapé, Abitbol !",
			Start: 47*time.Second + 312*time.Millisecond,
			End:   52*time.Second + 512*time.Millisecond,
		}, {
			Text:  "T’as été élu l’homme le plus classe du monde ? Laisse-moi rire !",
			Start: 52*time.Second + 558*time.Millisecond,
			End:   57*time.Second + 558*time.Millisecond,
		}, {
			Text:  "Le grand play-boy des fonds marins qui fait rêver les ménagères…",
			Start: 57*time.Second + 608*time.Millisecond,
			End:   1*time.Minute + 2*time.Second + 408*time.Millisecond,
		}},
		actual)
}
