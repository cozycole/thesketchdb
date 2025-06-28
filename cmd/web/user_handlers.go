package main

import (
	"errors"
	"fmt"
	"net/http"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Signup = &userSignupForm{}
	app.render(r, w, http.StatusOK, "signup.gohtml", "base", data)
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
		app.render(r, w, http.StatusUnprocessableEntity, "signup.gohtml", "base", data)
		return
	}

	activated := false
	user := &models.User{
		Username:  &form.Username,
		Email:     &form.Email,
		Activated: &activated,
	}

	err = user.Password.Set(form.Password)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.users.Insert(user)
	if err != nil {
		fmt.Println(err.Error())
		data := app.newTemplateData(r)
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Validator.AddFieldError("email", "A user with this email address already exists.")
			data.Forms.Signup = &form
			app.render(r, w, http.StatusUnprocessableEntity, "signup.gohtml", "base", data)
		} else if errors.Is(err, models.ErrDuplicateUsername) {
			form.Validator.AddFieldError("username", "A user with this username already exists.")
			data.Forms.Signup = &form
			app.render(r, w, http.StatusUnprocessableEntity, "signup.gohtml", "base", data)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Successful signup! Please log in.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Login = &userLoginForm{}
	app.render(r, w, http.StatusOK, "login.gohtml", "base", data)
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
		app.render(r, w, http.StatusUnprocessableEntity, "login.gohtml", "base", data)
		return
	}

	id, err := app.users.Authenticate(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Username or password is incorrect")
			data := app.newTemplateData(r)
			data.Forms.Login = &form
			app.render(r, w, http.StatusUnprocessableEntity, "login.gohtml", "base", data)
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
	if app.sessionManager.Exists(r.Context(), "postLoginRedirectURL") {
		url := app.sessionManager.Pop(r.Context(), "postLoginRedirectURL").(string)
		http.Redirect(w, r, url, http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/browse", http.StatusSeeOther)
	}
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

	favoriteSketches, err := app.sketches.GetByUserLikes(*user.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	page, err := views.UserPageView(
		user,
		favoriteSketches,
		app.baseImgUrl,
	)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-user.gohtml", "base", data)
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
