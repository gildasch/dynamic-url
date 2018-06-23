package gif

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"time"

	"github.com/andybons/gogif"
	"github.com/esimov/colorquant"
)

type FloydSteinberg struct{}

var floydSteinberg = colorquant.Dither{
	[][]float32{
		[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
		[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
		[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
	},
}

func (FloydSteinberg) Convert(src image.Image, bounds image.Rectangle, p color.Palette) *image.Paletted {
	startQuant := time.Now()
	palettedImage := floydSteinberg.Quantize(src, image.NewPaletted(bounds, palette.WebSafe), 256, true, true)
	fmt.Println("FloydSteinberg.Quantize:", time.Since(startQuant))

	return palettedImage.(*image.Paletted)
}

type Sierra2 struct{}

var sierra2 = colorquant.Dither{
	[][]float32{
		[]float32{0.0, 0.0, 0.0, 4.0 / 16.0, 3.0 / 16.0},
		[]float32{1.0 / 16.0, 2.0 / 16.0, 3.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0},
		[]float32{0.0, 0.0, 0.0, 0.0, 0.0},
	},
}

func (Sierra2) Convert(src image.Image, bounds image.Rectangle, p color.Palette) *image.Paletted {
	startQuant := time.Now()
	palettedImage := sierra2.Quantize(src, image.NewPaletted(bounds, palette.WebSafe), 256, true, true)
	fmt.Println("Sierra2.Quantize:", time.Since(startQuant))

	return palettedImage.(*image.Paletted)
}

type MedianCut struct{}

func (MedianCut) Convert(src image.Image, bounds image.Rectangle, p color.Palette) *image.Paletted {
	palettedImage := image.NewPaletted(bounds, nil)

	start := time.Now()
	quantizer := gogif.MedianCutQuantizer{NumColor: 256}
	quantizer.Quantize(palettedImage, bounds, src, image.ZP)
	fmt.Println("MedianCutQuantizer:", time.Since(start))

	return palettedImage
}
