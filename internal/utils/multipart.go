package utils

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func CreateMultipartForm(fields map[string]string, files map[string]string) (*bytes.Buffer, string, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for name, val := range fields {
		x, err := w.CreateFormField(name)
		if err != nil {
			return nil, "", err
		}
		x.Write([]byte(val))
	}

	for name, filepath := range files {
		if filepath == "" {
			continue
		}

		file, err := os.Open(filepath)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()

		part, err := w.CreateFormFile(name, path.Base(filepath))
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", err
		}
	}

	w.Close()
	return buf, w.FormDataContentType(), nil
}

// form validation tests just want to convert an *os.File to a *multipart.FileHeader
func CreateMultipartFileHeader(filePath string) (*multipart.FileHeader, error) {
	buf, contentType, err := CreateMultipartForm(map[string]string{}, map[string]string{"file": filePath})
	if err != nil {
		return nil, err
	}

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	buffReader := bytes.NewReader(buf.Bytes())
	formReader := multipart.NewReader(buffReader, params["boundary"])

	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, err
	}

	return files[0], nil
}

func GetMultipartFileMime(file multipart.File) (string, error) {
	buf := make([]byte, 512)

	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buf)
	return mimeType, nil
}
