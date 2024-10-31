package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type ImageData struct {
	Path string
	Img  image.Image
}

var (
	imageData [3]*ImageData
	progress  *widget.ProgressBar
	mu        sync.Mutex
)

func main() {
	a := app.New()
	w := a.NewWindow("HDR Image Processor")

	// Create image selection buttons
	imageButtons := make([]*widget.Button, 3)
	for i := 0; i < 3; i++ {
		index := i
		imageButtons[i] = widget.NewButton("Select Image", func() {
			openImageFileDialog(w, index)
		})
	}

	// Create process button
	processButton := widget.NewButton("Process", func() {
		go processImages(w)
	})

	// Create progress bar
	progress = widget.NewProgressBar()

	// Create image preview areas
	imagePreviews := make([]*canvas.Image, 3)
	for i := 0; i < 3; i++ {
		imagePreviews[i] = canvas.NewImageFromImage(nil)
		imagePreviews[i].FillMode = canvas.ImageFillContain
	}

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Select 3 images for HDR processing:"),
		container.NewHBox(imageButtons[0], imagePreviews[0]),
		container.NewHBox(imageButtons[1], imagePreviews[1]),
		container.NewHBox(imageButtons[2], imagePreviews[2]),
		processButton,
		progress,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func openImageFileDialog(w fyne.Window, index int) {
	dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			return
		}

		img, _, err := image.Decode(reader)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		mu.Lock()
		imageData[index] = &ImageData{
			Path: reader.URI().Path(),
			Img:  img,
		}
		mu.Unlock()

		reader.Close()
	}, w).Show()
}

func processImages(w fyne.Window) {
	mu.Lock()
	defer mu.Unlock()

	// Check if all images are selected
	for i := 0; i < 3; i++ {
		if imageData[i] == nil {
			dialog.ShowError(fyne.NewError("Error", "Please select all 3 images"), w)
			return
		}
	}

	// Simulate image processing
	progress.SetValue(0)
	for i := 0; i <= 100; i++ {
		progress.SetValue(float64(i) / 100)
	}

	// Save the resulting image
	outputPath := filepath.Join(os.TempDir(), "hdr_result.jpg")
	err := saveImage(outputPath)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	dialog.ShowInformation("Success", "HDR image saved to "+outputPath, w)
}

func saveImage(outputPath string) error {
	// Simulate saving the image
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	for y := 0; y < 600; y++ {
		for x := 0; x < 800; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	return nil
}
