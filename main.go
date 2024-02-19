package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
)

type Pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type Node struct {
	data     Pixel
	next     *Node
	previous *Node
}

type LinkedList struct {
	head *Node
	tail *Node
}

func (l *LinkedList) Append(newNode *Node) {

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.previous = l.tail
		l.tail.next = newNode
		l.tail = newNode
	}
}

// Function to be used for User Input/Config
func main() {
	//Create GUI
	a := app.New()

	w := a.NewWindow("Pix")

	//Container + Image + Window Setup
	var cont *fyne.Container
	var image *canvas.Image

	//This will be the path of the file the user uploads!
	var imageURI string

	//This will be used later to sort the image
	var errRange uint8
	fillerText := canvas.NewText("", color.White)
	insertBtn := widget.NewButton("1) Insert a File to Be Sorted", func() {
		//Read File Selection
		fileDialog := dialog.NewFileOpen(
			func(uc fyne.URIReadCloser, e error) {
				cont.Remove(fillerText)
				cont.Remove(image)
				//Check to make sure user selected something
				if e != nil || uc == nil {
					fillerText = canvas.NewText("Error selecting file", color.White)
					cont.Add(fillerText)
					return
				}

				//Read File Info
				data, err := ioutil.ReadAll(uc)
				if err != nil {
					fillerText = canvas.NewText("Error selecting file", color.White)
					cont.Add(fillerText)
					return
				}
				imageURI = uc.URI().Path()
				res := fyne.NewStaticResource(uc.URI().Name(), data)

				//Create Image widget and add to container
				image = canvas.NewImageFromResource(res)
				image.ScaleMode = canvas.ImageScaleFastest
				image.FillMode = canvas.ImageFillContain
				image.SetMinSize(fyne.NewSize(600, 600))
				cont.Add(image)
			}, w)

		//only allow png files
		fileDialog.SetFilter(
			storage.NewExtensionFileFilter([]string{".png"}))
		fileDialog.Show()
	})
	sortBtn := widget.NewButton("3) Sort Image by Color", func() {
		cont.Remove(fillerText)
		if imageURI == "" {
			fillerText = canvas.NewText("No File Selected!", color.White)
			cont.Add(fillerText)
			return
		}
		//Read file
		reader, err := os.Open(imageURI)

		defer func(reader *os.File) {
			err := reader.Close()
			if err != nil {
				fillerText = canvas.NewText("Error reading file", color.White)
				cont.Add(fillerText)
				return
			}
		}(reader)

		//Decode file into usable format
		m, err := png.Decode(reader)
		if err != nil {
			fillerText = canvas.NewText("Error decoding file", color.White)
			cont.Add(fillerText)
			return
		}

		pixArray := decodeImage(m)
		if err != nil {
			fillerText = canvas.NewText("Error decoding image", color.White)
			cont.Add(fillerText)
			return
		}

		//This is either an errRange with an invalid number or one that was never filled out
		if errRange == 0 {
			errRange = 15
		}
		pixArray = processImagePixels(m, pixArray, errRange)
		if pixArray == nil {
			fillerText = canvas.NewText("Error processing file", color.White)
			cont.Add(fillerText)
			return
		}
		m = encodeImage(m, pixArray)
		if m == nil {
			fillerText = canvas.NewText("Error encoding image", color.White)
			cont.Add(fillerText)
			return
		}

		//Create a new window for image
		//Sort Image
		//Put into that window
		newW := fyne.CurrentApp().NewWindow("Pix Output")

		image := canvas.NewImageFromImage(m)
		image.ScaleMode = canvas.ImageScaleFastest
		image.FillMode = canvas.ImageFillOriginal
		newW.SetContent(image)
		newW.CenterOnScreen()
		newW.SetFullScreen(true)
		newW.Show()

	})

	rangeBtn := widget.NewButton("2) Change the Error Range of Sorting", func() {
		errCont := container.NewVBox()
		errW := fyne.CurrentApp().NewWindow("Pix Error Range")
		rangeText := canvas.NewText("Pix uses an error range to determine how similar or how different colors can be before they are grouped together", color.White)
		rangeText1 := canvas.NewText("Insert an Error Range for sorting below", color.White)
		rangeText2 := canvas.NewText("Goes from Color Correctness ---> Compression", color.White)
		input := widget.NewEntry()
		input.SetPlaceHolder("Minimum 1, Maximum 255")
		saveBtn := widget.NewButton("Save", func() {
			temp, err := strconv.Atoi(input.Text)
			if err != nil || temp == 0 || temp > 255 {
				cont.Remove(fillerText)
				fillerText = canvas.NewText("Invalid error range!", color.White)
				cont.Add(fillerText)
				return
			}
			errRange = uint8(temp)
			errW.Close()
		})
		errCont = container.NewVBox(
			rangeText,
			rangeText1,
			rangeText2,
			input,
			saveBtn,
		)
		errW.SetContent(errCont)
		errW.Resize(fyne.NewSize(75, 150))
		errW.Show()
	})

	//Populate container
	cont = container.NewVBox(
		insertBtn,
		rangeBtn,
		sortBtn,
		fillerText,
	)

	//Display GUI and run
	w.SetContent(cont)
	w.ShowAndRun()
}

// Sorting pixels into color groups
func processImagePixels(m image.Image, pixArray [][]Pixel, errRange uint8) [][]Pixel {
	pixArray = sortImagePixels(m, pixArray, errRange)
	return pixArray
}

// Sort pixels into arrays of similar pixels (to reduce load on sorting algorithm)
// Merge back into an 2D array (could be just a 1D array, but 2D makes for easier proccessing!)
func sortImagePixels(m image.Image, pixArray [][]Pixel, errRange uint8) [][]Pixel {
	if pixArray == nil {
		return pixArray
	}

	var temp []Pixel
	for i := 0; i < len(pixArray); i++ {
		for j := 0; j < len(pixArray[i]); j++ {
			temp = append(temp, pixArray[i][j])
		}
	}

	startingPixel := &Node{temp[0], nil, nil}
	startingPoints := []*Node{startingPixel}
	sortedList := LinkedList{startingPixel, startingPixel}

	for i := 1; i < len(temp); i++ {
		pixel := &Node{temp[i], nil, nil}

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
				if j < len(startingPoints)-1 {
					pixel.next = startingPoints[j+1]
					pixel.previous = startingPoints[j+1].previous
					startingPoints[j+1].previous.next = pixel
					startingPoints[j+1].previous = pixel
					break
				} else {
					sortedList.Append(pixel)
					break
				}
			}

			if j == len(startingPoints)-1 {
				sortedList.Append(pixel)
				startingPoints = append(startingPoints, pixel)

				break
			}
		}
	}

	u := sortedList.head
	for i := 0; i < len(pixArray); i++ {
		for j := 0; j < len(pixArray[i]); j++ {
			pixArray[i][j] = u.data
			u = u.next
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
func encodeImage(m image.Image, pixArray [][]Pixel) image.Image {
	os.Remove("output.png")
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

	return img
}
