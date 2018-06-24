package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gildasch/dynamic-url/gif"
	"github.com/gildasch/dynamic-url/instagram"
	"github.com/gin-gonic/gin"
)

func oneOf(ss []string) string {
	if len(ss) == 0 {
		return ""
	}

	return ss[rand.Intn(len(ss)-1)]
}

func main() {
	instagramLoginPtr := flag.Bool("instagram-login", false, "log-in to instagram and export connexion file")
	flag.Parse()

	insta, err := instagram.NewClient(".goinsta", *instagramLoginPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()

	router.GET("/instagram/user/:username/10.jpg", func(c *gin.Context) {
		urls, err := insta.GetLatestPicturesFromUser(c.Param("username"), 10)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, oneOf(urls))
	})

	router.GET("/instagram/user/:username/10.gif", func(c *gin.Context) {
		urls, err := insta.GetLatestPicturesFromUser(c.Param("username"), 10)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

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
	})

	router.GET("/instagram/tag/:tag/10.jpg", func(c *gin.Context) {
		urls, err := insta.GetLatestPicturesFromTag(c.Param("tag"), 100, 1920, 1920)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, oneOf(urls))
	})

	router.GET("/instagram/tag/:tag/10.gif", func(c *gin.Context) {
		urls, err := insta.GetLatestPicturesFromTag(c.Param("tag"), 10, 640, 640)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

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
	})

	router.Run()
}
