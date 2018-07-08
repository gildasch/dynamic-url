package search

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
