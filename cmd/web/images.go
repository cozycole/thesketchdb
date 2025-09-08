package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/google/uuid"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

const (
	ThumbnailAspectRatio  = 16.0 / 9.0
	JPGQuality            = 82
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

func (app *application) deleteImage(prefix, imgName string) error {
	for _, size := range []string{"small", "medium", "large"} {
		imgSubPath := path.Join(prefix, size, imgName)
		app.infoLog.Printf("Deleting %s\n", imgSubPath)
		err := app.fileStorage.DeleteFile(imgSubPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) saveCastImages(member *models.CastMember) error {
	if member.ThumbnailFile != nil {
		err := app.saveMediumThumbnail(
			safeDeref(member.ThumbnailName),
			"/cast/thumbnail", member.ThumbnailFile,
		)

		if err != nil {
			return err
		}
	}

	if member.ProfileFile == nil {
		return nil
	}

	err := app.saveMediumProfile(
		safeDeref(member.ProfileImg),
		"/cast/profile", member.ProfileFile,
	)

	if err != nil {
		return err
	}

	return nil
}

// saveLargeThumbnail saves large, medium and small resolutions
func (app *application) saveLargeThumbnail(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	width, height, err := utils.GetImageDimensions(file)
	if err != nil {
		return err
	}
	width, height = GetLargest16x9Dimensions(width, height)

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	images := map[string]image.Image{}
	images["small"], err = processThumbnailImage(img, SmallThumbnailWidth, SmallThumbnailHeight)
	if err != nil {
		return err
	}

	if width > MediumThumbnailWidth {
		images["medium"], err = processThumbnailImage(img, MediumThumbnailWidth, MediumThumbnailHeight)
	} else if width > SmallThumbnailWidth {
		images["medium"], err = processThumbnailImage(img, width, height)
	} else {
		images["medium"] = images["small"]
	}

	if width > LargeThumbnailWidth {
		images["large"], err = processThumbnailImage(img, LargeThumbnailWidth, LargeThumbnailHeight)
	} else if width > MediumThumbnailWidth {
		images["large"], err = processThumbnailImage(img, width, height)
	} else {
		images["large"] = images["medium"]
	}

	if err != nil {
		return err
	}

	// 5) Save all files to /small /medium /large respectively
	for size, img := range images {
		var dstFile bytes.Buffer
		err := jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: JPGQuality})
		if err != nil {
			return err
		}

		err = app.fileStorage.SaveFile(path.Join(prefix, size, imgName), &dstFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// saveMediumThumbnail saves medium and small resolutions
func (app *application) saveMediumThumbnail(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	width, height, err := utils.GetImageDimensions(file)
	if err != nil {
		return err
	}

	width, height = GetLargest16x9Dimensions(width, height)

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	images := map[string]image.Image{}
	images["small"], err = processThumbnailImage(img, SmallThumbnailWidth, SmallThumbnailHeight)
	if err != nil {
		return err
	}

	if width > MediumThumbnailWidth {
		images["medium"], err = processThumbnailImage(img, MediumThumbnailWidth, MediumThumbnailHeight)
	} else if width > SmallThumbnailWidth {
		images["medium"], err = processThumbnailImage(img, width, height)
	} else {
		images["medium"] = images["small"]
	}

	if err != nil {
		return err
	}

	// 5) Save all files to /small and /medium respectively
	for size, img := range images {
		var dstFile bytes.Buffer
		err := jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: JPGQuality})
		if err != nil {
			return err
		}

		err = app.fileStorage.SaveFile(path.Join(prefix, size, imgName), &dstFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *application) saveLargeProfile(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	width, height, err := utils.GetImageDimensions(file)
	if err != nil {
		return err
	}

	width = min(width, height)

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	images := map[string]image.Image{}
	images["small"], err = processProfileImage(img, SmallProfileWidth)
	if err != nil {
		return err
	}

	if width > MediumProfileWidth {
		images["medium"], err = processProfileImage(img, MediumProfileWidth)
	} else if width > SmallProfileWidth {
		images["medium"], err = processProfileImage(img, width)
	} else {
		images["medium"] = images["small"]
	}

	if err != nil {
		return err
	}

	if width > LargeProfileWidth {
		images["large"], err = processProfileImage(img, LargeProfileWidth)
	} else if width > MediumProfileWidth {
		images["large"], err = processProfileImage(img, width)
	} else {
		images["large"] = images["medium"]
	}

	if err != nil {
		return err
	}

	// 5) Save all files to /small /medium /large respectively
	for size, img := range images {
		var dstFile bytes.Buffer
		err := jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: JPGQuality})
		if err != nil {
			return err
		}

		err = app.fileStorage.SaveFile(path.Join(prefix, size, imgName), &dstFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *application) saveMediumProfile(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	width, height, err := utils.GetImageDimensions(file)
	if err != nil {
		return err
	}

	width = min(width, height)

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	images := map[string]image.Image{}
	images["small"], err = processProfileImage(img, SmallProfileWidth)
	if err != nil {
		return err
	}

	if width > MediumProfileWidth {
		images["medium"], err = processProfileImage(img, MediumProfileWidth)
	} else if width > SmallProfileWidth {
		images["medium"], err = processProfileImage(img, width)
	} else {
		images["medium"] = images["small"]
	}

	if err != nil {
		return err
	}

	// 5) Save all files to /small /medium /large respectively
	for size, img := range images {
		var dstFile bytes.Buffer
		err := jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: JPGQuality})
		if err != nil {
			return err
		}

		err = app.fileStorage.SaveFile(path.Join(prefix, size, imgName), &dstFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateThumbnailName(fileHeader *multipart.FileHeader) (string, error) {
	thumbnailId := uuid.New().String()
	thumbnailExtension, err := getFileExtension(fileHeader)
	if err != nil {
		return "", err
	}
	return thumbnailId + thumbnailExtension, nil
}

func processThumbnailImage(img image.Image, width, height int) (image.Image, error) {

	img = utils.CenterCropToAspectRatio(img, ThumbnailAspectRatio)
	img = utils.ResizeImage(img, width, height)
	return img, nil
}

func processProfileImage(img image.Image, width int) (image.Image, error) {
	img = utils.CenterCropToAspectRatio(img, 1.0)
	img = utils.ResizeImage(img, width, width)
	return img, nil
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

func getFileExtension(header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file")
	}

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", err
	}
	defer file.Seek(0, 0)

	mime, ok := mimeToExt[http.DetectContentType(buf)]
	if !ok {
		return "", fmt.Errorf("Mime does not exist in extension table, bad file")
	}
	return mime, nil
}
