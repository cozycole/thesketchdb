package integration_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

// fakeGoogleProvider exercises the real OAuth handlers without contacting Google.
// Goth keeps a global provider registry, so these tests must NOT run in parallel.
type fakeGoogleProvider struct {
	user goth.User
}

func (p *fakeGoogleProvider) Name() string                { return "google" }
func (p *fakeGoogleProvider) SetName(string)              {}
func (p *fakeGoogleProvider) Debug(bool)                  {}
func (p *fakeGoogleProvider) RefreshTokenAvailable() bool { return false }
func (p *fakeGoogleProvider) RefreshToken(string) (*oauth2.Token, error) {
	return nil, fmt.Errorf("not implemented by test provider")
}
func (p *fakeGoogleProvider) BeginAuth(state string) (goth.Session, error) {
	return &fakeGoogleSession{state: state}, nil
}
func (p *fakeGoogleProvider) UnmarshalSession(data string) (goth.Session, error) {
	if data == "" {
		return nil, fmt.Errorf("empty OAuth session")
	}
	return &fakeGoogleSession{state: data}, nil
}
func (p *fakeGoogleProvider) FetchUser(goth.Session) (goth.User, error) {
	return p.user, nil
}

type fakeGoogleSession struct {
	state string
}

func (s *fakeGoogleSession) Marshal() string { return s.state }
func (s *fakeGoogleSession) GetAuthURL() (string, error) {
	return "https://accounts.google.test/authorize?state=" + url.QueryEscape(s.state), nil
}
func (s *fakeGoogleSession) Authorize(_ goth.Provider, params goth.Params) (string, error) {
	if params.Get("code") != "test-code" {
		return "", fmt.Errorf("invalid test authorization code")
	}
	return "test-access-token", nil
}

func installFakeGoogle(t *testing.T, user goth.User) {
	t.Helper()
	old, err := goth.GetProvider("google")
	goth.UseProviders(&fakeGoogleProvider{user: user})
	t.Cleanup(func() {
		if err == nil {
			goth.UseProviders(old)
		} else {
			goth.ClearProviders()
		}
	})
}

// oauthClient retains the SCS cookie between the login and callback requests.
// Redirects are disabled so tests can assert the intended destination.
func oauthClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func requireRedirect(
	t *testing.T,
	resp *http.Response,
	wantStatus int,
	wantLocation string,
) {
	t.Helper()
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		body, _ := io.ReadAll(resp.Body)

		t.Fatalf(
			"HTTP status: got %d, want %d\n"+
				"URL: %s\n"+
				"Body: %s",
			resp.StatusCode,
			wantStatus,
			resp.Request.URL,
			body,
		)
	}

	if location := resp.Header.Get("Location"); location != wantLocation {
		t.Fatalf(
			"redirect location: got %q, want %q",
			location,
			wantLocation,
		)
	}
}

// beginGoogleLogin performs the first leg of the OAuth flow and returns the
// state that the provider would return during its callback.
func beginGoogleLogin(t *testing.T, client *http.Client, baseURL string) string {
	t.Helper()
	resp, err := client.Get(baseURL + "/auth/google")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("Google login returned %d, want 302", resp.StatusCode)
	}
	location, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if location.Host != "accounts.google.test" {
		t.Fatalf("unexpected OAuth host: %q", location.Host)
	}
	state := location.Query().Get("state")
	if state == "" {
		t.Fatal("login omitted OAuth state")
	}
	return state
}

func googleCallbackURL(baseURL, state string) string {
	params := url.Values{"state": {state}, "code": {"test-code"}}
	return baseURL + "/auth/google/callback?" + params.Encode()
}

func uniqueSocialUser(t *testing.T) (email, username, googleID string) {
	t.Helper()

	id := strings.ReplaceAll(uuid.NewString(), "-", "")[:12]

	return "social-" + id + "@example.test",
		"social_" + id,
		"google-" + id
}

func TestGoogleAuth_ExistingLinkedIdentity(t *testing.T) {
	test := newTestApp(t, TestOptions{})
	email, username, googleID := uniqueSocialUser(t)
	ctx := context.Background()

	var userID int
	err := test.DB.QueryRow(ctx, `
		INSERT INTO users (username, email, activated)
		VALUES ($1, $2, true) RETURNING id
	`, username, email).Scan(&userID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = test.DB.Exec(context.Background(), "DELETE FROM user_identities WHERE user_id = $1", userID)
		_, _ = test.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	})

	_, err = test.DB.Exec(ctx, `
		INSERT INTO user_identities (user_id, provider, provider_user_id)
		VALUES ($1, 'google', $2)
	`, userID, googleID)
	if err != nil {
		t.Fatal(err)
	}

	installFakeGoogle(t, goth.User{UserID: googleID, Email: email, Provider: "google"})
	client := oauthClient(t)
	state := beginGoogleLogin(t, client, test.Server.URL)
	resp, err := client.Get(googleCallbackURL(test.Server.URL, state))
	if err != nil {
		t.Fatal(err)
	}
	requireRedirect(t, resp, http.StatusSeeOther, "/")

	var count int
	if err := test.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM user_identities WHERE user_id = $1 AND provider = 'google'", userID,
	).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("identity count: got %d, want 1", count)
	}

	// This verifies persistence of a session cookie, not the authenticated ID.
	// Add an authenticated endpoint assertion to verify login state end-to-end.
	serverURL, _ := url.Parse(test.Server.URL)
	if len(client.Jar.Cookies(serverURL)) == 0 {
		t.Fatal("no session cookie after Google login")
	}
}

func TestGoogleAuth_LinksExistingEmail(t *testing.T) {
	test := newTestApp(t, TestOptions{})
	email, username, googleID := uniqueSocialUser(t)
	ctx := context.Background()

	var userID int
	if err := test.DB.QueryRow(ctx, `
		INSERT INTO users (username, email, activated)
		VALUES ($1, $2, true) RETURNING id
	`, username, email).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = test.DB.Exec(context.Background(), "DELETE FROM user_identities WHERE user_id = $1", userID)
		_, _ = test.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	})

	installFakeGoogle(t, goth.User{UserID: googleID, Email: email, Provider: "google"})
	client := oauthClient(t)
	state := beginGoogleLogin(t, client, test.Server.URL)
	resp, err := client.Get(googleCallbackURL(test.Server.URL, state))
	if err != nil {
		t.Fatal(err)
	}
	requireRedirect(t, resp, http.StatusSeeOther, "/")

	var linkedUserID int
	if err := test.DB.QueryRow(ctx, `
		SELECT user_id FROM user_identities
		WHERE provider = 'google' AND provider_user_id = $1
	`, googleID).Scan(&linkedUserID); err != nil {
		t.Fatal(err)
	}
	if linkedUserID != userID {
		t.Fatalf("linked to user %d, want %d", linkedUserID, userID)
	}
}

func TestGoogleAuth_NewUserCompletesSignup(t *testing.T) {
	test := newTestApp(t, TestOptions{})
	client := oauthClient(t)

	email, username, googleID := uniqueSocialUser(t)

	// Configure the fake Google provider to return this identity.
	installFakeGoogle(t, goth.User{UserID: googleID, Email: email, Provider: "google"})

	// 1. Initiate OAuth.
	state := beginGoogleLogin(t, client, test.Server.URL)

	// 2. Simulate Google's callback.
	callbackURL := googleCallbackURL(test.Server.URL, state)

	resp, err := client.Get(callbackURL)
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()

	requireRedirect(t, resp, http.StatusSeeOther, "/complete-signup")

	// 3. submit the username using the same client/session.
	form := url.Values{
		"username": {username},
	}

	resp, err = client.PostForm(
		test.Server.URL+"/complete-signup",
		form,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	requireRedirect(t, resp, http.StatusSeeOther, "/")

	// 4. Verify the user and identity were created.
	var userID int
	var activated bool

	err = test.DB.QueryRow(
		context.Background(),
		`SELECT id, activated
         FROM users
         WHERE email = $1 AND username = $2`,
		email,
		username,
	).Scan(&userID, &activated)
	if err != nil {
		t.Fatalf("user not created: %v", err)
	}

	if !activated {
		t.Error("expected OAuth user to be activated")
	}

	var identityUserID int

	err = test.DB.QueryRow(
		context.Background(),
		`SELECT user_id
         FROM user_identities
         WHERE provider = 'google'
           AND provider_user_id = $1`,
		googleID,
	).Scan(&identityUserID)
	if err != nil {
		t.Fatalf("Google identity not created: %v", err)
	}

	if identityUserID != userID {
		t.Errorf(
			"identity linked to user %d, expected %d",
			identityUserID,
			userID,
		)
	}
}

func TestGoogleAuth_RejectsInvalidState(t *testing.T) {
	test := newTestApp(t, TestOptions{})
	email, _, googleID := uniqueSocialUser(t)
	installFakeGoogle(t, goth.User{UserID: googleID, Email: email, Provider: "google"})

	client := oauthClient(t)
	_ = beginGoogleLogin(t, client, test.Server.URL)
	resp, err := client.Get(googleCallbackURL(test.Server.URL, "tampered-state"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid state: got %d, want 400", resp.StatusCode)
	}
}

func TestGoogleAuth_RejectsCallbackWithoutSession(t *testing.T) {
	test := newTestApp(t, TestOptions{})
	installFakeGoogle(t, goth.User{UserID: "irrelevant", Email: "missing@example.test", Provider: "google"})

	client := oauthClient(t)
	resp, err := client.Get(googleCallbackURL(test.Server.URL, "missing-state"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("callback without session: got %d, want 400", resp.StatusCode)
	}
}

func TestGoogleAuth_RejectsReplayedCallback(t *testing.T) {
	test := newTestApp(t, TestOptions{})
	email, _, googleID := uniqueSocialUser(t)
	installFakeGoogle(t, goth.User{UserID: googleID, Email: email, Provider: "google"})

	client := oauthClient(t)
	state := beginGoogleLogin(t, client, test.Server.URL)
	callback := googleCallbackURL(test.Server.URL, state)
	first, err := client.Get(callback)
	if err != nil {
		t.Fatal(err)
	}
	requireRedirect(t, first, http.StatusSeeOther, "/complete-signup")

	second, err := client.Get(callback)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Body.Close()
	if second.StatusCode != http.StatusBadRequest {
		t.Fatalf("replayed callback: got %d, want 400", second.StatusCode)
	}
}

// Compile-time check against Goth's provider interface.
var _ goth.Provider = (*fakeGoogleProvider)(nil)

// Compile-time check against Goth's session interface.
var _ goth.Session = (*fakeGoogleSession)(nil)
