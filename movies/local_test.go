package movies

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestLocal(t *testing.T) {
	l, err := NewLocal(os.Getenv("FFMPEG_TEST_MOVIE"), "", 1024, 576)
	require.NoError(t, err)

	assert.Equal(t, 1*time.Hour+11*time.Minute+46*time.Second, l.Duration())
	inbound := rand.Intn(int(l.Duration()))
	assert.NotNil(t, l.Frame(time.Duration(inbound)))
	outbound := l.Duration() + 1*time.Minute
	assert.Nil(t, l.Frame(time.Duration(outbound)))
}
