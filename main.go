package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"time"

	docker "docker.io/go-docker"
	"github.com/gildasch/dynamic-url/faceswap"
	"github.com/gildasch/dynamic-url/instagram"
	"github.com/gildasch/dynamic-url/movies"
	"github.com/gildasch/dynamic-url/movies/ffmpeg"
	"github.com/gildasch/dynamic-url/script"
	"github.com/gildasch/dynamic-url/script/search"
	"github.com/gildasch/dynamic-url/utils"
	"github.com/gildasch/dynamic-url/utils/gif"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const framesPerSecond = 5

func main() {
	instagramLoginPtr := flag.Bool("instagram-login", false, "log-in to instagram and export connection file")
	confPath := flag.String("conf", "", "path to conf")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	ffmpeg.Debug = *debug

	var conf conf
	if confPath != nil {
		var err error
		conf, err = parseConfFile(*confPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	insta, err := instagram.NewClient(".goinsta", *instagramLoginPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var store persistence.CacheStore
	if conf.RedisCache != "" {
		store = persistence.NewRedisCache(conf.RedisCache, "", 365*24*time.Hour)
	} else {
		store = persistence.NewInMemoryStore(365 * 24 * time.Hour)
	}

	router := gin.Default()

	router.GET("/instagram/user/:username/10.jpg", instagramHandler(insta, "user", "jpg"))
	router.GET("/instagram/user/:username/10.gif",
		cache.CachePage(store, 12*time.Hour, instagramHandler(insta, "user", "gif")))
	router.GET("/instagram/tag/:tag/10.jpg", instagramHandler(insta, "tag", "jpg"))
	router.GET("/instagram/tag/:tag/10.gif",
		cache.CachePage(store, 12*time.Hour, instagramHandler(insta, "tag", "gif")))

	if conf.Movies != nil {
		var ms []movies.Movie
		indexes := make(map[string]*search.Index)

		for name, mc := range conf.Movies {
			var captions movies.Captions

			if mc.Script != "" {
				captions, err = script.NewScript(mc.Script, 10*time.Second)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else if mc.Subtitles != "" {
				subtitles, err := script.NewSubtitles(mc.Subtitles)
				if err != nil {
					fmt.Println(err)
					return
				}

				indexes[name] = search.NewIndex(subtitles)
				captions = subtitles
			}

			l, err := movies.NewLocal(name, mc.Movie, mc.Subtitles, captions, 1024/2, 576/2)
			if err != nil {
				fmt.Println(err)
				return
			}

			ms = append(ms, l)
		}

		dockerClient, err := docker.NewEnvClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		fswap := &faceswap.Wuhuikais{Client: dockerClient}

		router.GET("/movies/:name/at/:at/1.jpg", movieHandler(ms, fswap, "jpg"))
		router.GET("/movies/:name/at/:at/1.gif",
			cache.CachePage(store, 365*24*time.Hour, movieHandler(ms, fswap, "gif")))
		router.GET("/movies/:name/at/:at/1.webm",
			cache.CachePage(store, 365*24*time.Hour, movieHandler(ms, fswap, "webm")))
		router.GET("/movies/:name/all.html",
			cache.CachePage(store, 365*24*time.Hour, movieAllHandler(ms, "gif")))
		router.GET("/movies/:name/search.html",
			cache.CachePage(store, 365*24*time.Hour, movieSearchHandler(indexes)))
	}

	router.Run()
}

func instagramHandler(insta *instagram.Client, search string, format string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var urls []string
		var err error
		switch search {
		case "user":
			urls, err = insta.GetLatestPicturesFromUser(c.Param("username"), 10)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return
			}
		case "tag":
			urls, err = insta.GetLatestPicturesFromTag(c.Param("tag"), 100, 1920, 1920)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return
			}
		default:
			c.Status(http.StatusBadRequest)
			return
		}

		switch format {
		case "jpg":
			c.Redirect(http.StatusFound, utils.OneOf(urls))
			return
		case "gif":
			delay, err := time.ParseDuration(c.DefaultQuery("delay", "1s"))
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			var convert gif.Converter = gif.MedianCut{}
			if c.Query("dither") == "true" {
				convert = gif.FloydSteinberg{}
			}

			gif, err := gif.MakeGIFFromURLs(urls, delay, convert)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/gif", gif)
			return
		default:
		}

		c.Status(http.StatusBadRequest)
		return
	}
}

func movieHandler(ms []movies.Movie, fswap *faceswap.Wuhuikais, format string) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")

		var movie movies.Movie
		for _, m := range ms {
			if m.Name() == name {
				movie = m
				break
			}
		}

		if movie == nil {
			fmt.Println("movie not found")
			c.Status(http.StatusBadRequest)
			return
		}

		at, err := time.ParseDuration(c.Param("at"))
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		var caption string
		if c.Query("text") != "" {
			caption = c.Query("text")
		} else {
			caption = movie.Caption(at)
		}

		switch format {
		case "jpg":
			jpg := utils.WithCaption(movie.Frame(at), nil, caption)
			var buf bytes.Buffer
			err = jpeg.Encode(&buf, jpg, nil)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
			return
		case "gif":
			gab := c.Query("gab") == "1"

			nframes := 50
			if c.Query("frames") != "" {
				i, err := strconv.Atoi(c.Query("frames"))
				if err == nil && nframes <= 200 {
					nframes = i
				}
			}

			var frames []image.Image
			if gab {
				frames = movie.Frames(at, 1, framesPerSecond)
			} else {
				frames = movie.Frames(at, nframes, framesPerSecond)
			}
			var withCaption []image.Image

			startTimeInMovie := at
			endTimeInMovie := at + time.Duration(nframes)*time.Second/time.Duration(framesPerSecond)

			previousCaption := caption
			var bounds *image.Rectangle
			for i := 0; i < len(frames); i++ {
				f := frames[i]
				if bounds == nil {
					fb := f.Bounds()
					bounds = &fb
				}

				timeInMovie := at + time.Duration(i)*time.Second/time.Duration(framesPerSecond)

				if c.Query("text") != "" {
					caption = c.Query("text")
				} else {
					caption = movie.Caption(timeInMovie)

					if timeInMovie < startTimeInMovie+300*time.Millisecond &&
						caption != movie.Caption(startTimeInMovie+time.Second) {
						caption = ""
					}
					if timeInMovie > endTimeInMovie-300*time.Millisecond &&
						caption != movie.Caption(endTimeInMovie-time.Second) {
						caption = ""
					}
				}

				if !gab && caption != previousCaption {
					at = timeInMovie
					nframes = nframes - i
					frames = movie.Frames(at, nframes, framesPerSecond)
					bounds = nil
					i = -1
					previousCaption = caption
					continue
				}

				if gab {
					// Faceswap
					tmpImage, err := uuid.NewV4()
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					tmp, err := os.Create("/tmp/" + tmpImage.String() + ".jpg")
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					err = jpeg.Encode(tmp, f, &jpeg.Options{Quality: 90})
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					defer func() {
						os.Remove("/tmp/" + tmpImage.String() + ".jpg")
					}()

					err = tmp.Close()
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					err = fswap.FaceSwap("/tmp", "gab2.jpg", tmpImage.String()+".jpg", tmpImage.String()+".jpg")
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					tmp, err = os.Open("/tmp/" + tmpImage.String() + ".jpg")
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					f, err = jpeg.Decode(tmp)
					if err != nil {
						fmt.Println(err)
						c.Status(http.StatusInternalServerError)
						return
					}

					// Get next frame
					nextTimeInMovie := at + time.Duration(i+1)*time.Second/time.Duration(framesPerSecond)
					if nextTimeInMovie <= endTimeInMovie {
						at = nextTimeInMovie
						frames = movie.Frames(at, 1, framesPerSecond)
						i = -1
					}
				}

				previousCaption = caption
				withCaption = append(withCaption, utils.WithCaption(f, bounds, caption))
			}

			var convert gif.Converter = gif.StandardQuantizer{}
			if c.Query("dither") == "true" {
				convert = gif.FloydSteinberg{}
			}

			gif, err := gif.MakeGIFFromImages(withCaption, 150*time.Millisecond, convert)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/gif", gif)
			return
		case "webm":
			nframes := 50
			if c.Query("frames") != "" {
				i, err := strconv.Atoi(c.Query("frames"))
				if err == nil && nframes <= 200 {
					nframes = i
				}
			}

			webm, err := movie.WebM(at, nframes*25/framesPerSecond, 25)
			if err != nil {
				c.Data(http.StatusInternalServerError, "video/webm", []byte{})
				return
			}

			c.Data(http.StatusOK, "video/webm", webm)
			return
		}

		c.Status(http.StatusBadRequest)
	}
}

func movieAllHandler(ms []movies.Movie, format string) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")

		var movie movies.Movie
		for _, m := range ms {
			if m.Name() == name {
				movie = m
				break
			}
		}

		if movie == nil {
			fmt.Println("movie not found")
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
		fmt.Fprintf(c.Writer, "<table>")
		for i := 0 * time.Second; i < movie.Duration(); i += 10 * time.Second {
			fmt.Fprintf(c.Writer, `
<tr>
  <td>
    <img src='/movies/%s/at/%s/1.jpg' style='max-width:240px' />
  </td>
  <td>
    %s
  </td>
  <td>
    <a href='/movies/%s/at/%s/1.jpg'>JPEG</a>
  </td>
  <td>
    <a href='/movies/%s/at/%s/1.gif'>GIF</a>
  </td>
</tr>`,
				name, i, i, name, i, name, i)
		}
		fmt.Fprintf(c.Writer, "</table>")
	}
}

func movieSearchHandler(indexes map[string]*search.Index) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		query := c.Query("query")

		if len(query) <= 1 {
			c.Status(http.StatusBadRequest)
			return
		}

		index, ok := indexes[name]
		if !ok {
			fmt.Println("movie not found")
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
		fmt.Fprintf(c.Writer, "<table>")

		res := index.Search(query)
		for _, at := range res {
			fmt.Fprintf(c.Writer, `
<tr>
  <td>
    <img src='/movies/%s/at/%s/1.jpg' style='max-width:240px' />
  </td>
  <td>
    %s
  </td>
  <td>
    <a href='/movies/%s/at/%s/1.jpg'>JPEG</a>
  </td>
  <td>
    <a href='/movies/%s/at/%s/1.gif'>GIF</a>
  </td>
</tr>`,
				name, at, at, name, at, name, at)
		}
		fmt.Fprintf(c.Writer, "</table>")
	}
}
