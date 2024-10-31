package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/harperreed/hdarrrr/internal/processor"
	"github.com/harperreed/hdarrrr/pkg/align"
	"github.com/harperreed/hdarrrr/pkg/imaging"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"
)

// convertToHDR converts a regular image to HDR format
func convertToHDR(img image.Image) hdr.Image {
	bounds := img.Bounds()
	hdrImg := hdr.NewRGB(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			hdrImg.Set(x, y, hdrcolor.RGB{
				R: float64(r) / 0xffff,
				G: float64(g) / 0xffff,
				B: float64(b) / 0xffff,
			})
		}
	}

	return hdrImg
}

func main() {
	// Define command line flags
	img1Path := flag.String("low", "", "Path to low exposure image (required)")
	img2Path := flag.String("mid", "", "Path to mid exposure image (required)")
	img3Path := flag.String("high", "", "Path to high exposure image (required)")
	outputPath := flag.String("output", "hdr_output.jpg", "Path for output HDR image")
	tonemapperFlag := flag.String("tonemapper", "drago03", "Tone mapping operator (linear, logarithmic, drago03, durand, custom_reinhard05, reinhard05, icam06)")
	gammaFlag := flag.Float64("gamma", 1.0, "Gamma correction value")
	intensityFlag := flag.Float64("intensity", 1.0, "Intensity adjustment")
	lightFlag := flag.Float64("light", 0.0, "Light adaptation (Reinhard05 only)")

	// Parse command line arguments
	flag.Parse()

	// Validate required arguments
	if *img1Path == "" || *img2Path == "" || *img3Path == "" {
		fmt.Println("Error: All three exposure images are required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate tone mapping operator
	validToneMappers := map[string]bool{
		"linear":            true,
		"logarithmic":       true,
		"drago03":           true,
		"durand":            true,
		"custom_reinhard05": true,
		"reinhard05":        true,
		"icam06":            true,
	}
	if !validToneMappers[*tonemapperFlag] {
		fmt.Printf("Error: Invalid tone mapping operator '%s'\n", *tonemapperFlag)
		fmt.Println("Valid options are: linear, logarithmic, drago03, durand, custom_reinhard05, reinhard05, icam06")
		os.Exit(1)
	}

	// Validate gamma value
	if *gammaFlag <= 0 {
		fmt.Println("Error: Gamma value must be a positive number")
		os.Exit(1)
	}

	// Load images
	images, err := imaging.LoadImages(*img1Path, *img2Path, *img3Path)
	if err != nil {
		log.Fatal("Error loading images:", err)
	}

	// Convert regular images to HDR for alignment
	hdrImages := make([]hdr.Image, len(images))
	for i, img := range images {
		hdrImages[i] = convertToHDR(img)
	}

	// Align HDR images
	alignedImages, err := align.AlignImages(images) // Align regular images first
	if err != nil {
		log.Printf("Warning: Image alignment failed: %v", err)
		log.Println("Proceeding with unaligned images...")
		alignedImages = images
	}

	// Convert aligned images to HDR
	alignedHDRImages := make([]hdr.Image, len(alignedImages))
	for i, img := range alignedImages {
		alignedHDRImages[i] = convertToHDR(img)
	}

	// Create HDR processor with configured parameters
	hdrProc := processor.NewHDRProcessor().
		WithToneMapper(*tonemapperFlag).
		WithParams(map[string]float64{
			"gamma":     *gammaFlag,
			"intensity": *intensityFlag,
			"light":     *lightFlag,
		})

	// Process HDR image
	output, err := hdrProc.Process(alignedImages)
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
	fmt.Printf("Processing parameters:\n")
	fmt.Printf("- Tone mapper: %s\n", *tonemapperFlag)
	fmt.Printf("- Gamma: %.2f\n", *gammaFlag)
	fmt.Printf("- Intensity: %.2f\n", *intensityFlag)
	if *tonemapperFlag == "reinhard05" {
		fmt.Printf("- Light adaptation: %.2f\n", *lightFlag)
	}
}
