package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

func getLatestPicture(i *goinsta.Instagram, username string) (string, error) {
	user, err := i.Profiles.ByName(username)
	if err != nil {
		return "", err
	}

	media := user.Feed()

	if !media.Next() {
		return "", errors.Errorf("user %q have no post", username)
	}

	if len(media.Items) <= 0 {
		return "", errors.Errorf("user %q have no metia.Item", username)
	}

	if len(media.Items[0].Images.Versions) <= 0 {
		return "", errors.Errorf("user %q have no metia.Items[0].Images.Versions", username)
	}

	return media.Items[0].Images.Versions[0].URL, nil
}

func main() {
	i, err := goinsta.Import(".goinsta")
	if err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()
	router.GET("/instagram/:username/first.jpg", func(c *gin.Context) {
		u, err := getLatestPicture(i, c.Param("username"))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, u)
	})

	router.Run()
}
