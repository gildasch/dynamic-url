package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"time"

	"github.com/gildasch/dynamic-url/instagram"
	"github.com/gildasch/dynamic-url/movies"
	"github.com/gildasch/dynamic-url/script"
	"github.com/gildasch/dynamic-url/utils"
	"github.com/gildasch/dynamic-url/utils/gif"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

const framesPerSecond = 5

func main() {
	instagramLoginPtr := flag.Bool("instagram-login", false, "log-in to instagram and export connection file")
	lcaMovie := flag.String("lca-movie", "", "path to the movie file of lca")
	lcaScript := flag.String("lca-script", "", "path to the script file of lca")
	lcaSubtitles := flag.String("lca-subs", "", "path to the subtitle file of lca")
	flag.Parse()

	insta, err := instagram.NewClient(".goinsta", *instagramLoginPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	store := persistence.NewInMemoryStore(365 * 24 * time.Hour)

	router := gin.Default()

	router.GET("/instagram/user/:username/10.jpg", instagramHandler(insta, "user", "jpg"))
	router.GET("/instagram/user/:username/10.gif",
		cache.CachePage(store, 12*time.Hour, instagramHandler(insta, "user", "gif")))
	router.GET("/instagram/tag/:tag/10.jpg", instagramHandler(insta, "tag", "jpg"))
	router.GET("/instagram/tag/:tag/10.gif",
		cache.CachePage(store, 12*time.Hour, instagramHandler(insta, "tag", "gif")))

	var ms []movies.Movie

	var captions movies.Captions

	if *lcaScript != "" {
		captions, err = script.NewScript(*lcaScript, 10*time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if *lcaSubtitles != "" {
		captions, err = script.NewSubtitles(*lcaSubtitles)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	lca, err := movies.NewLocal("lca", *lcaMovie, captions, 1024/2, 576/2)
	if err != nil {
		fmt.Println(err)
		return
	}

	ms = append(ms, lca)

	router.GET("/movies/:name/at/:at/1.jpg", movieHandler(ms, "jpg"))
	router.GET("/movies/:name/at/:at/1.gif",
		cache.CachePage(store, 365*24*time.Hour, movieHandler(ms, "gif")))
	router.GET("/movies/:name/all.html",
		cache.CachePage(store, 365*24*time.Hour, movieAllHandler(ms, "gif")))

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

func movieHandler(ms []movies.Movie, format string) func(c *gin.Context) {
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
			jpg := utils.WithCaption(movie.Frame(at), caption)
			var buf bytes.Buffer
			err = jpeg.Encode(&buf, jpg, nil)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
			return
		case "gif":
			nframes := 50
			if c.Query("frames") != "" {
				i, err := strconv.Atoi(c.Query("frames"))
				if err == nil && nframes <= 200 {
					nframes = i
				}
			}

			frames := movie.Frames(at, nframes, framesPerSecond)
			var withCaption []image.Image

			previousCaption := caption
			for i := 0; i < len(frames); i++ {
				f := frames[i]

				if c.Query("text") != "" {
					caption = c.Query("text")
				} else {
					caption = movie.Caption(at + time.Duration(i)*time.Second/time.Duration(framesPerSecond))
				}

				if caption != previousCaption {
					at = at + time.Duration(i)*time.Second/time.Duration(framesPerSecond)
					nframes = nframes - i
					frames = movie.Frames(at, nframes, framesPerSecond)
					i = -1
					previousCaption = caption
					continue
				}

				previousCaption = caption
				withCaption = append(withCaption, utils.WithCaption(f, caption))
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
