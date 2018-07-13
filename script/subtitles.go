package script

import (
	"sort"
	"time"

	astisub "github.com/asticode/go-astisub"
	"github.com/gildasch/dynamic-url/movies"
	"github.com/pkg/errors"
)

type Subtitles struct {
	quotes []SubtitleQuote
}

type SubtitleQuote struct {
	Start time.Duration
	End   time.Duration
	Quote string
}

func NewSubtitles(path string) (*Subtitles, error) {
	s, err := astisub.OpenFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening subtitle file %q", path)
	}
	s.Order()

	var quotes []SubtitleQuote
	for _, i := range s.Items {
		t := ""
		for _, l := range i.Lines {
			if t != "" {
				t += "\n"
			}
			t += l.String()
		}

		if t == "" {
			continue
		}

		quotes = append(quotes, SubtitleQuote{
			Start: i.StartAt,
			End:   i.EndAt,
			Quote: t,
		})
	}

	return &Subtitles{quotes: quotes}, nil
}

func (s *Subtitles) At(t time.Duration) string {
	if len(s.quotes) == 0 {
		return ""
	}

	n := sort.Search(len(s.quotes), func(i int) bool {
		return s.quotes[i].Start > t
	})

	if n == 0 {
		return ""
	}

	n = n - 1

	if s.quotes[n].Start <= t && t <= s.quotes[n].End {
		return s.quotes[n].Quote
	}

	return ""
}

func (s *Subtitles) Between(start, end time.Duration) []movies.Caption {
	if len(s.quotes) == 0 {
		return nil
	}

	n := sort.Search(len(s.quotes), func(i int) bool {
		return s.quotes[i].Start > start
	})

	if n == 0 {
		return nil
	}

	n = n - 1

	var captions []movies.Caption
	for n < len(s.quotes) && (s.quotes[n].Start < end || start > s.quotes[n].End) {
		captions = append(captions, movies.Caption{
			Text:  s.quotes[n].Quote,
			Start: s.quotes[n].Start,
			End:   s.quotes[n].End,
		})
		n++
	}

	return captions
}

func (s *Subtitles) Len() int {
	return len(s.quotes)
}

func (s *Subtitles) Quote(i int) (string, time.Duration) {
	return s.quotes[i].Quote, s.quotes[i].Start
}
