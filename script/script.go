package script

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gildasch/dynamic-url/movies"
	"github.com/pkg/errors"
)

type Script struct {
	quotes []Quote
}

type Quote struct {
	At    time.Duration
	Time  string
	Quote string
}

func NewScript(path string, correction time.Duration) (*Script, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var script Script
	err = json.NewDecoder(f).Decode(&script.quotes)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(script.quotes); i++ {
		splitted := strings.Split(script.quotes[i].Time, ":")
		if len(splitted) != 3 {
			return nil, errors.Errorf("Time %q from json is malformed", script.quotes[i].Time)
		}
		h, err := strconv.Atoi(splitted[0])
		if err != nil {
			return nil, errors.Wrapf(err, "Time %q from json is malformed", script.quotes[i].Time)
		}
		m, err := strconv.Atoi(splitted[1])
		if err != nil {
			return nil, errors.Wrapf(err, "Time %q from json is malformed", script.quotes[i].Time)
		}
		s, err := strconv.Atoi(splitted[2])
		if err != nil {
			return nil, errors.Wrapf(err, "Time %q from json is malformed", script.quotes[i].Time)
		}
		script.quotes[i].At = time.Duration(h)*time.Hour +
			time.Duration(m)*time.Minute +
			time.Duration(s)*time.Second +
			correction
	}

	sort.Slice(script.quotes, func(i, j int) bool {
		return script.quotes[i].At < script.quotes[j].At
	})

	return &script, nil
}

func (s *Script) At(t time.Duration) string {
	if len(s.quotes) == 0 {
		return ""
	}
	if len(s.quotes) == 1 {
		return s.quotes[0].Quote
	}

	n := sort.Search(len(s.quotes), func(i int) bool {
		return s.quotes[i].At > t
	})

	if n == 0 {
		return s.quotes[0].Quote
	}

	if n >= len(s.quotes) {
		return s.quotes[len(s.quotes)-1].Quote
	}

	if t-s.quotes[n-1].At < s.quotes[n].At-t {
		return s.quotes[n-1].Quote
	}

	return s.quotes[n].Quote
}

func (s *Script) Between(start, end time.Duration) []movies.Caption {
	if len(s.quotes) == 0 {
		return nil
	}

	n := sort.Search(len(s.quotes), func(i int) bool {
		return s.quotes[i].At > start
	})

	if n >= len(s.quotes) {
		return nil
	}

	n = n - 1

	var starts []time.Duration
	var texts []string
	for n+1 < len(s.quotes) && s.quotes[n].At < end {
		starts = append(starts, s.quotes[n].At)
		texts = append(texts, s.quotes[n].Quote)
		n++
	}

	var captions []movies.Caption
	i := 0
	for ; i < len(starts)-1; i++ {
		captions = append(captions, movies.Caption{
			Text:  texts[i],
			Start: starts[i],
			End:   starts[i+1],
		})
	}
	captions = append(captions, movies.Caption{
		Text:  texts[i],
		Start: starts[i],
		End:   end,
	})

	return captions
}
