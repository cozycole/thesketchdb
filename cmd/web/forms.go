package main

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
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
	Title               string                  `form:"title"`
	URL                 string                  `form:"url"`
	Rating              string                  `form:"rating"`
	UploadDate          string                  `form:"uploadDate"`
	Thumbnail           *multipart.FileHeader   `img:"thumbnail"`
	CreatorID           int                     `form:"creator"`
	PersonIDs           []int                   `form:"peopleId"`
	PersonInputs        []string                `form:"peopleText"`
	CharacterIDs        []int                   `form:"characterId"`
	CharacterInputs     []string                `form:"characterText"`
	CharacterThumbnails []*multipart.FileHeader `img:"characterThumbnail"`
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

	form.CheckField(form.Thumbnail != nil, "thumbnail", "Please upload an image")
	if form.Thumbnail == nil {
		return
	}

	// The following fields are slices are related records
	// (PersondIDs[0], CharacterIDs[0], CharacterThumbnails[0])
	minLength := min(len(form.PersonIDs), len(form.CharacterIDs), len(form.CharacterThumbnails))
	for i := 0; i < minLength; i++ {

		// if they're all empty, ignore this input
		if form.IsEmptyActorInput(i) {
			continue
		}

		htmlPeopleIdField := fmt.Sprintf("peopleId[%d]", i)
		htmlCharIdField := fmt.Sprintf("characterId[%d]", i)
		htmlCharThumbField := fmt.Sprintf("characterThumbnail[%d]", i)

		pid := form.PersonIDs[i]
		if !validator.IsZero(pid) {
			form.CheckField(
				validator.BoolWithError(app.people.Exists(pid)),
				htmlPeopleIdField,
				"Person does not exist. Please add it, then resubmit video!",
			)
		} else {
			form.CheckField(true, htmlPeopleIdField, "This field cannot be blank")
		}

		cid := form.CharacterIDs[i]
		if !validator.IsZero(cid) {
			form.CheckField(
				validator.BoolWithError(app.characters.Exists(cid)),
				htmlCharIdField,
				"Character does not exist. Please add it, then resubmit video!",
			)
		} else {
			form.CheckField(true, htmlCharIdField, "This field cannot be blank")
		}

		thumb := form.CharacterThumbnails[i]
		if thumb == nil {
			form.CheckField(true, htmlCharThumbField, "Please upload character thumbnail")
		}
		thumbnail, err := thumb.Open()
		if err != nil {
			form.AddFieldError(htmlCharThumbField, "Unable to open file, ensure it is a jpg or png")
			return
		}
		defer thumbnail.Close()

		form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"),
			htmlCharThumbField, "Uploaded file must be jpg or png")
	}

	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a jpg or png")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg or png")
	// Might want to check dimension ratios to make sure they work?
}

func (f *addVideoForm) IsEmptyActorInput(index int) bool {
	switch {
	case f.PersonIDs[index] != 0:
	case f.CharacterIDs[index] != 0:
	case f.CharacterThumbnails[index] != nil:
	default:
		return true
	}
	return false
}

func convertFormToVideo(form *addVideoForm) (models.Video, error) {
	if len(form.PersonIDs) != len(form.CharacterIDs) {
		return models.Video{}, fmt.Errorf("mismatched number of people and characters")
	}

	uploadDate, err := time.Parse(time.DateOnly, form.UploadDate)
	if err != nil {
		return models.Video{}, fmt.Errorf("unable to parse date")
	}

	creator := &models.Creator{
		ID: form.CreatorID,
	}

	slug := models.CreateSlugName(form.Title, maxFileNameLength)

	var cast []*models.CastMember
	for i := range form.PersonIDs {
		p := &models.Person{ID: &form.PersonIDs[i]}
		c := &models.Character{}
		if form.CharacterIDs[i] != 0 {
			c.ID = &form.CharacterIDs[i]
		}

		cm := &models.CastMember{
			Position:      &i,
			Actor:         p,
			Character:     c,
			ThumbnailFile: form.CharacterThumbnails[i],
		}

		cast = append(cast, cm)
	}

	return models.Video{
		Title:         form.Title,
		URL:           form.URL,
		Slug:          slug,
		ThumbnailFile: form.Thumbnail,
		Rating:        form.Rating,
		UploadDate:    &uploadDate,
		Creator:       creator,
		Cast:          cast,
	}, nil
}

