package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"path"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

// Functions used within handlers that save images

func (app *application) saveVideoImages(video *models.Video) error {
	// save video thumbnail	
	vidDirSubPath := path.Join("video", strconv.Itoa(video.ID))
	file, err := video.ThumbnailFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	width, height, err := utils.GetImageDimensions(file)
	if err != nil {
		return err
	}

	var dstFile io.Reader
	dstFile = file
	// This is stock youtube thumbnail dimensions but has black
	// top/bottom borders that need to be removed
	subStrs := strings.Split(video.ThumbnailName, ".")
	ext := "." +subStrs[len(subStrs)-1]
	if width == 480 && height == 360 {
		rect := image.Rect(0, 45, 480, 315)
		dstFile, err = utils.CropImg(file, ext, rect)
		if err != nil {
			return err
		}
	}

	err = app.fileStorage.SaveFile(path.Join(vidDirSubPath, video.ThumbnailName), dstFile)
	if err != nil {
		return err
	}

	for _, c := range video.Cast {
		err = app.saveCharacterThumbnail(c.ThumbnailFile, path.Join(vidDirSubPath, *c.ThumbnailName))
		if err != nil {
			return err
		}
	}
	return nil
}

// Generates names for video thumbnail and character thumbnails
func addVideoImageNames(video *models.Video) error {
	file, err := video.ThumbnailFile.Open()
	if err != nil {
		return fmt.Errorf("unable to open file")
	}
	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		return fmt.Errorf("unable to identify mime")
	}
	file.Close()
	if video.Title == "" {
		return fmt.Errorf("unable to name video thumbnail: title not defined")
	}
	slug := models.CreateSlugName(video.Title, maxFileNameLength)
	video.ThumbnailName = fmt.Sprintf("%s-%d%s", slug, 1, mimeToExt[mimeType])

	for _, c := range video.Cast {
		file, err = c.ThumbnailFile.Open()
		if err != nil {
			return fmt.Errorf("unable to identify mime")
		}

		mimeType, err = utils.GetMultipartFileMime(file)
		file.Close()
		if err != nil {
			return fmt.Errorf("unable to identify mime")
		}

		var thumbName string
		if c.Character.ID != nil {
			thumbName = fmt.Sprintf("%d-%d-%d", video.ID, *c.Actor.ID, *c.Character.ID)
		} else {
			thumbName = fmt.Sprintf("%d-%d", video.ID, *c.Actor.ID)
		}

		thumbName = toURLSafeBase64MD5(thumbName) + mimeToExt[mimeType]
		c.ThumbnailName = &thumbName
	}
	return nil
}

func toURLSafeBase64MD5(input string) string {
	hash := md5.Sum([]byte(input))
	urlSafe := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])
	return urlSafe
}

func (app *application) saveCharacterThumbnail(header *multipart.FileHeader, path string) error {
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("unable to open file: %s", path)
	}
	defer file.Close()
	
	app.fileStorage.SaveFile(path, file)
	return nil
}