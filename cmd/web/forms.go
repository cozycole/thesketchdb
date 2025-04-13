package main

import (
	"fmt"
	"mime/multipart"
	"strings"

	"sketchdb.cozycole.net/internal/utils"
	"sketchdb.cozycole.net/internal/validator"
)

type Forms struct {
	Cast      *castForm
	Category  *categoryForm
	Creator   *creatorForm
	Episode   *episodeForm
	Login     *userLoginForm
	Person    *personForm
	Show      *showForm
	Signup    *userSignupForm
	Tag       *tagForm
	Video     *videoForm
	VideoTags *videoTagsForm
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

type categoryForm struct {
	Name                string `form:"categoryName"`
	ParentId            int    `form:"parentId"`
	ParentInput         string `form:"parentInput"`
	validator.Validator `form:"-"`
}

func (app *application) validateCategoryForm(form *categoryForm) {
	form.CheckField(validator.NotBlank(form.Name), "categoryName", "Field cannot be blank")
	if form.ParentId != 0 {
		form.CheckField(
			validator.BoolWithError(app.categories.Exists(form.ParentId)),
			"parentId",
			"Category does not exist. Please add it, then resubmit category!",
		)
	}
}

type tagForm struct {
	Name                string `form:"tagName"`
	CategoryId          int    `form:"categoryId"`
	CategoryInput       string `form:"categoryInput"`
	validator.Validator `form:"-"`
}

func (app *application) validateTagForm(form *tagForm) {
	form.CheckField(validator.NotBlank(form.Name), "tagName", "Field cannot be blank")
	if form.CategoryId != 0 {
		form.CheckField(
			validator.BoolWithError(app.categories.Exists(form.CategoryId)),
			"categoryId",
			"Category does not exist. Please add it, then resubmit tag!",
		)
	}
}

type videoTagsForm struct {
	TagIds              []int    `form:"tagId"`
	TagInputs           []string `form:"tagName"`
	validator.Validator `form:"-"`
}

func (app *application) validateVideoTagsForm(form *videoTagsForm) {
	for i, tagId := range form.TagIds {
		if tagId == 0 {
			form.AddMultiFieldError("tagId", i, "Please input tag")
		}

		if exists, _ := app.tags.Exists(tagId); !exists {
			form.AddMultiFieldError("tagId", i, "Tag does not exist")
		}
	}
}

type showForm struct {
	Name                string                `form:"name"`
	Slug                string                `form:"slug"`
	ProfileImg          *multipart.FileHeader `img:"profileImg"`
	validator.Validator `form:"-"`
}

func (app *application) validateShowForm(form *showForm) {
	form.CheckField(validator.NotBlank(form.Name), "name", "Field cannot be blank")

	thumbnail, err := form.ProfileImg.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a valid jpg or png")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg")
}

func (app *application) validateUpdateShowForm(form *showForm) {
	app.infoLog.Printf("%+v\n", form)
	form.CheckField(validator.NotBlank(form.Name), "name", "Field cannot be blank")
	form.CheckField(validator.NotBlank(form.Slug), "slug", "Field cannot be blank")

	if form.ProfileImg == nil {
		return
	}

	thumbnail, err := form.ProfileImg.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a valid jpg or png")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg")
}

type episodeForm struct {
	ID                  int                   `form:"id"`
	Number              int                   `form:"number"`
	Title               string                `form:"title"`
	AirDate             string                `form:"airDate"`
	Thumbnail           *multipart.FileHeader `img:"thumbnail"`
	SeasonId            int                   `form:"seasonId"`
	validator.Validator `form:"-"`
}

func (app *application) validateEpisodeForm(form *episodeForm) {
	form.CheckField(form.Number != 0, "number", "Please enter a valid number")

	// validate episode number
	if form.SeasonId == 0 {
		form.AddNonFieldError("Season ID not defined")
	}
	season, err := app.shows.GetSeason(form.SeasonId)
	if err != nil {
		app.errorLog.Printf("Error getting season for episode form validation: $s", err)
		form.AddNonFieldError("Error getting season")
	}

	for _, ep := range season.Episodes {
		if ep.Number != nil && *ep.Number == form.Number {
			form.AddFieldError("number", fmt.Sprintf("Episode number %d already exists for season %d", *ep.Number, *season.Number))
		}
	}

	form.CheckField(validator.NotBlank(form.AirDate), "airDate", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.AirDate), "airDate", "Date must be of the format YYYY-MM-DD")

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

func (app *application) validateUpdateEpisodeForm(form *episodeForm) {
	if form.ID == 0 {
		app.errorLog.Println("Episode ID not defined in form")
		form.AddNonFieldError("Episode ID not defined in form")
		return
	}

	form.CheckField(form.Number != 0, "number", "Please enter a valid number")

	episode, err := app.shows.GetEpisode(form.ID)
	if err != nil {
		app.errorLog.Printf("Error getting episode, bad episode ID: %s", err)
		form.AddNonFieldError("Error getting episode, bad episode ID")
		return
	}

	if episode.Number != nil && *episode.Number != form.Number {
		season, err := app.shows.GetSeason(form.SeasonId)
		if err != nil {
			app.errorLog.Printf("Error getting season for episode form validation: $s", err)
			form.AddNonFieldError("Error getting season")
		}

		for _, ep := range season.Episodes {
			if ep.Number != nil && *ep.Number == form.Number {
				form.AddFieldError("number", fmt.Sprintf("Episode number %d already exists for season %d", *ep.Number, *season.Number))
			}
		}
	}

	form.CheckField(validator.NotBlank(form.AirDate), "airDate", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.AirDate), "airDate", "Date must be of the format YYYY-MM-DD")

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
