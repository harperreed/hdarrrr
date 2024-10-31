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
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"

	// Import HDR codecs
	_ "github.com/mdouchement/hdr/codec/rgbe"
	_ "github.com/mdouchement/tiff"
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

// convertToHDRSlice converts a slice of regular images to HDR images
func convertToHDRSlice(images []image.Image) []hdr.Image {
	hdrImages := make([]hdr.Image, len(images))
	for i, img := range images {
		hdrImages[i] = convertToHDR(img)
	}
	return hdrImages
}

// convertToRegularSlice converts a slice of HDR images to regular images
func convertToRegularSlice(images []hdr.Image) []image.Image {
	regularImages := make([]image.Image, len(images))
	for i, img := range images {
		regularImages[i] = img.(image.Image)
	}
	return regularImages
}

func main() {
	// Define command line flags
	img1Path := flag.String("low", "", "Path to low exposure image (required)")
	img2Path := flag.String("mid", "", "Path to mid exposure image (required)")
	img3Path := flag.String("high", "", "Path to high exposure image (required)")
	outputPath := flag.String("output", "hdr_output.jpg", "Path for output HDR image")
	tonemapperFlag := flag.String("tonemapper", "reinhard05", "Tone mapping operator (reinhard05, drago03)")
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

	// Validate file extensions
	paths := []string{*img1Path, *img2Path, *img3Path}
	for _, path := range paths {
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

	// Convert to HDR for alignment
	hdrImages := convertToHDRSlice(images)

	// Align images
	alignedHDRImages, err := align.AlignImages(hdrImages)
	if err != nil {
		log.Printf("Warning: Image alignment failed: %v", err)
		log.Println("Proceeding with unaligned images...")
		alignedHDRImages = hdrImages // Use original images if alignment fails
	}

	// Convert back to regular images for processing
	alignedImages := convertToRegularSlice(alignedHDRImages)

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
