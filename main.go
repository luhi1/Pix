package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

// Function to be used for User Input/Config
func main() {
	//Intro Script
	fmt.Println("Welcome to Pix, the best pixel sorter of all time!")
	var in string
	fmt.Print("Enter the name of a file to be sorted (must be in the test directory and do not include the .png!) ")
	_, err := fmt.Scan(&in)
	if err != nil {
		return
	}

	//Read file
	reader, err := os.Open("testing/" + in + ".png")
	defer func(reader *os.File) {
		err := reader.Close()
		check(err, "Error opening file.")
	}(reader)
	check(err, "Error opening file.")

	//Decode file into usable format
	m, err := png.Decode(reader)
	check(err, "Error decoding file.")

	if err != nil {
		panic(err)
	}

	pixArray := decodeImage(m)
	pixArray = processImagePixels(pixArray)
	encodeImage(m, pixArray)
}

// Basic error handling
func check(err error, errCode string) {
	if err != nil {
		print(errCode)
	}
}

// Sorting pixels into color groups
func processImagePixels(pixArray [][]Pixel) [][]Pixel {
	var in uint8

	fmt.Print("Enter an error range for the pixel sorter: ")
	_, err := fmt.Scan(&in)
	if err != nil {
		return nil
	}

	pixArray = sortImagePixels(pixArray, in)
	return pixArray
}

// Sort pixels into arrays of similar pixels (to reduce load on sorting algorithm)
// Merge back into an 2D array (could be just a 1D array, but 2D makes for easier proccessing!)
func sortImagePixels(pixArray [][]Pixel, errRange uint8) [][]Pixel {
	var sortedArrays [][]Pixel

	for y := 0; y < len(pixArray); y++ {
		for x := 0; x < len(pixArray[y]); x++ {
			pixel := pixArray[y][x]
			if len(sortedArrays) == 0 {
				sortedArrays = append(sortedArrays, []Pixel{pixel})
				continue
			}
			for i := 0; i < len(sortedArrays); i++ {
				indexPixel := sortedArrays[i][0]
				indexPixelRedMin := indexPixel.r - errRange
				if indexPixelRedMin > indexPixel.r {
					indexPixelRedMin = 0
				}
				indexPixelGreenMin := indexPixel.g - errRange
				if indexPixelGreenMin > indexPixel.g {
					indexPixelGreenMin = 0
				}
				indexPixelBlueMin := indexPixel.b - errRange
				if indexPixelBlueMin > indexPixel.b {
					indexPixelBlueMin = 0
				}

				indexPixelRedMax := indexPixel.r + errRange
				if indexPixelRedMax < indexPixel.r {
					indexPixelRedMax = 255
				}
				indexPixelGreenMax := indexPixel.g + errRange
				if indexPixelGreenMax < indexPixel.g {
					indexPixelGreenMax = 255
				}
				indexPixelBlueMax := indexPixel.b + errRange
				if indexPixelBlueMax < indexPixel.b {
					indexPixelBlueMax = 255
				}

				if ((pixel.r <= indexPixelRedMax) && (pixel.r >= indexPixelRedMin)) &&
					((pixel.g <= indexPixelGreenMax) && (pixel.g >= indexPixelGreenMin)) &&
					((pixel.b <= indexPixelBlueMax) && (pixel.b >= indexPixelBlueMin)) {

					sortedArrays[i] = append(sortedArrays[i], pixel)
					break
				}

				if i == len(sortedArrays)-1 {
					sortedArrays = append(sortedArrays, []Pixel{pixel})
					break
				}
			}

		}
	}

	var sortedArraysOneArrayLol []Pixel
	for i := 0; i < len(sortedArrays); i++ {
		for j := 0; j < len(sortedArrays[i]); j++ {
			sortedArraysOneArrayLol = append(sortedArraysOneArrayLol, sortedArrays[i][j])
		}
	}

	pixCounter := 0
	for y := 0; y < len(pixArray); y++ {
		for x := 0; x < len(pixArray[y]); x++ {
			pixArray[y][x] = sortedArraysOneArrayLol[pixCounter]
			pixCounter++
		}
	}

	return pixArray
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
