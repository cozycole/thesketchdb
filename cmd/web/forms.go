package main

import (
	"mime/multipart"

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

func validateAddCreatorForm(form *addCreatorForm) {
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

type addActorForm struct {
	First               string                `form:"first"`
	Last                string                `form:"last"`
	BirthDate           string                `form:"birthDate"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	validator.Validator `form:"-"`
}

func validateAddActorForm(form *addActorForm) {
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

// func validateProfileImage[F FormInterface](form F, )
