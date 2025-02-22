package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
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

// Functions used within handlers that save images
func (app *application) saveVideoThumbnail(video *models.Video) error {
	file, err := video.ThumbnailFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := processThumbnailImage(file)
	if err != nil {
		return err
	}

	var dstFile bytes.Buffer
	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join("video", video.ThumbnailName), &dstFile)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) deleteVideoThumbnail(video *models.Video) error {
	// within wherever images are stored, you get video images by
	// appending /video/ to the thumbnailName e.g. /video/asdasdfjkl.jpg
	thumbnailSubPath := path.Join("video", video.ThumbnailName)
	return app.fileStorage.DeleteFile(thumbnailSubPath)
}

func (app *application) saveCastImages(member *models.CastMember) error {
	thumbFile, err := member.ThumbnailFile.Open()
	if err != nil {
		return err
	}
	defer thumbFile.Close()

	img, err := processThumbnailImage(thumbFile)
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

	profileFile, err := member.ProfileFile.Open()
	if err != nil {
		return err
	}
	defer profileFile.Close()

	img, err = processProfileImage(profileFile)
	if err != nil {
		app.errorLog.Print("error processing profile image")
		return err
	}

	jpeg.Encode(&dstFile, img, &jpeg.Options{Quality: 85})
	err = app.fileStorage.SaveFile(path.Join("cast", "profile", *member.ThumbnailName), &dstFile)
	if err != nil {
		app.errorLog.Print("error saving profile img")
		return err
	}
	return nil
}

func generateThumbnailName(id int, fileHeader *multipart.FileHeader) (string, error) {
	thumbnailHash := generateThumbnailHash(id)
	thumbnailExtension, err := getFileExtension(fileHeader)
	if err != nil {
		return "", err
	}
	return thumbnailHash + thumbnailExtension, nil
}

func processThumbnailImage(file io.Reader) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	img = utils.CenterCropToAspectRatio(img, 4.0/3.0)
	img = utils.ResizeImage(img, TargetThumbnailWidth, TargetThumbnailHeight)
	img = utils.CropFinalSize(img, FinalThumbnailHeight)
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

func generateThumbnailHash(id int) string {
	data := fmt.Sprintf("%d-%d", id, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded[:22], "=")
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

	mime, ok := mimeToExt[http.DetectContentType(buf)]
	if !ok {
		return "", fmt.Errorf("Mime does not exists in extension table, bad file")
	}
	return mime, nil
}
