package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
)

var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func GenerateFileName(fileHeader *multipart.FileHeader) (string, error) {
	thumbnailId := uuid.New().String()
	thumbnailExtension, err := GetFileExtension(fileHeader)
	if err != nil {
		return "", err
	}
	return thumbnailId + thumbnailExtension, nil
}

func GetFileExtension(header *multipart.FileHeader) (string, error) {
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
