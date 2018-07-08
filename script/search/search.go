package search

import (
	"strings"
	"time"
)

type Index struct {
	quotes []quote
}

type quote struct {
	text string
	at   time.Duration
}

type Quotes interface {
	Len() int
	Quote(i int) (string, time.Duration)
}

func NewIndex(qs Quotes) *Index {
	index := Index{}

	for i := 0; i < qs.Len(); i++ {
		text, at := qs.Quote(i)
		index.quotes = append(index.quotes, quote{
			text: normalize(text),
			at:   at,
		})
	}

	return &index
}

func (i *Index) Search(query string) []time.Duration {
	query = normalize(query)

	res := []time.Duration{}

	for _, q := range i.quotes {
		if strings.Contains(q.text, query) {
			res = append(res, q.at)
		}
	}

	return res
}
