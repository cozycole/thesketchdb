package media

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"math"
	"path"

	"sketchdb.cozycole.net/internal/fileStore"
)

const (
	ThumbnailAspectRatio  = 16.0 / 9.0
	JPGQuality            = 85
	LargeThumbnailWidth   = 1280
	LargeThumbnailHeight  = 720
	MediumThumbnailWidth  = 640
	MediumThumbnailHeight = 360
	SmallThumbnailWidth   = 320
	SmallThumbnailHeight  = 180
	LargeProfileWidth     = 512
	MediumProfileWidth    = 256
	SmallProfileWidth     = 88
)

// Given a src image, saves cropped/covered images based on small / medium / large dimensions
// to imgStore
//
// e.g. if maxSize == Medium, then only small / medium images will be saved
//
// imgName is the baseName of the file path and prefix is the path to it WITHOUT the size
// So if prefix == "/cast/profile" and imgName == "abcdefg.jpg" then it will be saved as
// /cast/profile/{size}/abcdefg.jpg
func RunImagePipeline(
	src []byte,
	maxSize Size,
	imgType ImageType,
	imgName, prefix string,
	imgStore fileStore.FileStorageInterface,
) error {
	variants, err := CreateImageVariants(src, maxSize, imgType)
	if err != nil {
		return err
	}

	err = SaveImageVariants(imgStore, prefix, imgName, variants)
	if err != nil {
		return err
	}

	return nil
}

func CreateImageVariants(src []byte, maxSize Size, imgType ImageType) ([]Variant, error) {
	img, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	img = RemoveBorders(img, 15, 2)

	specs, err := createImageVariantSpec(img, maxSize, imgType)
	if err != nil {
		return nil, err
	}

	variants, err := Process(img, specs)
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func createImageVariantSpec(img image.Image, maxSize Size, imgType ImageType) ([]VariantSpec, error) {
	var smallOutputWidth, mediumOutputWidth, largeOutputWidth int
	var smallOutputHeight, mediumOutputHeight, largeOutputHeight int
	if imgType == Profile {
		smallOutputWidth = SmallProfileWidth
		smallOutputHeight = SmallProfileWidth
		mediumOutputWidth = MediumProfileWidth
		mediumOutputHeight = MediumProfileWidth
		largeOutputWidth = LargeProfileWidth
		largeOutputHeight = LargeProfileWidth
	} else {
		smallOutputWidth = SmallThumbnailWidth
		smallOutputHeight = SmallThumbnailHeight
		mediumOutputWidth = MediumThumbnailWidth
		mediumOutputHeight = MediumThumbnailHeight
		largeOutputWidth = LargeThumbnailWidth
		largeOutputHeight = LargeThumbnailHeight
	}

	var w, h int
	if imgType == Profile {
		w = min(img.Bounds().Dx(), img.Bounds().Dy())
		h = w
	} else {
		w, h = GetLargest16x9Dimensions(img.Bounds().Dx(), img.Bounds().Dy())
	}

	specs := []VariantSpec{}
	// we always at the very least output a small version
	specs = append(specs, VariantSpec{
		Name:    "small",
		Width:   smallOutputWidth,
		Height:  smallOutputHeight,
		Mode:    FitCover,
		Format:  FormatJPEG,
		Quality: JPGQuality,
	})

	mWidth := mediumOutputWidth
	mHeight := mediumOutputHeight
	if w > smallOutputWidth && w < mediumOutputWidth {
		mWidth = w
		mHeight = h
	} else if w < smallOutputWidth {
		mWidth = smallOutputWidth
		mHeight = smallOutputHeight
	}

	if maxSize == Large || maxSize == Medium {
		specs = append(specs, VariantSpec{
			Name:    "medium",
			Width:   mWidth,
			Height:  mHeight,
			Mode:    FitCover,
			Format:  FormatJPEG,
			Quality: JPGQuality,
		})
	}

	lWidth := largeOutputWidth
	lHeight := largeOutputHeight
	if w > mediumOutputWidth && w < largeOutputWidth {
		lWidth = w
		lHeight = h
	} else if w < mediumOutputWidth {
		lWidth = mediumOutputWidth
		lHeight = mediumOutputHeight
	}

	if maxSize == Large {
		specs = append(specs, VariantSpec{
			Name:    "large",
			Width:   lWidth,
			Height:  lHeight,
			Mode:    FitCover,
			Format:  FormatJPEG,
			Quality: JPGQuality,
		})
	}

	return specs, nil

}

func GetLargest16x9Dimensions(width, height int) (int, int) {
	const aspectRatio = 16.0 / 9.0
	widthBasedHeight := int(math.Round(float64(width) / aspectRatio))
	heightBasedWidth := int(math.Round(float64(height) * aspectRatio))

	if widthBasedHeight <= height {
		return width, widthBasedHeight
	} else {
		return heightBasedWidth, height
	}
}

func SaveImageVariants(imgStore fileStore.FileStorageInterface, prefix, fileName string, variants []Variant) error {
	for _, v := range variants {
		fileName := fmt.Sprintf("%s/%s/%s", prefix, v.Name, fileName)
		err := imgStore.SaveFile(fileName, bytes.NewBuffer(v.Bytes))
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteImageVariants(imgStore fileStore.FileStorageInterface, prefix, fileName string) error {
	for _, size := range []string{"small", "medium", "large"} {
		imgSubPath := path.Join(prefix, size, fileName)
		err := imgStore.DeleteFile(imgSubPath)
		if err != nil {
			return err
		}
	}
	return nil
}
