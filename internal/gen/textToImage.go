package gen

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	ErrBadFont     = errors.New("bad font")
	ErrInvalidFile = errors.New("invalid file")
	ErrBadEncoding = errors.New("bad encoding")
)

type Config struct {
	Width  int
	Height int
}

func TextToImage(text string, c *Config) (string, error) {
	// Load the font file
	fontFname := path.Join(".", "asset", "Lora-VariableFont_wght.ttf")
	fntBytes, err := os.ReadFile(fontFname)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadFont, err)
	}

	// Parse the font
	fnt, err := opentype.Parse(fntBytes)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadFont, err)
	}

	// Set up the image
	rect := image.Rect(0, 0, c.Width, c.Height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Set up the font face
	fntSize := 32
	fntFace, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(fntSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadFont, err)
	}

	x := int(math.Round(0.1 * float64(img.Bounds().Max.X)))
	y := int(math.Round(0.1 * float64(img.Bounds().Max.Y)))
	maxw := int(math.Round(0.9 * float64(img.Bounds().Max.X)))
	maxh := int(math.Round(0.9 * float64(img.Bounds().Max.Y)))

	for _, line := range splitText(text, fntFace, maxw) {

		dr := font.Drawer{
			Dst:  img,
			Src:  image.White,
			Face: fntFace,
			Dot:  fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)},
		}

		dr.DrawString(line)

		y += int(fntFace.Metrics().Height / 64)

		if y > maxh {
			break
		}
	}

	// Write the image to file
	dir := path.Join(os.TempDir(), "ttb")
	os.MkdirAll(dir, 0777)
	file, err := os.CreateTemp(dir, "output_*.png")
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidFile, err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadEncoding, err)
	}

	return file.Name(), nil
}

func splitText(text string, fntFace font.Face, maxWidth int) []string {
	var lines []string
	words := strings.Fields(text)
	line := ""
	for _, word := range words {
		if len(line)+len(word)+1 > maxWidth/fntFace.Metrics().XHeight.Round() {
			lines = append(lines, line)
			line = word
		} else {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		}
	}
	lines = append(lines, line)
	return lines
}
