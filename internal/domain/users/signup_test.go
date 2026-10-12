package users

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func TestUserService_Signup(t *testing.T) {
	// Sentinel errors used to simulate dependency failures.
	errDatabase := errors.New("database unavailable")
	errInsert := errors.New("insert failed")
	errDelete := errors.New("token deletion failed")
	errToken := errors.New("token creation failed")
	errMail := errors.New("email delivery failed")

	const (
		newUserID      = 100
		existingUserID = 200
		email          = "test@example.com"
		username       = "testuser"
		password       = "secure-password-123"
		origin         = "https://example.com"
	)

	tokenTTL := 24 * time.Hour

	// Test user factory.
	existingUser := func(activated bool) *models.User {
		id := existingUserID
		active := activated
		userEmail := email

		return &models.User{
			ID:        &id,
			Email:     &userEmail,
			Activated: &active,
		}
	}

	tests := []struct {
		name string

		// Configure dependency behavior.
		existingUser *models.User
		findErr      error
		insertErr    error
		deleteErr    error
		tokenErr     error
		mailErr      error

		// Expected outcome.
		wantErr    error
		wantInsert bool
		wantDelete bool
		wantToken  bool
		wantMail   bool
		wantUserID int
	}{
		{
			name:       "new user signs up",
			findErr:    models.ErrNoRecord,
			wantInsert: true,
			wantToken:  true,
			wantMail:   true,
			wantUserID: newUserID,
		},
		{
			name:         "existing activated user",
			existingUser: existingUser(true),
		},
		{
			name:         "existing unactivated user",
			existingUser: existingUser(false),
			wantDelete:   true,
			wantToken:    true,
			wantMail:     true,
			wantUserID:   existingUserID,
		},
		{
			name:    "database lookup fails",
			findErr: errDatabase,
			wantErr: errDatabase,
		},
		{
			name:       "duplicate username",
			findErr:    models.ErrNoRecord,
			insertErr:  models.ErrDuplicateUsername,
			wantErr:    ErrUsernameUnavailable,
			wantInsert: true,
		},
		{
			name:       "user insertion fails",
			findErr:    models.ErrNoRecord,
			insertErr:  errInsert,
			wantErr:    errInsert,
			wantInsert: true,
		},
		{
			name:         "deleting old tokens fails",
			existingUser: existingUser(false),
			deleteErr:    errDelete,
			wantErr:      errDelete,
			wantDelete:   true,
		},
		{
			name:       "creating activation token fails",
			findErr:    models.ErrNoRecord,
			tokenErr:   errToken,
			wantErr:    errToken,
			wantInsert: true,
			wantToken:  true,
			wantUserID: newUserID,
		},
		{
			name:         "existing user token creation fails",
			existingUser: existingUser(false),
			tokenErr:     errToken,
			wantErr:      errToken,
			wantDelete:   true,
			wantToken:    true,
			wantUserID:   existingUserID,
		},
		{
			name:       "sending activation email fails",
			findErr:    models.ErrNoRecord,
			mailErr:    errMail,
			wantErr:    errMail,
			wantInsert: true,
			wantToken:  true,
			wantMail:   true,
			wantUserID: newUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &spyUserRepository{
				foundUser: tt.existingUser,
				findErr:   tt.findErr,
				insertErr: tt.insertErr,
				nextID:    newUserID,
			}

			tokens := &spyTokenRepository{
				deleteErr: tt.deleteErr,
				createErr: tt.tokenErr,
			}

			mailer := &spyMailer{
				sendErr: tt.mailErr,
			}

			service := &UserService{
				users:  users,
				tokens: tokens,
				mailer: mailer,
				config: Config{
					origin:             origin,
					activationTokenTTL: tokenTTL,
				},
			}

			err := service.Signup(SignupInput{
				Username: username,
				Email:    email,
				Password: password,
			})

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Signup() unexpected error: %v", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Signup() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			// Assert: email lookup
			if !reflect.DeepEqual(
				users.getByEmailCalls,
				[]string{email},
			) {
				t.Errorf(
					"GetByEmail calls = %v, want [%q]",
					users.getByEmailCalls,
					email,
				)
			}

			// Assert: user insertion
			utils.AssertCallCount(
				t,
				"Insert",
				len(users.insertedUsers),
				tt.wantInsert,
			)

			if tt.wantInsert {
				inserted := users.insertedUsers[0]

				if inserted.Username == nil ||
					*inserted.Username != username {
					t.Errorf(
						"Insert username = %v, want %q",
						inserted.Username,
						username,
					)
				}

				if inserted.Email == nil ||
					*inserted.Email != email {
					t.Errorf(
						"Insert email = %v, want %q",
						inserted.Email,
						email,
					)
				}
			}

			// Assert: token deletion
			utils.AssertCallCount(
				t,
				"DeleteAllForUser",
				len(tokens.deleteCalls),
				tt.wantDelete,
			)

			if tt.wantDelete {
				got := tokens.deleteCalls[0]

				if got.scope != models.ScopeActivation {
					t.Errorf(
						"DeleteAllForUser scope = %q, want %q",
						got.scope,
						models.ScopeActivation,
					)
				}

				if got.userID != existingUserID {
					t.Errorf(
						"DeleteAllForUser userID = %d, want %d",
						got.userID,
						existingUserID,
					)
				}
			}

			// Assert: token creation
			utils.AssertCallCount(
				t,
				"Token.New",
				len(tokens.createCalls),
				tt.wantToken,
			)

			if tt.wantToken {
				got := tokens.createCalls[0]

				want := createTokenCall{
					userID: tt.wantUserID,
					ttl:    tokenTTL,
					scope:  models.ScopeActivation,
				}

				if got != want {
					t.Errorf(
						"Token.New call = %+v, want %+v",
						got,
						want,
					)
				}
			}

			// Assert: email delivery
			utils.AssertCallCount(
				t,
				"Mailer.Send",
				len(mailer.calls),
				tt.wantMail,
			)

			if tt.wantMail {
				got := mailer.calls[0]

				if got.to != email {
					t.Errorf(
						"Mailer.Send to = %q, want %q",
						got.to,
						email,
					)
				}

				if got.template != "user_welcome.gohtml" {
					t.Errorf(
						"Mailer.Send template = %q",
						got.template,
					)
				}

				wantData := map[string]string{
					"ActivationURL": origin +
						"/verify-email?token=test-token",
				}

				if !reflect.DeepEqual(got.data, wantData) {
					t.Errorf(
						"Mailer.Send data = %v, want %v",
						got.data,
						wantData,
					)
				}
			}
		})
	}
}
