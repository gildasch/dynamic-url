package script

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

type Script struct {
	quotes []Quote
}

type Quote struct {
	Time  time.Time
	Quote string
}

func NewScript(path string) (*Script, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var script Script
	err = json.NewDecoder(f).Decode(&script)
	if err != nil {
		return nil, err
	}

	sort.Slice(script.quotes, func(i, j int) bool {
		return script.quotes[i].Time.Before(script.quotes[j].Time)
	})

	return &script, nil
}

func (s *Script) At(t time.Time) string {
	if len(s.quotes) == 0 {
		return ""
	}
	if len(s.quotes) == 1 {
		return s.quotes[0].Quote
	}

	n := sort.Search(len(s.quotes), func(i int) bool {
		return s.quotes[i].Time.After(t)
	})

	if n == 0 {
		return s.quotes[0].Quote
	}

	if n >= len(s.quotes) {
		return s.quotes[len(s.quotes)-1].Quote
	}

	if t.Sub(s.quotes[n-1].Time) < s.quotes[n].Time.Sub(t) {
		return s.quotes[n-1].Quote
	}

	return s.quotes[n].Quote
}
