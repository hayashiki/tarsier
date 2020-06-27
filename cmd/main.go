package main

import (
	"os"
	"tarsier"
)

func main() {

	//gen := tarsier.NewGenerator(tarsier.Option{11})
	//err := gen.Generate(os.Stdin, os.Stdout)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	os.Exit(1)
	//}
	tarsier.GenA(os.Stdin)
}
