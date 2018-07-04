package utils

import (
	"image"
	"image/draw"

	"github.com/fogleman/gg"
)

const (
	TOP = iota
	MIDDLE
	BOTTOM
)

const DefaultFontSize = 70.0
const MinFontSize = 20
const FontFamily = "LiberationSans-Regular.ttf"
const Padding = 15
const DefaultPos = BOTTOM
const TextHeightPercent = 0.3

func WithCaption(in image.Image, caption string) image.Image {
	w, h := in.Bounds().Dx(), in.Bounds().Dy()

	c := gg.NewContext(w, h)
	c.DrawImage(in, 0, 0)

	fontSize := DefaultFontSize
	c.LoadFontFace(FontFamily, fontSize)

	textHeight, fontSize := adjustFontSize(c, caption, fontSize, float64(w), float64(h))

	pos := DefaultPos
	if textHeight > float64(h)/2 {
		pos = MIDDLE
	}

	drawStroke(c, caption, fontSize, float64(w), float64(h), pos)
	drawText(c, caption, fontSize, float64(w), float64(h), pos)

	if paletted, ok := in.(*image.Paletted); ok {
		out := image.NewPaletted(c.Image().Bounds(), paletted.Palette)
		draw.Draw(out, out.Bounds(), c.Image(), c.Image().Bounds().Min, draw.Src)
		return out
	}

	return c.Image()
}

func adjustFontSize(c *gg.Context, text string, fontSize, w, h float64) (txtHeight, newFontSize float64) {
	c.LoadFontFace(FontFamily, fontSize)
	lines := c.WordWrap(text, w)

	if text == "" {
		return textHeight(len(lines), fontSize), fontSize
	}

	for textHeight(len(lines), fontSize) > h*TextHeightPercent && fontSize > MinFontSize {
		fontSize--
		c.LoadFontFace(FontFamily, fontSize)
		lines = c.WordWrap(text, w)
	}

	for textHeight(len(lines), fontSize) < h*TextHeightPercent {
		fontSize++
		c.LoadFontFace(FontFamily, fontSize)
		lines = c.WordWrap(text, w)
	}

	return textHeight(len(lines), fontSize), fontSize
}

func textHeight(nlines int, fontSize float64) float64 {
	return float64(nlines) * 1.5 * fontSize
}

func drawStroke(c *gg.Context, text string, fontSize, w, h float64, position int) {
	c.SetHexColor("#000")
	px, py, ax, ay := getPxPyAxAy(position, w, h)
	n := 6 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := px + float64(dx)
			y := py + float64(dy)
			c.DrawStringWrapped(text, x, y, ax, ay, float64(w-10), 1.5, gg.AlignCenter)
		}
	}
}

func drawText(c *gg.Context, text string, fontSize, w, h float64, position int) {
	c.SetHexColor("#FFF")
	px, py, ax, ay := getPxPyAxAy(position, w, h)
	c.DrawStringWrapped(text, px, py, ax, ay,
		float64(w-10), 1.5, gg.AlignCenter)
}

func getPxPyAxAy(position int, w, h float64) (px float64, py float64,
	ax float64, ay float64) {
	px = w / 2
	ax = 0.5
	switch position {
	case TOP:
		py = Padding
		ay = 0.0
	case MIDDLE:
		py = h / 2
		ay = 0.5
	case BOTTOM:
		py = h - 2*Padding
		ay = 1.0
	}

	return
}
