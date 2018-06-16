package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

func getLatestPictures(i *goinsta.Instagram, username string, n int) ([]string, error) {
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
	router.GET("/instagram/:username/10.jpg", func(c *gin.Context) {
		urls, err := getLatestPictures(i, c.Param("username"), 10)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, oneOf(urls))
	})

	router.Run()
}
