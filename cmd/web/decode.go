package main

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
)

// see ~/go/pkg/mod/github.com/go-playground/form/v4@v4.2.1/README.md
// for doc on the form decoder
func (app *application) decodePostForm(r *http.Request, dst any) error {
	var err error
	contentType := r.Header.Get("Content-Type")
	isMultipart := strings.Contains(contentType, "multipart/form-data")
	if isMultipart {
		err = r.ParseMultipartForm(10 << 20)
	} else {
		err = r.ParseForm()
	}

	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	if isMultipart {

		// checks if struct tag with key "img" exists
		// If it does, we take the value of the struct tag and find any
		// files headers within r.MultipartForm.File with key that either equals it
		// (if the field with the struct is of type *multipart.FileHeader)
		// or prefixed with it (if field is of type []*multipart.FileHeader and the form files
		// are keyed with tag_value[x], see getIndexedValuesWithGaps for example)
		// The file headers are then set as the value in the form struct
		v := reflect.ValueOf(dst).Elem()
		structType := v.Type()
		fmt.Println(structType)

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := structType.Field(i)

			if tagValue := fieldType.Tag.Get("img"); tagValue != "" {
				if field.Kind() == reflect.Slice {
					fileHeaders, err := getIndexedValuesWithGaps(r.MultipartForm.File, tagValue)
					if err != nil {
						return err
					}

					field.Set(reflect.ValueOf(fileHeaders))
				} else {
					fileHeaders, ok := r.MultipartForm.File[tagValue]
					if ok && len(fileHeaders) > 0 {
						newVal := reflect.ValueOf(fileHeaders[0])
						field.Set(newVal)
					}
				}
			}
		}
	}

	return nil
}

// This function returns an array of values from any keys of format prefix[x]
// The returned array is in sorted order and zero filled if x index
// have spaces between.
// Example (Strings used here but values will be of type *multipart.FileHeader):
//
//	inputMap := map[string]string{
//		"example[0]": "zero",
//		"example[1]": "one",
//		"example[3]": "three",
//		"example[5]": "five",
//		"other[0]":   "not included",
//	}
//
// getIndexedValuesWithGaps(inputMap, "example")
// -> ["zero", "one", "", "three", "", "five"]
func getIndexedValuesWithGaps(inputMap map[string][]*multipart.FileHeader, prefix string) ([]*multipart.FileHeader, error) {

	pattern := fmt.Sprintf(`^%s\[(\d+)\]$`, regexp.QuoteMeta(prefix))
	re := regexp.MustCompile(pattern)

	indexedValues := make(map[int][]*multipart.FileHeader)
	maxIndex := -1 // since it's 0 indexed
	for key, value := range inputMap {
		matches := re.FindStringSubmatch(key)
		if matches != nil {
			index, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, err
			}
			indexedValues[index] = value
			if index > maxIndex {
				maxIndex = index
			}
		}
	}

	// no keys found with substring
	if maxIndex == -1 {
		return nil, nil
	}

	result := make([]*multipart.FileHeader, maxIndex+1)
	for index := 0; index <= maxIndex; index++ {
		if fileHeaders, exists := indexedValues[index]; exists {
			// just use the first file header with that name
			// (there shouldn't be multiple)
			result[index] = fileHeaders[0]
		}
	}

	return result, nil
}
