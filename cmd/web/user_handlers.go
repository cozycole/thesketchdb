package main

import (
	"errors"
	"net/http"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Signup = &userSignupForm{}
	app.render(r, w, http.StatusOK, "signup.tmpl.html", "base", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	app.infoLog.Println(err)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if app.validateUserSignupForm(&form); !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Signup = &form
		app.render(r, w, http.StatusUnprocessableEntity, "signup.tmpl.html", "base", data)
		return
	}

	user := &models.User{
		Username:  form.Username,
		Email:     form.Email,
		Activated: true,
	}

	err = user.Password.Set(form.Password)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.users.Insert(user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Validator.AddFieldError("email", "a user with this email address already exists")
			data := app.newTemplateData(r)
			data.Forms.Signup = &form
			app.render(r, w, http.StatusUnprocessableEntity, "signup.tmpl.html", "base", data)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Successful signup! Please log in.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Login = &userLoginForm{}
	app.render(r, w, http.StatusOK, "login.tmpl.html", "base", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateUserLoginForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Login = &form
		app.render(r, w, http.StatusUnprocessableEntity, "login.tmpl.html", "base", data)
		return
	}

	id, err := app.users.Authenticate(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Forms.Login = &form
			app.render(r, w, http.StatusUnprocessableEntity, "login.tmpl.html", "base", data)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/search", http.StatusSeeOther)
}

type UserPage struct {
	User      *models.User
	Favorited []*models.Video
}

func (app *application) userView(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	user, err := app.users.GetByUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	videos, err := app.videos.GetByUserLikes(user.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.UserPage.User = user
	data.UserPage.Favorited = videos
	app.render(r, w, http.StatusOK, "view-user.tmpl.html", "base", data)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
