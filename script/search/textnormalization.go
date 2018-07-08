package search

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func normalize(in string) string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	in = strings.ToLower(in)
	inb := []byte(in)

	out := make([]byte, len(inb))
	n, _, err := t.Transform(out, inb, true)
	if err != nil || n == 0 {
		return in
	}
	return string(out[:n])
}
