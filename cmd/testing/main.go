package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"
	"sketchdb.cozycole.net/internal/utils"
)

const (
	MinThumbnailWidth     = 480
	MinThumbnailHeight    = 360
	TargetThumbnailWidth  = 480
	TargetThumbnailHeight = 360
	FinalThumbnailHeight  = 270
	MinProfileWidth       = 256
	TargetProfileWidth    = 256
)

// Validate minimum image dimensions
func validateImageSize(img image.Image, minWidth, minHeight int) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < minWidth || height < minHeight {
		return errors.New("image too small for resizing")
	}
	return nil
}

// Center crop to 4:3 aspect ratio (before resizing)
func centerCropToAspectRatio(img image.Image, targetRatio float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Determine target dimensions while keeping the correct ratio
	targetWidth := width
	targetHeight := int(float64(width) / targetRatio)

	if targetHeight > height {
		targetHeight = height
		targetWidth = int(float64(height) * targetRatio)
	}

	// Calculate crop box (centered)
	x0 := (width - targetWidth) / 2
	y0 := (height - targetHeight) / 2
	x1 := x0 + targetWidth
	y1 := y0 + targetHeight

	// Crop and return
	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x0, y0, x1, y1))
}

// Resize image using Go standard library
func resizeImage(img image.Image, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

// Crop top and bottom to final size (480x270)
func cropFinalSize(img image.Image, finalHeight int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	y0 := (height - finalHeight) / 2
	y1 := y0 + finalHeight

	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, y0, width, y1))
}

// Process image: Validate, crop, resize, and final crop
func processImage(inputFile, outputFile string) error {
	// Open input image
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Validate image size
	if err := validateImageSize(img, MinThumbnailWidth, MinThumbnailHeight); err != nil {
		return err
	}

	// Step 1: Center crop to 4:3 aspect ratio
	img = centerCropToAspectRatio(img, 4.0/3.0)

	// Step 2: Resize to 480x360
	img = resizeImage(img, TargetThumbnailWidth, TargetThumbnailHeight)

	// Step 3: Crop top and bottom to 480x270
	img = cropFinalSize(img, FinalThumbnailHeight)

	// Save output
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode image in the correct format
	if format == "png" {
		return png.Encode(outFile, img)
	}
	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: 85})
}

func main() {
	inputFile := "input.jpg"
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("error opening file")
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("error decoding file")
		return
	}
	var dstFile bytes.Buffer

	img = utils.CenterCropToAspectRatio(img, 1.0)
	img = utils.ResizeImage(img, TargetProfileWidth, TargetProfileWidth)
	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})

	err = os.WriteFile("profile_output.jpg", dstFile.Bytes(), 0644)
	if err != nil {
		fmt.Printf("error creating file")
		return
	}
	println("Thumbnail successfully processed and saved.")
}
