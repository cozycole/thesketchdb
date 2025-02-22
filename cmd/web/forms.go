package main

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
	"sketchdb.cozycole.net/internal/validator"
)

type Forms struct {
	Creator *creatorForm
	Login   *userLoginForm
	Person  *personForm
	Signup  *userSignupForm
	Video   *videoForm
	Cast    *castForm
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
	ID                  int                   `form:"vidId"`
	Title               string                `form:"title"`
	URL                 string                `form:"url"`
	Slug                string                `form:"slug"`
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
	form.CheckField(form.CreatorID != 0, "creatorId", "This field cannot be blank")
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
	width, height, err := utils.GetImageDimensions(thumbnail)
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to determine image dimensions")
		return
	}

	form.CheckField(width >= MinThumbnailWidth && height >= MinThumbnailHeight, "thumbnail", "Thumbnail dimensions must be at least 480x360")
}

func (app *application) validateUpdateVideoForm(form *videoForm) {
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

	form.CheckField(validator.NotBlank(form.Slug), "slug", "This field cannot be blank")
	form.CheckField(
		!app.videos.IsSlugDuplicate(form.ID, form.Slug),
		"slug",
		"Slug already exists",
	)

	if form.CreatorID != 0 {
		form.CheckField(
			validator.BoolWithError(app.creators.Exists(form.CreatorID)),
			"creator",
			"Unable to find creator, please add them",
		)
	}

	if form.Thumbnail == nil {
		return
	}

	// NOTE: not DRY, see above
	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a valid jpg or png")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg")
	width, height, err := utils.GetImageDimensions(thumbnail)
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to determine image dimensions")
		return
	}

	form.CheckField(width >= MinThumbnailWidth && height >= MinThumbnailHeight,
		"thumbnail", "Thumbnail dimensions must be at least 480x360")
}

type castForm struct {
	PersonID            int                   `form:"personId"`
	PersonInput         string                `form:"personInput"`
	CharacterName       string                `form:"characterName"`
	CharacterID         int                   `form:"characterId"`
	CharacterInput      string                `form:"characterInput"`
	CharacterThumbnail  *multipart.FileHeader `img:"characterThumbnail"`
	CharacterProfile    *multipart.FileHeader `img:"characterProfile"`
	validator.Validator `form:"-"`
}

func (app *application) validateAddCast(form *castForm) {
	pid := form.PersonID
	form.CheckField(pid != 0, "personId", "This field cannot be blank. Select a person from dropdown.")
	if pid != 0 {
		form.CheckField(
			validator.BoolWithError(app.people.Exists(pid)),
			"personId",
			"Person does not exist. Please add it, then resubmit video!",
		)
	}

	form.CheckField(validator.NotBlank(form.CharacterName), "characterName", "This field cannot be blank.")

	cid := form.CharacterID
	if cid != 0 {
		form.CheckField(
			validator.BoolWithError(app.characters.Exists(cid)),
			"characterId",
			"Character does not exist. Please add it, then resubmit video!",
		)
	}

	thumb := form.CharacterThumbnail
	form.CheckField(thumb != nil, "characterThumbnail", "Please upload an image")
	if thumb != nil {
		thumbnail, err := thumb.Open()
		if err != nil {
			form.AddFieldError("characterThumbnail", "File error, ensure it is a jpg")
			return
		}

		form.CheckField(validator.IsMime(thumbnail, "image/jpeg"),
			"characterThumbnail", "Uploaded file must be jpg or png")

		width, height, err := utils.GetImageDimensions(thumbnail)
		if err != nil {
			app.errorLog.Print(err)
			form.AddFieldError("characterThumbnail", "Unable to determine image dimensions")
			return
		}

		form.CheckField(width >= MinThumbnailWidth && height >= MinThumbnailHeight,
			"characterThumbnail", "Thumbnail dimensions must be at least 480x360")
	}

	profile := form.CharacterProfile
	form.CheckField(thumb != nil, "characterProfile", "Please upload an image")
	if profile == nil {
		return
	}

	thumbnail, err := profile.Open()
	if err != nil {
		form.AddFieldError("characterProfile", "File error, ensure it is a jpg")
		return
	}

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg"),
		"characterProfile", "Uploaded file must be jpg")

	width, _, err := utils.GetImageDimensions(thumbnail)
	form.CheckField(width >= MinProfileWidth, "characterProfile", "Ensure photo width is larger than 256")

	thumbnail.Close()
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

func convertFormToVideo(form *videoForm) (models.Video, error) {
	uploadDate, err := time.Parse(time.DateOnly, form.UploadDate)
	if err != nil {
		return models.Video{}, fmt.Errorf("unable to parse date")
	}

	creator := &models.Creator{
		ID: &form.CreatorID,
	}

	return models.Video{
		Title:         form.Title,
		URL:           &form.URL,
		Slug:          form.Slug,
		ThumbnailFile: form.Thumbnail,
		Rating:        form.Rating,
		UploadDate:    &uploadDate,
		Creator:       creator,
	}, nil
}

func convertFormtoCastMember(form *castForm) models.CastMember {
	actor := models.Person{ID: &form.PersonID}
	character := models.Character{}
	if form.CharacterID != 0 {
		character.ID = &form.CharacterID
	}
	return models.CastMember{
		Actor:         &actor,
		Character:     &character,
		CharacterName: &form.CharacterName,
		ThumbnailFile: form.CharacterThumbnail,
		ProfileFile:   form.CharacterProfile,
	}
}
