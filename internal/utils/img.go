package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
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

func GetImageDimensions(file io.ReadSeekCloser) (int, int, error) {
	imgConf, _, err := image.DecodeConfig(file)
	defer file.Seek(0, 0)

	if err != nil {
		return 0, 0, err
	}
	return imgConf.Width, imgConf.Height, err
}
