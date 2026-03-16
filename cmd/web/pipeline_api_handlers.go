package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const MAX_FILE_SIZE = 262_144_000 // 250 MiB

var allowedMIMETypes = map[string]bool{
	"video/mp4":        true,
	"video/quicktime":  true,
	"video/x-msvideo":  true, // avi
	"video/webm":       true,
	"video/x-matroska": true,
}

func (app *application) generateSketchVideoS3PutUrl(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("sketch id param not defined"))
		return
	}

	videos, err := app.services.Sketches.GetVideos(sketchId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// currently not accepting more than one video per sketch
	if len(videos) > 0 {
		app.badRequestResponse(w, r, fmt.Errorf("sketch already has a video"))
		return
	}

	var input struct {
		FileName    string `json:"fileName"`
		ContentType string `json:"contentType"`
		FileSize    int    `json:"fileSize"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FileName == "" {
		app.failedValidationResponse(w, r, map[string]string{"error": "file name must be specified"})
		return
	}

	if !allowedMIMETypes[input.ContentType] {
		app.failedValidationResponse(w, r, map[string]string{"error": "unsupported file type"})
		return
	}

	if input.FileSize > MAX_FILE_SIZE {
		app.failedValidationResponse(w, r, map[string]string{"error": "max file upload size 250 MiB"})
		return
	}

	ext := filepath.Ext(input.FileName)
	s3Key := fmt.Sprintf("video/%s%s", uuid.New().String(), ext)
	url, err := app.fileStorage.PresignedUploadURL(s3Key, 30*time.Minute, input.FileSize)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"uploadUrl": url,
		"s3Key":     s3Key,
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}

func (app *application) sketchVideoUploaded(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("sketch id param not defined"))
		return
	}

	exists, err := app.sketches.Exists(sketchId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		app.failedValidationResponse(w, r, map[string]string{"error": "sketch does not exist"})
		return
	}

	var input struct {
		S3Key string `json:"s3Key"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.S3Key == "" {
		app.failedValidationResponse(w, r, map[string]string{"error": "s3 key must be defined"})
		return
	}

	exists, err = app.fileStorage.Exists(input.S3Key)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		app.failedValidationResponse(w, r, map[string]string{"error": "s3 key does not exist in bucket"})
		return
	}

	video, err := app.services.Sketches.AddVideo(sketchId, input.S3Key)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_, err = app.services.Pipeline.AddPipelineJob(safeDeref(video.ID))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	videos, err := app.services.Sketches.GetVideos(sketchId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"videos": videos,
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}

func (app *application) getSketchVideos(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("sketch id param not defined"))
		return
	}

	videos, err := app.services.Sketches.GetVideos(sketchId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"videos": videos,
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}
