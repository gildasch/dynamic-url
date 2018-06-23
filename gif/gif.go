package gif

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"net/http"
	"time"
)

type Converter interface {
	Convert(src image.Image, bounds image.Rectangle, p color.Palette) *image.Paletted
}

func MakeGIFFromURLs(urls []string, converter Converter) ([]byte, error) {
	start := time.Now()
	subImages, err := fetchImages(urls)
	fmt.Println("fetchImages:", time.Since(start))
	if err != nil {
		return nil, err
	}

	start = time.Now()
	width, height := -1, -1
	for _, i := range subImages {
		if width == -1 || width > i.Bounds().Dx() {
			width = i.Bounds().Dx()
		}
		if height == -1 || height > i.Bounds().Dy() {
			height = i.Bounds().Dy()
		}
	}
	bounds := image.Rect(0, 0, width, height)

	outGif := &gif.GIF{}
	for _, simage := range subImages {
		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, converter.Convert(simage, bounds, nil))
		outGif.Delay = append(outGif.Delay, 100)
	}
	fmt.Println("appends:", time.Since(start))

	var buf bytes.Buffer
	start = time.Now()
	err = gif.EncodeAll(&buf, outGif)
	fmt.Println("gif encode:", time.Since(start))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func fetchImages(urls []string) ([]image.Image, error) {
	var imgs []image.Image

	for _, u := range urls {
		resp, err := http.Get(u)
		if err != nil {
			fmt.Println("error fetching", u, err)
			continue
		}
		defer resp.Body.Close()

		i, err := jpeg.Decode(resp.Body)
		if err != nil {
			fmt.Println("error decoding", u, err)
			continue
		}

		imgs = append(imgs, i)
	}

	return imgs, nil
}
