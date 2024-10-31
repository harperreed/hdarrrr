package main

import (
	"image"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestOpenImageFileDialog(t *testing.T) {
	a := test.NewApp()
	w := a.NewWindow("Test")

	// Simulate opening an image file dialog
	openImageFileDialog(w, 0)

	// Check if the image data is set correctly
	if imageData[0] == nil {
		t.Error("Expected image data to be set, but it is nil")
	}
}

func TestProcessImages(t *testing.T) {
	a := test.NewApp()
	w := a.NewWindow("Test")

	// Set up dummy image data
	imageData[0] = &ImageData{Img: image.NewRGBA(image.Rect(0, 0, 100, 100))}
	imageData[1] = &ImageData{Img: image.NewRGBA(image.Rect(0, 0, 100, 100))}
	imageData[2] = &ImageData{Img: image.NewRGBA(image.Rect(0, 0, 100, 100))}

	// Simulate processing images
	go processImages(w)

	// Check if the progress bar is updated correctly
	if progress.Value != 1.0 {
		t.Errorf("Expected progress bar value to be 1.0, but got %f", progress.Value)
	}
}
