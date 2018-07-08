package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	assert.Equal(t, "je m'appelle gilas", normalize("je m'appèllê Giläs"))
	assert.Equal(t, `<i>dieu de misericorde, ce ne sont pas
des choses qu'on fait.</i>
`, normalize(
		`<i>Dieu de miséricorde, ce ne sont pas
des choses qu'on fait.</i>
`))
}
