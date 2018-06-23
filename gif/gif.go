package gif

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"net/http"
	"time"

	"github.com/esimov/colorquant"
)

var floydSteinberg = colorquant.Dither{
	[][]float32{
		[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
		[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
		[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
	},
}

func MakeGIFFromURLs(urls []string) ([]byte, error) {
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
		startQuant := time.Now()
		palettedImage := floydSteinberg.Quantize(simage, image.NewPaletted(bounds, palette.WebSafe), 256, true, true)
		fmt.Println("floydSteinberg.Quantize:", time.Since(startQuant))

		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, palettedImage.(*image.Paletted))
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
