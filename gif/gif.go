package gif

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"net/http"

	"github.com/andybons/gogif"
)

func MakeGIFFromURLs(urls []string) ([]byte, error) {
	subImages, err := fetchImages(urls)
	if err != nil {
		return nil, err
	}

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
		// simage, err := cutter.Crop(simage, cutter.Config{
		// 	Width:  250,
		// 	Height: 250,
		// 	Mode:   cutter.Centered,
		// })
		// if err != nil {
		// 	return nil, err
		// }

		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 100)
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, outGif)
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
