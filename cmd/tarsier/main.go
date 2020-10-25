package main

import (
	"fmt"
	"github.com/hayashiki/tarsier"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

// eg. go run cmd/main.go --in=./testdata/dog.jpg" --out="./output"
// eg. go run cmd/main.go --in=./testdata/gopher.png" --out="./output"
func main() {
	app := cli.NewApp()
	app.Version = "v1.0.0" // TODO versionうめこむ
	app.Name = "tarsier"
	app.Usage = "Put text in the image"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "in, i",
			Required: true,
			Usage:    "target input image",
		},
		cli.StringFlag{
			Name:     "out, o",
			Value:    ".",
			Usage:    "target output image",
		},
		cli.StringFlag{
			Name:     "text, t",
			Value:    tarsier.DefaultOverlayText,
			Usage:    "append text",
		},
	}

	app.Action = func(c *cli.Context) error {
		in := c.String("in")
		outFlag := c.String("out")
		text := c.String("text")

		file, err := os.OpenFile(in, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		defer func() {
			if err := file.Close(); err != nil {
				return
			}
		}()

		img, err := decode(file)
		if err != nil {
			return fmt.Errorf("could read image '%s': %v", in, err)
		}
		outImg, err := tarsier.Print(img, text)
		if err != nil {
			return fmt.Errorf("could print text '%s': %v", in, err)
		}

		// TODO: pngで出力するのでjpgの場合拡張子を変更する
		out, err := os.Create(fmt.Sprintf("%s/%s", outFlag, filepath.Base(file.Name())))
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := out.Close(); err != nil {
				return
			}
		}()

		return png.Encode(out, outImg)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func decode(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)

	return img, err
}
