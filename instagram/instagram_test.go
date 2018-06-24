package instagram

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T) *Client {
	if os.Getenv("GOINSTA_CONF") == "" {
		t.SkipNow()
		return nil
	}

	c, err := NewClient(os.Getenv("GOINSTA_CONF"), false)
	require.NoError(t, err)
	return c
}

func TestGetLatestPicturesFromUser(t *testing.T) {
	c := testClient(t)

	urls, err := c.GetLatestPicturesFromUser("kimkardashian", 10)
	assert.NoError(t, err)
	assert.Len(t, urls, 10)
}

func TestGetLatestPicturesFromTag(t *testing.T) {
	c := testClient(t)

	urls, err := c.GetLatestPicturesFromTag("vinnature", 10, 1920, 1920)
	assert.NoError(t, err)
	assert.Len(t, urls, 10)
}
