package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/google/uuid"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

const (
	LargeThumbnailWidth     = 640
	LargeThumbnailHeight    = 360
	StandardThumbnailWidth  = 480
	StandardThumbnailHeight = 270
	MinProfileWidth         = 160
	TargetProfileWidth      = 176
)

// Functions used within handlers that save images
func (app *application) deleteImage(prefix, imgName string) error {
	imgSubPath := path.Join(prefix, imgName)
	app.infoLog.Printf("Deleting %s\n", imgSubPath)
	return app.fileStorage.DeleteFile(imgSubPath)
}

func (app *application) saveCastImages(member *models.CastMember) error {
	if member.ThumbnailFile != nil {
		thumbFile, err := member.ThumbnailFile.Open()
		if err != nil {
			return err
		}
		defer thumbFile.Close()

		img, err := processThumbnailImage(thumbFile, StandardThumbnailWidth, StandardThumbnailHeight)
		if err != nil {
			app.errorLog.Print("error processing thumbnail image")
			return err
		}

		var dstFile bytes.Buffer
		jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
		err = app.fileStorage.SaveFile(path.Join("cast", "thumbnail", *member.ThumbnailName), &dstFile)
		if err != nil {
			app.errorLog.Print("error saving thumbnail img")
			return err
		}
	}

	if member.ProfileFile == nil {
		return nil
	}

	profileFile, err := member.ProfileFile.Open()
	if err != nil {
		return err
	}
	defer profileFile.Close()

	img, err := processProfileImage(profileFile)
	if err != nil {
		app.errorLog.Print("error processing profile image")
		return err
	}

	var profileDstFile bytes.Buffer
	jpeg.Encode(&profileDstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join("cast", "profile", *member.ProfileImg), &profileDstFile)
	if err != nil {
		app.errorLog.Print("error saving profile img")
		return err
	}
	return nil
}

func (app *application) saveLargeThumbnail(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := processThumbnailImage(file, LargeThumbnailWidth, LargeThumbnailHeight)
	if err != nil {
		return err
	}

	var dstFile bytes.Buffer
	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join(prefix, "large", imgName), &dstFile)
	if err != nil {
		return err
	}

	file.Seek(0, 0)

	img, err = processThumbnailImage(file, StandardThumbnailWidth, StandardThumbnailHeight)
	if err != nil {
		return err
	}

	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join(prefix, imgName), &dstFile)
	if err != nil {
		return err
	}

	return nil

}

func (app *application) saveThumbnail(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := processThumbnailImage(file, StandardThumbnailWidth, StandardThumbnailHeight)
	if err != nil {
		return err
	}

	var dstFile bytes.Buffer
	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join(prefix, "large", imgName), &dstFile)
	return err
}

func (app *application) saveProfileImage(imgName string, prefix string, fileHeader *multipart.FileHeader) error {
	if imgName == "" {
		return fmt.Errorf("Image name cannot be blank")
	}

	profileFile, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer profileFile.Close()

	img, err := processProfileImage(profileFile)
	if err != nil {
		return err
	}

	var dstFile bytes.Buffer
	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join(prefix, imgName), &dstFile)
	if err != nil {
		return err
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

func processThumbnailImage(file io.Reader, width, height int) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	img = utils.CenterCropToAspectRatio(img, 16.0/9.0)
	img = utils.ResizeImage(img, width, height)
	return img, nil
}

func processProfileImage(file io.Reader) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	img = utils.CenterCropToAspectRatio(img, 1.0)
	img = utils.ResizeImage(img, TargetProfileWidth, TargetProfileWidth)
	return img, nil
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
