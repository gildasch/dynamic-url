package search

import (
	"testing"
	"time"

	"github.com/gildasch/dynamic-url/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockQuotes struct {
	at  []time.Duration
	txt []string
}

func (mq *MockQuotes) Len() int {
	return len(mq.at)
}

func (mq *MockQuotes) Quote(i int) (string, time.Duration) {
	return mq.txt[i], mq.at[i]
}

func TestIndexSearch(t *testing.T) {
	quotes := &MockQuotes{
		at: []time.Duration{
			1 * time.Second,
			2 * time.Second,
			3 * time.Second,
			4 * time.Second,
			5 * time.Second,
			6 * time.Second,
		},
		txt: []string{
			`<i>Dieu de miséricorde, ce ne sont pas
des choses qu'on fait.</i>
`,
			"Applaudissements",
			"Je l'ai pas vu.",
			"Bravo!",
			`En coupant le 3e temps,
à la fin.`,
			"Bien net, comme ça.",
		},
	}

	i := NewIndex(quotes)

	assert.Equal(t, []time.Duration{1 * time.Second}, i.Search("misericorde"))
	assert.Equal(t, []time.Duration{3 * time.Second}, i.Search("l'ai"))
	assert.Equal(t, []time.Duration{5 * time.Second}, i.Search("3"))
	assert.Equal(t, []time.Duration{6 * time.Second}, i.Search("bien net"))
}

func TestIndexSearchSubtitle(t *testing.T) {
	s, err := script.NewSubtitles("test.srt")
	require.NoError(t, err)

	i := NewIndex(s)

	assert.Equal(t, []time.Duration{
		15*time.Minute + 47*time.Second + 711*time.Millisecond,
		57*time.Minute + 52*time.Second + 714*time.Millisecond,
		58*time.Minute + 21*time.Second + 631*time.Millisecond,
		58*time.Minute + 45*time.Second + 298*time.Millisecond,
		59*time.Minute + 11*time.Second + 007*time.Millisecond,
		59*time.Minute + 21*time.Second + 674*time.Millisecond,
		59*time.Minute + 28*time.Second + 507*time.Millisecond,
		1*time.Hour + 16*time.Minute + 38*time.Second + 359*time.Millisecond,
		1*time.Hour + 28*time.Minute + 38*time.Second + 913*time.Millisecond,
	}, i.Search("drole"))
}

func BenchmarkIndexSearch(b *testing.B) {
	s, err := script.NewSubtitles("test.srt")
	require.NoError(b, err)

	i := NewIndex(s)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i.Search("drole")
	}
}
