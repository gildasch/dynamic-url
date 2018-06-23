package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gildasch/dynamic-url/gif"
	"github.com/gin-gonic/gin"
	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

func getLatestPicturesFromUser(i *goinsta.Instagram, username string, n int) ([]string, error) {
	user, err := i.Profiles.ByName(username)
	if err != nil {
		return nil, err
	}

	media := user.Feed()

	var urls []string

	for media.Next() {
		for _, i := range media.Items {
			if len(i.Images.Versions) == 0 {
				continue
			}
			urls = append(urls, i.Images.Versions[0].URL)
			if len(urls) == n {
				return urls, nil
			}
		}
	}

	return urls, nil
}

func getLatestPicturesFromTag(i *goinsta.Instagram, tag string, n int, maxWidth, maxHeight int) ([]string, error) {
	feed, err := i.Search.FeedTags(tag)
	if err != nil {
		return nil, err
	}

	var urls []string

	for _, i := range feed.Images {
		u := best(i.Images, maxWidth, maxHeight)
		if u == "" {
			continue
		}
		urls = append(urls, u)
	}

	return urls, nil
}

func best(images goinsta.Images, maxWidth, maxHeight int) string {
	fmt.Printf("%#v\n", images)

	if len(images.Versions) == 0 {
		return ""
	}

	ret := images.Versions[0].URL

	for _, c := range images.Versions {
		if c.Width <= maxWidth && c.Height <= maxHeight {
			return c.URL
		}
	}

	return ret
}

func oneOf(ss []string) string {
	if len(ss) == 0 {
		return ""
	}

	return ss[rand.Intn(len(ss)-1)]
}

func main() {
	i, err := goinsta.Import(".goinsta")
	if err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()
	router.GET("/instagram/user/:username/10.jpg", func(c *gin.Context) {
		urls, err := getLatestPicturesFromUser(i, c.Param("username"), 10)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, oneOf(urls))
	})

	router.GET("/instagram/tag/:tag/10.jpg", func(c *gin.Context) {
		urls, err := getLatestPicturesFromTag(i, c.Param("tag"), 100, 1920, 1920)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, oneOf(urls))
	})

	router.GET("/instagram/tag/:tag/10.gif", func(c *gin.Context) {
		urls, err := getLatestPicturesFromTag(i, c.Param("tag"), 10, 640, 640)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		gif, err := gif.MakeGIFFromURLs(urls, gif.MedianCut{})
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "image/gif", gif)
	})

	router.Run()
}
