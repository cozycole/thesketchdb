package main

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/validator"
)

type Forms struct {
	Creator    creatorForm
	Login      userLoginForm
	Person     personForm
	Signup     userSignupForm
	Video      videoForm
	VideoActor videoActorForm
}

// Changes to the form fields must be updated in their respective
// validate functions
type creatorForm struct {
	Name                string                `form:"name"`
	URL                 string                `form:"url"`
	EstablishedDate     string                `form:"establishedDate"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	validator.Validator `form:"-"`
}

func (app *application) validateAddCreatorForm(form *creatorForm) {
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

type personForm struct {
	First               string                `form:"first"`
	Last                string                `form:"last"`
	BirthDate           string                `form:"birthDate"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	validator.Validator `form:"-"`
}

func (app *application) validateAddPersonForm(form *personForm) {
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
type videoForm struct {
	Title               string                `form:"title"`
	URL                 string                `form:"url"`
	Rating              string                `form:"rating"`
	UploadDate          string                `form:"uploadDate"`
	Thumbnail           *multipart.FileHeader `img:"thumbnail"`
	CreatorID           int                   `form:"creatorId"`
	CreatorInput        string                `form:"creatorInput"`
	validator.Validator `form:"-"`
}

// We need this function to have access to the apps state
// to validate based on database queries
func (app *application) validateAddVideoForm(form *videoForm) {
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

	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a valid jpg or png")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg or png")
	// Might want to check dimension ratios to make sure they work?
}

type videoActorForm struct {
	PersonIDs           []int                   `form:"peopleId"`
	PersonInputs        []string                `form:"peopleText"`
	CharacterIDs        []int                   `form:"characterId"`
	CharacterInputs     []string                `form:"characterText"`
	CharacterThumbnails []*multipart.FileHeader `img:"characterThumbnail"`
	validator.Validator `form:"-"`
}

func (app *application) validateAddActor(form *videoActorForm) {
	// The following fields are slices and are related records
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
		}

		cid := form.CharacterIDs[i]
		if !validator.IsZero(cid) {
			form.CheckField(
				validator.BoolWithError(app.characters.Exists(cid)),
				htmlCharIdField,
				"Character does not exist. Please add it, then resubmit video!",
			)
		}

		thumb := form.CharacterThumbnails[i]
		form.CheckField(thumb != nil, htmlCharThumbField, "Please upload an image")
		if thumb != nil {
			thumbnail, err := thumb.Open()
			if err != nil {
				form.AddFieldError(htmlCharThumbField, "Unable to open file, ensure it is a jpg or png")
				return
			}

			form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"),
				htmlCharThumbField, "Uploaded file must be jpg or png")

			thumbnail.Close()
		}
	}
}

type userSignupForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) validateUserSignupForm(form *userSignupForm) {
	form.CheckField(validator.NotBlank(form.Username), "username", "Field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegEx), "email", "Please enter a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "Field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Must be 8-20 chararacters")
	form.CheckField(validator.MaxChars(form.Password, 20), "password", "Must be 8-20 chararacters")
}

type userLoginForm struct {
	Username            string `form:"username"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) validateUserLoginForm(form *userLoginForm) {
	form.CheckField(validator.NotBlank(form.Username), "username", "Field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Field cannot be blank")
}

func (f *videoActorForm) IsEmptyActorInput(index int) bool {
	switch {
	case f.PersonIDs[index] != 0:
	case f.CharacterIDs[index] != 0:
	case f.CharacterThumbnails[index] != nil:
	default:
		return true
	}
	return false
}

func convertFormToVideo(form *videoForm) (models.Video, error) {
	uploadDate, err := time.Parse(time.DateOnly, form.UploadDate)
	if err != nil {
		return models.Video{}, fmt.Errorf("unable to parse date")
	}

	creator := &models.Creator{
		ID: &form.CreatorID,
	}

	slug := models.CreateSlugName(form.Title, maxFileNameLength)
	slug = slug + "-" + models.GetTimeStampHash()

	return models.Video{
		Title:         form.Title,
		URL:           &form.URL,
		Slug:          slug,
		ThumbnailFile: form.Thumbnail,
		Rating:        form.Rating,
		UploadDate:    &uploadDate,
		Creator:       creator,
	}, nil
}
