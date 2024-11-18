package main

import (
	"fmt"
	"mime/multipart"
	"strings"

	"sketchdb.cozycole.net/internal/validator"
)

// Functions for form validation within handlers
type FormInterface interface {
	validator.Validator
}

// Changes to the form fields must be updated in their respective
// validate functions
type addCreatorForm struct {
	Name                string                `form:"name"`
	URL                 string                `form:"url"`
	EstablishedDate     string                `form:"establishedDate"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	validator.Validator `form:"-"`
}

func (app *application) validateAddCreatorForm(form *addCreatorForm) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.URL), "url", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.EstablishedDate), "establishedDate", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.EstablishedDate), "establishedDate", "Date not of correct format YYYY-MM-DD")
	form.CheckField(form.ProfileImage != nil, "profileImg", "Please upload an image")

	if form.ProfileImage == nil {
		return
	}

	profileImg, err := form.ProfileImage.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a jpg or png")
		return
	}
	defer profileImg.Close()

	form.CheckField(validator.IsMime(profileImg, "image/jpeg", "image/png"), "profileImg", "Uploaded file must be jpg or png")
}

type addPersonForm struct {
	First               string                `form:"first"`
	Last                string                `form:"last"`
	BirthDate           string                `form:"birthDate"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	validator.Validator `form:"-"`
}

func (app *application) validateAddPersonForm(form *addPersonForm) {
	form.CheckField(validator.NotBlank(form.First), "first", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Last), "last", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.BirthDate), "birthDate", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.BirthDate), "birthDate", "Date not of correct format YYYY-MM-DD")
	form.CheckField(form.ProfileImage != nil, "profileImg", "Please upload an image")

	if form.ProfileImage == nil {
		return
	}

	profileImg, err := form.ProfileImage.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a jpg or png")
		return
	}
	defer profileImg.Close()

	form.CheckField(validator.IsMime(profileImg, "image/jpeg", "image/png"), "profileImg", "Uploaded file must be jpg or png")
}

// PersonIDs have name form field names of form peopleId[i]
// if there are spaces between indexes say peopleId[0] : 1, peopleId[2] : 3
// the result is zero filled so []int{1,0,3}
type addVideoForm struct {
	Title               string                `form:"title"`
	URL                 string                `form:"url"`
	Rating              string                `form:"rating"`
	UploadDate          string                `form:"uploadDate"`
	Thumbnail           *multipart.FileHeader `img:"thumbnail"`
	CreatorID           int                   `form:"creator"`
	PersonIDs           []int                 `form:"peopleId"`
	PersonInputs        []string              `form:"peopleText"`
	validator.Validator `form:"-"`
}

// We need this function to have access to the apps state
// to validate based on database queries
func (app *application) validateAddVideoForm(form *addVideoForm) {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.URL), "url", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Rating), "rating", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.UploadDate), "uploadDate", "This field cannot be blank")
	form.CheckField(form.CreatorID != 0, "creator", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.UploadDate), "uploadDate", "Date must be of the format YYYY-MM-DD")
	form.CheckField(
		validator.PermittedValue(strings.ToLower(form.Rating), "pg", "pg-13", "r"),
		"rating",
		"Rating must be PG, PG-13 or R (case insensitive)",
	)

	if form.CreatorID != 0 {
		form.CheckField(
			validator.BoolWithError(app.creators.Exists(form.CreatorID)),
			"creator",
			"Unable to find creator, please add them",
		)
	}

	for i, a := range form.PersonIDs {
		htmlFieldName := fmt.Sprintf("peopleId[%d]", i)
		if !validator.IsZero(a) {
			form.CheckField(
				validator.BoolWithError(app.people.Exists(a)),
				htmlFieldName,
				"Person does not exist. Please add it, then resubmit video!",
			)
		}
	}

	form.CheckField(form.Thumbnail != nil, "thumbnail", "Please upload an image")
	if form.Thumbnail == nil {
		return
	}

	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a jpg or png")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg or png")
}

// func validateProfileImage[F FormInterface](form F, )
