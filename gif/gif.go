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

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type Converter interface {
	Convert(src image.Image, bounds image.Rectangle, p color.Palette) *image.Paletted
}

const (
	width  = 240
	height = 240
)

func MakeGIFFromURLs(urls []string, converter Converter) ([]byte, error) {
	start := time.Now()
	fetched, err := fetchImages(urls)
	fmt.Println("fetchImages:", time.Since(start))
	if err != nil {
		return nil, err
	}

	start = time.Now()
	var normalized []image.Image
	for _, f := range fetched {
		cropped, err := cutter.Crop(f, cutter.Config{
			Width:   1,
			Height:  1,
			Mode:    cutter.Centered,
			Options: cutter.Ratio,
		})
		if err != nil {
			return nil, err
		}

		resized := resize.Resize(width, height, cropped, resize.Bilinear)

		normalized = append(normalized, resized)
	}
	bounds := image.Rect(0, 0, width, height)

	outGif := &gif.GIF{}
	for _, n := range normalized {
		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, converter.Convert(n, bounds, nil))
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
