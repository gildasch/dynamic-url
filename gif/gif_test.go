package gif

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testURLs = []string{
	"https://scontent-cdg2-1.cdninstagram.com/vp/0851b45f6256d567a4c93dd6aee7893e/5BEA6009/t51.2885-15/s750x750/sh0.08/e35/35520976_2091303054444657_7754583553374945280_n.jpg",
	"https://scontent-cdg2-1.cdninstagram.com/vp/533f8d48ecb3660ff09777cbae1ef12f/5BB8EC67/t51.2885-15/s750x750/sh0.08/e35/35549378_947247585448907_3986929591536058368_n.jpg",
	"https://scontent-cdg2-1.cdninstagram.com/vp/7a6cd317859eb5e2f6f13ee8ad33ca2d/5BB0593F/t51.2885-15/s750x750/sh0.08/e35/34859331_204316296864393_4017932503724589056_n.jpg",
	"https://scontent-cdg2-1.cdninstagram.com/vp/c739ff8c13493657d9f386519b389c06/5BB724E6/t51.2885-15/s750x750/sh0.08/e35/34395139_1765017003578723_8080607223664869376_n.jpg",
	"https://scontent-cdg2-1.cdninstagram.com/vp/0e840903d7b49bb2e9fb1fd92d11c7ea/5BA628BA/t51.2885-15/e15/34823383_477366369350498_1710509268668514304_n.jpg",
}

func TestGIFFromURLs(t *testing.T) {
	b, err := MakeGIFFromURLs(testURLs, time.Second, Sierra2{})
	require.NoError(t, err)

	err = ioutil.WriteFile("out.gif", b, 0644)
	require.NoError(t, err)
}

func TestQualityOfConverters(t *testing.T) {
	b, err := MakeGIFFromURLs(testURLs, time.Second, StandardQuantizer{})
	require.NoError(t, err)

	err = ioutil.WriteFile("StandardQuantizer.gif", b, 0644)
	require.NoError(t, err)

	b, err = MakeGIFFromURLs(testURLs, time.Second, MedianCut{})
	require.NoError(t, err)

	err = ioutil.WriteFile("MedianCut.gif", b, 0644)
	require.NoError(t, err)

	b, err = MakeGIFFromURLs(testURLs, time.Second, FloydSteinberg{})
	require.NoError(t, err)

	err = ioutil.WriteFile("FloydSteinberg.gif", b, 0644)
	require.NoError(t, err)

	b, err = MakeGIFFromURLs(testURLs, time.Second, Sierra2{})
	require.NoError(t, err)

	err = ioutil.WriteFile("sierra2.gif", b, 0644)
	require.NoError(t, err)
}
