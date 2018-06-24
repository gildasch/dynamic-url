package gif

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
	defaultWidth  = 240
	defaultHeight = 240
)

func MakeGIFFromImages(in []image.Image, delay time.Duration, converter Converter) ([]byte, error) {
	outGif := &gif.GIF{}
	start := time.Now()
	for _, i := range in {
		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, converter.Convert(i, i.Bounds(), nil))
		outGif.Delay = append(outGif.Delay, int(delay.Seconds()*100)) // delay is in 100th of second
	}
	fmt.Println("appends:", time.Since(start))

	var buf bytes.Buffer
	start = time.Now()
	err := gif.EncodeAll(&buf, outGif)
	fmt.Println("gif encode:", time.Since(start))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func MakeGIFFromURLs(urls []string, delay time.Duration, converter Converter) ([]byte, error) {
	start := time.Now()
	fetched, err := fetchImages(urls)
	fmt.Println("fetchImages:", time.Since(start))
	if err != nil {
		return nil, err
	}

	start = time.Now()
	var normalized []image.Image
	for _, f := range fetched {
		n, err := normalize(f, defaultWidth, defaultHeight)
		if err != nil {
			return nil, err
		}

		normalized = append(normalized, n)
	}
	fmt.Println("normalize:", time.Since(start))

	return MakeGIFFromImages(normalized, delay, converter)
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

func normalize(in image.Image, width, height int) (image.Image, error) {
	maxWidth, maxHeight := uint(0), uint(0)
	if in.Bounds().Dx() > in.Bounds().Dy() {
		maxHeight = uint(height)
	} else {
		maxWidth = uint(width)
	}
	resized := resize.Resize(maxWidth, maxHeight, in, resize.Bilinear)

	cropped, err := cutter.Crop(resized, cutter.Config{
		Width:   1,
		Height:  1,
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})
	if err != nil {
		return nil, err
	}

	out := image.NewRGBA(image.Rect(0, 0, 240, 240))
	draw.Draw(out, out.Bounds(), cropped, cropped.Bounds().Min, draw.Src)

	return out, nil
}
