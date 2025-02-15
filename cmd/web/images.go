package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

// Functions used within handlers that save images
func (app *application) saveVideoThumbnail(video *models.Video) error {
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
	if width == 480 && height == 360 {
		subStrs := strings.Split(video.ThumbnailName, ".")
		ext := "." + subStrs[len(subStrs)-1]
		rect := image.Rect(0, 45, 480, 315)
		dstFile, err = utils.CropImg(file, ext, rect)
		if err != nil {
			return err
		}
	}

	err = app.fileStorage.SaveFile(path.Join("video", video.ThumbnailName), dstFile)
	if err != nil {
		return err
	}

	return nil
}

// Generates names for video thumbnail and character thumbnails
// func addVideoImageNames(video *models.Video) error {
// 	file, err := video.ThumbnailFile.Open()
// 	if err != nil {
// 		return fmt.Errorf("unable to open file")
// 	}
// 	mimeType, err := utils.GetMultipartFileMime(file)
// 	if err != nil {
// 		return fmt.Errorf("unable to identify mime")
// 	}
// 	file.Close()
// 	if video.Title == "" {
// 		return fmt.Errorf("unable to name video thumbnail: title not defined")
// 	}
// 	slug := models.CreateSlugName(video.Title, maxFileNameLength)
// 	video.ThumbnailName = fmt.Sprintf("%s-%d%s", slug, 1, mimeToExt[mimeType])

// for _, c := range video.Cast {
// 	file, err = c.ThumbnailFile.Open()
// 	if err != nil {
// 		return fmt.Errorf("unable to identify mime")
// 	}
//
// 	mimeType, err = utils.GetMultipartFileMime(file)
// 	file.Close()
// 	if err != nil {
// 		return fmt.Errorf("unable to identify mime")
// 	}
//
// 	var thumbName string
// 	if c.Character.ID != nil {
// 		thumbName = fmt.Sprintf("%d-%d-%d", video.ID, *c.Actor.ID, *c.Character.ID)
// 	} else {
// 		thumbName = fmt.Sprintf("%d-%d", video.ID, *c.Actor.ID)
// 	}
//
// 	thumbName = toURLSafeBase64MD5(thumbName) + mimeToExt[mimeType]
// 	c.ThumbnailName = &thumbName
// }
// 	return nil
// }

func generateThumbnailHash(videoID int) string {
	data := fmt.Sprintf("%d-%d", videoID, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded[:22], "=")
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

