package tarsier

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	//"image/draw"
	"image/jpeg"
	"golang.org/x/image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"

	pkgFont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
)

func Gen() {
	//img, _ := image.Decode(os.Stdin)
	img, _, err := image.Decode(os.Stdin)

	if err != nil {
		return
	}

	bounds := img.Bounds()
	dest := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(img.At(x, y))
			gray, _ := c.(color.Gray16)
			dest.Set(x, y, gray)
		}
	}
	png.Encode(os.Stdout, dest)
}

type Option struct {
	FontSize float64
}

type Generator struct {
	Option Option
}

func NewGenerator(opt Option) *Generator {
	return &Generator{Option: opt}
}

func (g *Generator) Generate(source io.Reader, dest io.Writer) error {
	font, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return err
	}

	fontOpt := truetype.Options{
		Size:              g.Option.FontSize,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face := truetype.NewFace(font, &fontOpt)
	defer face.Close()

	srcImg, _, err := image.Decode(source)
	if err != nil {
		return err
	}

	dstImg := image.NewRGBA(image.Rect(0, 0, srcImg.Bounds().Dx(), srcImg.Bounds().Dy()))
	draw.Draw(dstImg, dstImg.Rect, srcImg, image.Point{}, draw.Src)
	dot := fixed.P(0, 0)
	fontDrawer := pkgFont.Drawer{Dst: dstImg, Src: image.White, Face: face, Dot: dot}

	text := "LGTM"

	bounds, _ := fontDrawer.BoundString(text)
	fmt.Printf("%+v\n", bounds)

	dx := int((bounds.Max.X - bounds.Min.X) / 60)
	dy := int((bounds.Max.Y - bounds.Min.Y) / 60)

	fmt.Printf("%v, %v\n", dx, dy)

	x := (srcImg.Bounds().Dx() - dx) / 2
	y := (srcImg.Bounds().Dy() - dy) / 2

	fontDrawer.Dot = fixed.P(x, y)
	fmt.Printf("%+v", fontDrawer.Dot)
	fontDrawer.DrawString(text)

	o := jpeg.Options{Quality: 100}
	err = jpeg.Encode(dest, dstImg, &o)

	//err = png.Encode(os.Stdout, dstImg)
	if err != nil {
		return err
	}

	return nil
}

func GenA(source io.Reader) error {

	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return err
	}

	srcImg, _, err := image.Decode(source)
	if err != nil {
		return err
	}

	// 0はデフォルト値
	opt := truetype.Options{
		Size:              110,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	//imageWidth := 500
	//imageHeight := 10s0
	//textTopMargin := 100
	text := "LGTM"

	bounds := srcImg.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, srcImg.Bounds().Dx(), srcImg.Bounds().Dy()))
	draw.Draw(img, img.Bounds(), srcImg, image.Pt(0, 0), draw.Src)
	//draw.CatmullRom.Scale(img, img.Bounds(), srcImg, bounds, draw.Over, nil)

	face := truetype.NewFace(ft, &opt)
	//image.Pt(10, 40)
	dr := &pkgFont.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	// ど真ん中
	// 横が動いた
	dr.Dot.X = (fixed.I(bounds.Dx()) - dr.MeasureString(text)) / 2
	// 縦が動いた
	dr.Dot.Y = fixed.I(bounds.Dy() + 250) / 2

	dr.DrawString(text)

	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.Create(`test.png`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(buf.Bytes())

	return nil
}
