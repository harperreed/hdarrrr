package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/harperreed/hdarrrr/internal/processor"
	"github.com/harperreed/hdarrrr/pkg/imaging"
)

func main() {
	// Define command line flags
	img1Path := flag.String("low", "", "Path to low exposure image (required)")
	img2Path := flag.String("mid", "", "Path to mid exposure image (required)")
	img3Path := flag.String("high", "", "Path to high exposure image (required)")
	outputPath := flag.String("output", "hdr_output.jpg", "Path for output HDR image")

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

	// Process HDR
	hdrProcessor := processor.NewHDRProcessor()
	output, err := hdrProcessor.Process(images)
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