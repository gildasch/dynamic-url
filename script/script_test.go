package script

import (
	"testing"
	"time"

	"github.com/gildasch/dynamic-url/movies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScript(t *testing.T) {
	s, err := NewScript("script.json", 10*time.Second)
	require.NoError(t, err)
	assert.Len(t, s.quotes, 12)
	assert.Equal(t,
		"Attention, ce flim n’est pas un flim sur le cyclimse. Merci de votre compréhension.",
		s.quotes[0].Quote)
	assert.Equal(t, "1m13s", s.quotes[11].At.String())
}

func TestAt(t *testing.T) {
	script := Script{quotes: []Quote{
		Quote{
			At:    11*time.Hour + 22*time.Minute + 33*time.Second,
			Quote: "première"},
		Quote{
			At:    11*time.Hour + 22*time.Minute + 43*time.Second,
			Quote: "deuxième"},
		Quote{
			At:    11*time.Hour + 22*time.Minute + 45*time.Second,
			Quote: "troisième"},
		Quote{
			At:    11*time.Hour + 22*time.Minute + 50*time.Second,
			Quote: "quatrième"},
		Quote{
			At:    11*time.Hour + 23*time.Minute + 33*time.Second,
			Quote: "cinquième"}}}

	assert.Equal(t, "première", script.At(10*time.Hour+22*time.Minute+33*time.Second+000*time.Millisecond))
	assert.Equal(t, "première", script.At(11*time.Hour+22*time.Minute+37*time.Second+590*time.Millisecond))
	assert.Equal(t, "deuxième", script.At(11*time.Hour+22*time.Minute+38*time.Second+010*time.Millisecond))
	assert.Equal(t, "troisième", script.At(11*time.Hour+22*time.Minute+44*time.Second+320*time.Millisecond))
	assert.Equal(t, "quatrième", script.At(11*time.Hour+23*time.Minute+00*time.Second+320*time.Millisecond))
	assert.Equal(t, "cinquième", script.At(11*time.Hour+23*time.Minute+20*time.Second+000*time.Millisecond))
	assert.Equal(t, "cinquième", script.At(11*time.Hour+23*time.Minute+59*time.Second+000*time.Millisecond))
	assert.Equal(t, "cinquième", script.At(11*time.Hour+23*time.Minute+59*time.Second+000*time.Millisecond))
}

func TestBetween(t *testing.T) {
	script := Script{quotes: []Quote{
		Quote{
			At:    11*time.Hour + 22*time.Minute + 33*time.Second,
			Quote: "première"},
		Quote{
			At:    11*time.Hour + 22*time.Minute + 43*time.Second,
			Quote: "deuxième"},
		Quote{
			At:    11*time.Hour + 22*time.Minute + 45*time.Second,
			Quote: "troisième"},
		Quote{
			At:    11*time.Hour + 22*time.Minute + 50*time.Second,
			Quote: "quatrième"},
		Quote{
			At:    11*time.Hour + 23*time.Minute + 33*time.Second,
			Quote: "cinquième"}}}

	assert.Equal(t,
		[]movies.Caption{{
			Text:  "première",
			Start: 11*time.Hour + 22*time.Minute + 33*time.Second,
			End:   11*time.Hour + 22*time.Minute + 43*time.Second}, {
			Text:  "deuxième",
			Start: 11*time.Hour + 22*time.Minute + 43*time.Second,
			End:   40965000000000}, {
			Text:  "troisième",
			Start: 11*time.Hour + 22*time.Minute + 45*time.Second,
			End:   11*time.Hour + 22*time.Minute + 50*time.Second}},
		script.Between(
			11*time.Hour+22*time.Minute+43*time.Second-1*time.Millisecond,
			11*time.Hour+22*time.Minute+50*time.Second))
	assert.Equal(t,
		[]movies.Caption{{
			Text:  "deuxième",
			Start: 11*time.Hour + 22*time.Minute + 43*time.Second,
			End:   40965000000000}, {
			Text:  "troisième",
			Start: 11*time.Hour + 22*time.Minute + 45*time.Second,
			End:   11*time.Hour + 22*time.Minute + 49*time.Second}},
		script.Between(
			11*time.Hour+22*time.Minute+43*time.Second,
			11*time.Hour+22*time.Minute+49*time.Second))
}
