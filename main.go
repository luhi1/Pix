package main

import (
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
)

type Pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func main() {
	reader, err := os.Open("tester.png")
	if err != nil {
		print("Error opening file.")
		panic(err)
	}
	defer reader.Close()

	m, errStr, _ := image.Decode(reader)
	if errStr != "" {
		print("Error decoding file.")
		panic(err)
	}
	bounds := m.Bounds()

	var pixArray [][]Pixel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			nuevo := Pixel{uint8(r), uint8(g), uint8(b), uint8(a)}
			pixArray[y] = append(pixArray[y], nuevo)
		}
	}

	origin := image.Point{}
	last := image.Point{X: bounds.Max.X, Y: bounds.Max.Y}
	img := image.NewRGBA(image.Rectangle{Min: origin, Max: last})
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{R: pixArray[x][y].r, G: pixArray[x][y].g, B: pixArray[x][y].b, A: pixArray[x][y].a})
		}
	}

	output, _ := os.Create("output.png")
	finErr := png.Encode(output, img)
	if finErr != nil {
		return
	}
}
