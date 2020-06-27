package main

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
	pkgFont "golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
	"path/filepath"

	//"image/jpeg"
	"image/png"
	"log"
	"os"
)

const OVERLAY_TEXT = "LGTM"

// https://qiita.com/tanksuzuki/items/7866768c36e13f09eedb
func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func main() {
	imagePaths := dirwalk("./input")

	for _, imagePath := range imagePaths {
		file, err := os.OpenFile(imagePath, os.O_RDONLY, os.ModePerm)

		log.Printf("file name is %s", file.Name())
		defer func() {
			if err := file.Close(); err != nil {
				return
			}
		}()

		img, _, err := image.Decode(file)

		if err != nil {
			log.Fatal(err)
		}

		img = run(img)

		out, err := os.Create(fmt.Sprintf("output/%s", file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		png.Encode(out, img)
	}
	//png.Encode(os.Stdout, m)
}

func overlayText(imgSrc image.Image, text string) (image.Image, error) {
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
	//image.Pt(10, 40)
	dr := &pkgFont.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = (fixed.I(bounds.Dx()) - dr.MeasureString(text)) / 2
	// 縦が動いた
	dr.Dot.Y = fixed.I(bounds.Dy()) / 2

	dr.DrawString(text)

	buf := &bytes.Buffer{}
	err = png.Encode(buf, img)

	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		//os.Exit(1)
	}

	return img, nil

}

func run(img image.Image) image.Image {
	//_, _, err = image.DecodeConfig(buf)
	processedImg := resize.Resize(512, 0, img, resize.Lanczos3)
	processedImg, err := overlayText(processedImg, OVERLAY_TEXT)
	if err != nil {
		log.Fatal(err)
	}
	return 	processedImg
}

//out, err := os.Create("test_resized.png")
//if err != nil {
//	log.Fatal(err)
//}
//defer out.Close()

// write new image to file
//png.Encode(out, m, nil)

// open "test.jpg"
//file, err := os.Open("./input/gopher.png")
//if err != nil {
//	log.Fatal(err)
//}
//
//buf := new(bytes.Buffer)
//io.Copy(buf, file)


// decode jpeg into image.Image
//buf := new(bytes.Buffer)
//buf := bytes.NewBuffer(os.Stdin)
//io.Copy(buf, os.Stdin)
