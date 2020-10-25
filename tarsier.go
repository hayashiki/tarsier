package tarsier

import (
	"bytes"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	pkgFont "golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
)

const DefaultOverlayText = "LGTM"


func Print(img image.Image, text string) (image.Image, error) {

	processedImg := resize.Resize(512, 0, img, resize.Lanczos3)

	if text == "" {
		text = DefaultOverlayText
	}
	// TODO only resize
	processedImg, err := print(processedImg, text)
	if err != nil {
		return nil, err
	}
	return 	processedImg, err
}

func print(imgSrc image.Image, text string) (image.Image, error) {
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return nil, err
	}

	opt := truetype.Options{
		Size:              110,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	bounds := imgSrc.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(img, img.Bounds(), imgSrc, image.Pt(0, 0), draw.Src)

	face := truetype.NewFace(ft, &opt)
	dr := &pkgFont.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(bounds.Dx()) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I(bounds.Dy()) / 2

	dr.DrawString(text)

	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)

	if err != nil {
		return nil, err
	}
	return img, nil
}
