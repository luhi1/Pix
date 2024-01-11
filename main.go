package main

import (
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

type Node struct {
	data Pixel
	next *Node
}

type LinkedList struct {
	head *Node
}

func (l *LinkedList) Append(p Pixel) {
	list := &Node{data: p, next: nil}
	if l.head == nil {
		l.head = list
	} else {
		p := l.head
		for p.next != nil {
			p = p.next
		}
		p.next = list
	}
}

// Function to be used for User Input/Config
func main() {
	/*
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
	*/
	reader, err := os.Open("testing/tester.png")

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
	/*
		var in uint8

		fmt.Print("Enter an error range for the pixel sorter: ")
		poo, err := fmt.Scan(&in)
		if err != nil {
			return nil
		}*/

	pixArray = sortImagePixels(pixArray, 15)
	return pixArray
}

// Sort pixels into arrays of similar pixels (to reduce load on sorting algorithm)
// Merge back into an 2D array (could be just a 1D array, but 2D makes for easier proccessing!)
func sortImagePixels(pixArray [][]Pixel, errRange uint8) [][]Pixel {
	if pixArray == nil {
		return pixArray
	}

	var temp []Pixel
	for i := 0; i < len(pixArray); i++ {
		for j := 0; j < len(pixArray[i]); j++ {
			temp = append(temp, pixArray[i][j])
		}
	}

	//Starting points[0] has a different memory address from linkedList[0]!!!!!
	startingPoints := []Node{{temp[0], nil}}
	sortedList := LinkedList{&startingPoints[0]}

	for i := 1; i < len(temp); i++ {
		pixel := Node{temp[i], nil}

		for j := 0; j < len(startingPoints); j++ {
			indexPixel := startingPoints[j].data
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

			if ((pixel.data.r <= indexPixelRedMax) && (pixel.data.r >= indexPixelRedMin)) &&
				((pixel.data.g <= indexPixelGreenMax) && (pixel.data.g >= indexPixelGreenMin)) &&
				((pixel.data.b <= indexPixelBlueMax) && (pixel.data.b >= indexPixelBlueMin)) {
				pixel.next = startingPoints[j].next
				startingPoints[j].next = &pixel
				break
			}

			if j == len(startingPoints)-1 {
				temp := &startingPoints[j]
				for temp.next != nil {
					temp = temp.next
				}
				temp.next = &pixel
				startingPoints = append(startingPoints, pixel)
				break
			}
		}
	}
	//Todo: Finish LL implementation
	pixCounter := 0
	u := sortedList.head
	pixArray[0][0] = u.data
	for u.next != nil {
		u = u.next
		pixArray[pixCounter/len(pixArray)][pixCounter%len(pixArray)] = u.data
		pixCounter++
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
