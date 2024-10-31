package align

import (
	"errors"
	"image"
	"log"

	"gocv.io/x/gocv"
)

// AlignImages aligns multiple images using feature matching (e.g., SIFT, ORB)
func AlignImages(images []image.Image) ([]image.Image, error) {
	if len(images) < 2 {
		return nil, errors.New("at least two images are required for alignment")
	}

	// Convert images to gocv.Mat
	mats := make([]gocv.Mat, len(images))
	for i, img := range images {
		mat, err := gocv.ImageToMatRGB(img)
		if err != nil {
			return nil, err
		}
		mats[i] = mat
		defer mat.Close()
	}

	// Use the first image as the reference
	refMat := mats[0]

	// Initialize ORB detector
	orb := gocv.NewORB()
	defer orb.Close()

	// Detect keypoints and compute descriptors for the reference image
	refKeypoints, refDescriptors := orb.DetectAndCompute(refMat, gocv.NewMat())

	alignedImages := make([]image.Image, len(images))
	alignedImages[0] = images[0]

	for i := 1; i < len(images); i++ {
		// Detect keypoints and compute descriptors for the current image
		keypoints, descriptors := orb.DetectAndCompute(mats[i], gocv.NewMat())

		// Match descriptors using BFMatcher
		matcher := gocv.NewBFMatcher()
		matches := matcher.KnnMatch(refDescriptors, descriptors, 2)

		// Filter good matches using the ratio test
		goodMatches := make([]gocv.DMatch, 0)
		for _, m := range matches {
			if len(m) == 2 && m[0].Distance < 0.75*m[1].Distance {
				goodMatches = append(goodMatches, m[0])
			}
		}

		if len(goodMatches) < 4 {
			log.Printf("Warning: Not enough good matches for image %d", i+1)
			alignedImages[i] = images[i]
			continue
		}

		// Extract matched keypoints
		refPoints := make([]gocv.Point2f, len(goodMatches))
		imgPoints := make([]gocv.Point2f, len(goodMatches))
		for j, m := range goodMatches {
			refPoints[j] = refKeypoints[m.QueryIdx].Pt
			imgPoints[j] = keypoints[m.TrainIdx].Pt
		}

		// Find homography matrix
		homography := gocv.FindHomography(refPoints, imgPoints, gocv.Ransac, 3.0)

		// Warp the current image to align with the reference image
		alignedMat := gocv.NewMat()
		gocv.WarpPerspective(mats[i], &alignedMat, homography, refMat.Size())

		// Convert aligned Mat back to image.Image
		alignedImg, err := alignedMat.ToImage()
		if err != nil {
			return nil, err
		}
		alignedImages[i] = alignedImg
	}

	return alignedImages, nil
}
