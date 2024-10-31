package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/harperreed/hdarrrr/internal/processor"
	"github.com/harperreed/hdarrrr/pkg/align"
	"github.com/harperreed/hdarrrr/pkg/imaging"
)

func main() {
	// Define command line flags
	img1Path := flag.String("low", "", "Path to low exposure image (required)")
	img2Path := flag.String("mid", "", "Path to mid exposure image (required)")
	img3Path := flag.String("high", "", "Path to high exposure image (required)")
	outputPath := flag.String("output", "hdr_output.jpg", "Path for output HDR image")
	hdrMethod := flag.String("method", "tone-mapping", "HDR method: tone-mapping or exposure-fusion")

	// Parse command line arguments
	flag.Parse()

	// Validate required arguments
	if *img1Path == "" || *img2Path == "" || *img3Path == "" {
		fmt.Println("Error: All three exposure images are required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate file extensions
	for _, path := range []string{*img1Path, *img2Path, *img3Path} {
		ext := strings.ToLower(filepath.Ext(path))
		if !imaging.SupportedFormats[ext] {
			log.Fatalf("Error: Unsupported image format for file %s. Supported formats: PNG, JPEG", path)
		}
	}

	// Load images
	images, err := imaging.LoadImages(*img1Path, *img2Path, *img3Path)
	if err != nil {
		log.Fatal("Error loading images:", err)
	}

	// Align images
	alignedImages, err := align.AlignImages(images)
	if err != nil {
		log.Printf("Warning: Image alignment failed: %v", err)
		alignedImages = images // Use original images if alignment fails
	}

	// Validate image properties
	if err := validateImageProperties(alignedImages); err != nil {
		log.Fatal("Error validating image properties:", err)
	}

	// Process HDR
	var output image.Image
	switch *hdrMethod {
	case "tone-mapping":
		hdrProcessor := processor.NewHDRProcessor()
		output, err = hdrProcessor.Process(alignedImages)
	case "exposure-fusion":
		output, err = processExposureFusion(alignedImages)
	default:
		log.Fatalf("Error: Unsupported HDR method %s. Supported methods: tone-mapping, exposure-fusion", *hdrMethod)
	}
	if err != nil {
		log.Fatal("Error processing HDR:", err)
	}

	// Create output directory if it doesn't exist
	if dir := filepath.Dir(*outputPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("Error creating output directory:", err)
		}
	}

	// Save the result
	if err := imaging.SaveImage(output, *outputPath); err != nil {
		log.Fatal("Error saving output image:", err)
	}

	fmt.Printf("HDR image successfully saved to %s\n", *outputPath)
}

func validateImageProperties(images []image.Image) error {
	if len(images) < 2 {
		return fmt.Errorf("at least two images are required for validation")
	}

	// Check for nil images first
	for i, img := range images {
		if img == nil {
			return fmt.Errorf("image %d is nil", i+1)
		}
	}

	baseImg := images[0]
	baseProps := imaging.GetImageProperties(baseImg, filepath.Ext(baseImg.Bounds().String()))
	baseColorModel := baseImg.ColorModel()

	for i, img := range images[1:] {
		currentProps := imaging.GetImageProperties(img, filepath.Ext(img.Bounds().String()))

		// Check dimensions
		if img.Bounds() != baseImg.Bounds() {
			return fmt.Errorf("image %d has different dimensions than the first image", i+2)
		}

		// Check color model
		currentColorModel := img.ColorModel()
		if currentColorModel != baseColorModel {
			return fmt.Errorf("image %d has a different color model (%T) than the first image (%T)",
				i+2, currentColorModel, baseColorModel)
		}

		// Compare other properties
		if !imaging.ValidateImageProperties(baseProps, currentProps) {
			return fmt.Errorf("image %d has different properties than the first image", i+2)
		}
	}

	return nil
}

func processExposureFusion(images []image.Image) (image.Image, error) {
	// Placeholder for the actual implementation of the MKVR algorithm
	// This function should implement the Mertens-Kautz-Van Reeth (MKVR) algorithm for exposure fusion
	return nil, fmt.Errorf("exposure fusion method not yet implemented")
}
