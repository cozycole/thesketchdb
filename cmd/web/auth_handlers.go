package main

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"time"

	"github.com/markbates/goth"

	"github.com/cozycole/thesketchdb/internal/domain/users"
	"github.com/cozycole/thesketchdb/internal/models"
	"github.com/cozycole/thesketchdb/internal/utils"
)

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Signup = &userSignupForm{}
	app.render(r, w, http.StatusOK, "signup.gohtml", "base", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	app.InfoLog.Println(err)
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

	input := users.SignupInput{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}

	err = app.Services.Users.Signup(input)
	if err != nil {
		data := app.newTemplateData(r)
		if errors.Is(err, users.ErrUsernameUnavailable) {
			form.Validator.AddFieldError("username", "A user with this username already exists.")
			data.Forms.Signup = &form
			app.render(r, w, http.StatusUnprocessableEntity, "signup.gohtml", "base", data)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	http.Redirect(w, r, "/verify-email", http.StatusSeeOther)
}

func (app *Application) verifyEmail(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Token = r.URL.Query().Get("token")

	// if there's no token, we are just rendering the "check your email" blurb
	if data.Token == "" {
		app.render(r, w, http.StatusOK, "verify-email.gohtml", "base", data)
		return
	}

	_, err := app.Users.GetByToken(data.Token, models.ScopeActivation)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			data := app.newTemplateData(r)
			data.TokenScope = "activation"
			app.render(r, w, http.StatusOK, "invalid-token.gohtml", "base", data)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.render(r, w, http.StatusOK, "verify-email.gohtml", "base", data)
}

func (app *Application) activateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token := r.PostForm.Get("token")
	user, err := app.Services.Users.Activate(token)
	if err != nil {
		if errors.Is(err, users.ErrInvalidToken) {
			data := app.newTemplateData(r)
			data.TokenScope = "activation"
			app.render(r, w, http.StatusOK, "invalid-token.gohtml", "base", data)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	// login user
	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "authenticatedUserID", *user.ID)
	flash := flashMessage{
		Message: "Account activated!",
		Level:   FlashSuccess,
	}
	app.SessionManager.Put(r.Context(), "flash", flash)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) resendActivationLink(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.ResendActivation = &userResendActivationForm{}
	app.render(r, w, http.StatusOK, "resend-activation.gohtml", "base", data)
}

func (app *Application) resendActivationLinkPost(w http.ResponseWriter, r *http.Request) {
	var form userResendActivationForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateUserResendActivation(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.ResendActivation = &form
		app.render(r, w, http.StatusUnprocessableEntity, "resend-activation.gohtml", "base", data)
		return
	}

	err = app.Services.Users.ResendActivationLink(form.Email)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	http.Redirect(w, r, "/verify-email", http.StatusSeeOther)
}

func (app *Application) forgotPassword(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.ForgotPassword = &userForgotPasswordForm{}
	app.render(r, w, http.StatusOK, "forgot-password.gohtml", "base", data)
}

func (app *Application) forgotPasswordPost(w http.ResponseWriter, r *http.Request) {
	var form userForgotPasswordForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateUserForgotPassword(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.ForgotPassword = &form
		app.render(r, w, http.StatusUnprocessableEntity, "forgot-password.gohtml", "base", data)
		return
	}

	err = app.Services.Users.ForgotPassword(form.Email)

	http.Redirect(w, r, "/password-reset-sent", http.StatusSeeOther)
}

func (app *Application) passwordResetSent(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(r, w, http.StatusOK, "password-reset-sent.gohtml", "base", data)
}

func (app *Application) resetPassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := app.newTemplateData(r)
	data.Token = r.URL.Query().Get("token")
	data.Forms.ResetPassword = &userResetPasswordForm{}

	_, err := app.Users.GetByToken(data.Token, models.ScopePasswordReset)
	if errors.Is(err, models.ErrNoRecord) {
		data.TokenScope = "password-reset"
		app.render(r, w, http.StatusOK, "invalid-token.gohtml", "base", data)
		return
	}

	app.render(r, w, http.StatusOK, "reset-password.gohtml", "base", data)
}

func (app *Application) resetPasswordPost(w http.ResponseWriter, r *http.Request) {
	var form userResetPasswordForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	data := app.newTemplateData(r)
	data.Token = form.Token
	app.validateUserPasswordReset(&form)
	if !form.Valid() {
		data.Forms.ResetPassword = &form
		app.render(r, w, http.StatusUnprocessableEntity, "reset-password.gohtml", "base", data)
		return
	}

	err = app.Services.Users.ResetPassword(form.Token, form.Password)
	if err != nil {
		if errors.Is(err, users.ErrInvalidToken) {
			data.TokenScope = "password-reset"
			app.render(r, w, http.StatusOK, "invalid-token.gohtml", "base", data)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.SessionManager.Put(r.Context(), "flash", flashMessage{
		Message: "Password reset! Please log in.",
		Level:   "success",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Login = &userLoginForm{}
	app.render(r, w, http.StatusOK, "login.gohtml", "base", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
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

	input := users.LoginInput{
		Email:    form.Email,
		Password: form.Password,
	}
	user, err := app.Services.Users.Login(input)
	if err != nil {
		if errors.Is(err, users.ErrInvalidLoginCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Forms.Login = &form
			app.render(r, w, http.StatusUnprocessableEntity, "login.gohtml", "base", data)
		} else if errors.Is(err, users.ErrUserNotActivated) {
			http.Redirect(w, r, "/verify-email", http.StatusSeeOther)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "authenticatedUserID", utils.SafeDeref(user.ID))
	if app.SessionManager.Exists(r.Context(), "postLoginRedirectURL") {
		url := app.SessionManager.Pop(r.Context(), "postLoginRedirectURL").(string)
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.SessionManager.Remove(r.Context(), "authenticatedUserID")

	flash := flashMessage{
		Message: "You've been logged out successfully!",
		Level:   FlashSuccess,
	}
	app.SessionManager.Put(r.Context(), "flash", flash)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SOCIAL LOGIN

// Used for storing info before redirecting user
// to username form if they need to create an account
type PendingRegistration struct {
	GoogleID  string
	Email     string
	ExpiresAt time.Time
}

func (app *Application) googleLogin(w http.ResponseWriter, r *http.Request) {
	provider, err := goth.GetProvider("google")
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	state, err := generateSecureState()
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	session, err := provider.BeginAuth(state)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	// Store temporary OAuth information using SCS.
	app.SessionManager.Put(r.Context(), "oauth_state", state)
	app.SessionManager.Put(r.Context(), "oauth_session", session.Marshal())

	authURL, err := session.GetAuthURL()
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (app *Application) googleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	provider, err := goth.GetProvider("google")
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	// Validate OAuth state.
	expectedState := app.SessionManager.PopString(ctx, "oauth_state")
	receivedState := r.URL.Query().Get("state")

	if expectedState == "" ||
		receivedState == "" ||
		subtle.ConstantTimeCompare(
			[]byte(expectedState),
			[]byte(receivedState),
		) != 1 {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	// Retrieve Goth's temporary session.
	sessionData := app.SessionManager.PopString(ctx, "oauth_session")
	if sessionData == "" {
		http.Error(w, "OAuth session expired", http.StatusBadRequest)
		return
	}

	session, err := provider.UnmarshalSession(sessionData)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	// Exchange authorization code for tokens.
	_, err = session.Authorize(provider, r.URL.Query())
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	// Retrieve the Google user's identity.
	googleUser, err := provider.FetchUser(session)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	result, err := app.Services.Users.AuthenticateSocial(
		"google",
		googleUser.UserID,
		googleUser.Email,
	)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	switch result.Status {

	case users.SocialAuthAuthenticated:
		if err := app.SessionManager.RenewToken(ctx); err != nil {
			app.serverError(r, w, err)
			return
		}

		app.SessionManager.Put(
			ctx,
			"authenticatedUserID",
			result.User.ID,
		)

		flash := flashMessage{
			Message: "Account linked!",
			Level:   FlashSuccess,
		}

		app.SessionManager.Put(r.Context(), "flash", flash)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case users.SocialAuthRegistrationRequired:
		pending := PendingRegistration{
			GoogleID:  googleUser.UserID,
			Email:     googleUser.Email,
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}

		app.SessionManager.Put(ctx, "pendingRegistration", pending)

		http.Redirect(w, r, "/complete-signup", http.StatusSeeOther)
	}
}
