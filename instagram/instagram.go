package instagram

import (
	"bufio"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

type Client struct {
	*goinsta.Instagram
}

func NewClient(path string, login bool) (*Client, error) {
	if login {
		err := loginAndExport(path)
		if err != nil {
			return nil, errors.Wrap(err, "error logging in")
		}
	}

	i, err := goinsta.Import(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading login config file %q", path)
	}

	return &Client{i}, nil
}

func loginAndExport(path string) error {
	fmt.Println("Creating the login config file for Instagram.")
	fmt.Printf("Warning: it will override any existing %s file.\n", path)

	// The following login sequence is largely inspired from
	// goinsta.v2/utils/New function

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	l, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	user := string(l)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return err
	}

	inst := goinsta.New(user, string(pass))
	err = inst.Login()
	if err != nil {
		return err
	}
	if inst.Account == nil {
		return errors.New("login failed")
	}

	err = inst.Export(path)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetLatestPicturesFromUser(username string, n int) ([]string, error) {
	user, err := c.Profiles.ByName(username)
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

func (c *Client) GetLatestPicturesFromTag(tag string, n int, maxWidth, maxHeight int) ([]string, error) {
	feed, err := c.Feed.Tags(tag)
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
		if len(urls) == n {
			return urls, nil
		}
	}

	return urls, nil
}

func best(images goinsta.Images, maxWidth, maxHeight int) string {
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
