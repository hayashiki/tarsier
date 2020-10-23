package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	publicCloud = ""
)

func main() {
	app := cli.NewApp()
	app.Version = "v1.0.0"
	app.Name = "tarisier"
	app.Usage = "for template"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "dir, d",
			Value:    "./input",
			EnvVar:   "INPUT_IMAGEPATH",
			Required: true,
			Usage:    "target input image path",
		},
		cli.StringFlag{
			Name:     "output, o",
			Value:    "./output",
			EnvVar:   "OUTPUT_IMAGEPATH",
			Required: true,
			Usage:    "target output image path",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Println("say cli")
		dirFlag := c.String("dir")
		outputflag := c.String("output")

		log.Println(dirFlag)
		log.Println(outputflag)

		// do action

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
