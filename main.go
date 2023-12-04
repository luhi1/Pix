package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"fmt"
)

type Pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

// Function to be used for User Input/Config
func main() {
	//Read file
	reader, err := os.Open("tester.png")
	defer reader.Close()
	check(err, "Error opening file.")

	//Decode file into usable format
	m, err := png.Decode(reader)
	check(err, "Error decoding file.")

	if err != nil {
		panic(err)
	}

	pixArray := decodeImage(m)
	pixArray = proccessImagePixels(pixArray)
	encodeImage(m, pixArray)
}

// Basic error handling
func check(err error, errCode string) {
	if err != nil {
		print(errCode)
	}
}

// Sorting pixels into color groups
func sortImagePixels(pixArray [][]Pixel) [][]Pixel {
	pixArray = proccessImagePixels(pixArray)
	return pixArray
}

// Sort pixels into arrays of similar pixels (to reduce load on sorting algorithm)
// Merge back into an 2D array (could be just a 1D array, but 2D makes for easier proccessing!)
func proccessImagePixels(pixArray [][]Pixel) [][]Pixel {
	sortedArray := [][]Pixel{}

	for y := 0; y < len(pixArray); y++ {
		for x := 0; x < len(pixArray[y]); x++ {
			pixel := pixArray[x][y]
			if len(sortedArray) == 0 {
				sortedArray = append(sortedArray, []Pixel{pixel})
			}
			for i := 0; i < len(sortedArray); i++ {
				if ((pixel.r <= sortedArray[i][0].r+15) || (pixel.r >= sortedArray[i][0].r-15)) &&
					((pixel.g <= sortedArray[i][0].g+15) || (pixel.g >= sortedArray[i][0].g-15)) &&
					((pixel.b <= sortedArray[i][0].b+15) || (pixel.b >= sortedArray[i][0].b-15)) {

					alpha := pixel.a
					lengthi := len(sortedArray[i])
					for j := 0; j < lengthi; j++ {
					println(sortedArray[i][j].r)
						if alpha >= sortedArray[i][j].a {
							sortedArray[i] = insert(sortedArray[i], j, pixel)
						}
					}
					break
				}

				if i == len(sortedArray) {
					sortedArray = append(sortedArray, []Pixel{pixel})
				}
			}

		}
	}
	return sortedArray
}

func insert(a []Pixel, index int, value Pixel) []Pixel {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	fmt.Println(a)
	return a
}

// Convert file into 2D Array of Pixels
func decodeImage(m image.Image) [][]Pixel {
	//Define dimensions of image
	bounds := m.Bounds()

	//Create array for pixels.
	var pixArray [][]Pixel

	//Populate array of pixels using data from read file.
	//Starts from top left and goes to bottom right BY ROW.
	//Top Left and Bottom Right are defined by (Min.x, Min.y) and (Max.X, Max.Y) of bounds, respectively.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		pixArray = append(pixArray, []Pixel{})
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			nuevo := Pixel{uint8(r), uint8(g), uint8(b), uint8(a)}
			pixArray[y] = append(pixArray[y], nuevo)
		}
	}

	return pixArray
}

// Convert 2D Array of Pixels into Image
func encodeImage(m image.Image, pixArray [][]Pixel) {
	//Get Dimensions of image and create a new image for output based on dimensions
	//Image is defined as a rectangle with W: Max.X and L: Max.Y based on dimensions
	bounds := m.Bounds()
	origin := image.Point{}
	last := image.Point{X: bounds.Max.X, Y: bounds.Max.Y}
	img := image.NewRGBA(image.Rectangle{Min: origin, Max: last})

	//Populate output with pixels from pixArray
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			currPix := pixArray[y][x]
			img.SetRGBA(x, y, color.RGBA{R: currPix.r, G: currPix.g, B: currPix.b, A: currPix.a})
		}
	}

	//Create file for output
	output, err := os.Create("output.png")

	//Write output image to "output.png" and encode to .png format
	err = png.Encode(output, img)
	check(err, "Error encoding image")
}
