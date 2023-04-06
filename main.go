package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ppicom/ttb/gen"
)

func main() {

	w := flag.Int("w", 400, "Width of the picture")
	h := flag.Int("h", 400, "Hight of the picture")
	flag.Parse()

	fname, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fname = path.Join(fname, "Library", "Fonts", "Lora-VariableFont_wght.ttf")

	outputName, err := gen.TextToImage(strings.Join(flag.Args(), " "), &gen.Config{
		Width:  *w,
		Height: *h,
	})
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, outputName)
}
