package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

func CropImg(file io.Reader, ext string, cropRect image.Rectangle) (io.Reader, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	subImager, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if !ok {
		return nil, fmt.Errorf("image type does not support cropping")
	}

	croppedImg := subImager.SubImage(cropRect)

	var buf bytes.Buffer

	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(&buf, croppedImg, nil)
	case ".png":
		err = png.Encode(&buf, croppedImg)
	default:
		return nil, fmt.Errorf("unsupported image output format")
	}

	return &buf, err
}

func CenterCropToAspectRatio(img image.Image, targetRatio float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	targetWidth := width
	targetHeight := int(float64(width) / targetRatio)

	if targetHeight > height {
		targetHeight = height
		targetWidth = int(float64(height) * targetRatio)
	}

	x0 := (width - targetWidth) / 2
	y0 := (height - targetHeight) / 2
	x1 := x0 + targetWidth
	y1 := y0 + targetHeight

	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x0, y0, x1, y1))
}

func ResizeImage(img image.Image, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func GetImageDimensions(file io.ReadSeekCloser) (int, int, error) {
	imgConf, _, err := image.DecodeConfig(file)
	defer file.Seek(0, 0)

	if err != nil {
		return 0, 0, err
	}
	return imgConf.Width, imgConf.Height, err
}
