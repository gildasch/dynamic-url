package script

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAt(t *testing.T) {
	script := Script{quotes: []Quote{
		Quote{
			Time:  time.Date(2000, 1, 1, 11, 22, 33, 0, time.UTC),
			Quote: "première"},
		Quote{
			Time:  time.Date(2000, 1, 1, 11, 22, 43, 0, time.UTC),
			Quote: "deuxième"},
		Quote{
			Time:  time.Date(2000, 1, 1, 11, 22, 45, 0, time.UTC),
			Quote: "troisième"},
		Quote{
			Time:  time.Date(2000, 1, 1, 11, 22, 50, 0, time.UTC),
			Quote: "quatrième"},
		Quote{
			Time:  time.Date(2000, 1, 1, 11, 23, 33, 0, time.UTC),
			Quote: "cinquième"},
	}}

	assert.Equal(t, "première", script.At(time.Date(2000, 1, 1, 10, 22, 33, 0, time.UTC)))
	assert.Equal(t, "première", script.At(time.Date(2000, 1, 1, 11, 22, 37, 59, time.UTC)))
	assert.Equal(t, "deuxième", script.At(time.Date(2000, 1, 1, 11, 22, 38, 01, time.UTC)))
	assert.Equal(t, "troisième", script.At(time.Date(2000, 1, 1, 11, 22, 44, 32, time.UTC)))
	assert.Equal(t, "quatrième", script.At(time.Date(2000, 1, 1, 11, 23, 00, 32, time.UTC)))
	assert.Equal(t, "cinquième", script.At(time.Date(2000, 1, 1, 11, 23, 20, 00, time.UTC)))
	assert.Equal(t, "cinquième", script.At(time.Date(2000, 1, 1, 11, 23, 59, 00, time.UTC)))
	assert.Equal(t, "cinquième", script.At(time.Date(2029, 1, 1, 11, 23, 59, 00, time.UTC)))
}
