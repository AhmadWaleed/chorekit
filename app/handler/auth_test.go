package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/session"
	"github.com/ahmadwaleed/choreui/app/database/model"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/gorilla/sessions"
)

var testuser = User{
	Name:  "john",
	Email: "j.doe@example.com",
}

type FakeHasher struct{}

func (c *FakeHasher) Generate(password string) (string, error) {
	return string("hash-password"), nil
}

func (c *FakeHasher) Verify(hash, passowrd string) bool {
	return true
}

type FakeSession struct{}

func (fs *FakeSession) Get(r *http.Request, name string) (*sessions.Session, error) {
	return &sessions.Session{}, nil
}

func (fs *FakeSession) New(r *http.Request, name string) (*sessions.Session, error) {
	return &sessions.Session{}, nil
}

func (fs *FakeSession) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

func TestSignupGet(t *testing.T) {
	a := srv.app.Echo.Group("/auth")
	a.GET("/signup", SignupGet)

	req := httptest.NewRequest("GET", "/auth/signup", nil)
	rec := httptest.NewRecorder()

	srv.app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("could not visit '/auth/signup' page, status want %d got %d", http.StatusOK, rec.Code)
	}
}

func TestSignupPost(t *testing.T) {
	a := srv.app.Echo.Group("/auth")
	a.POST("/signup", SignupPost)

	cc := core.AppContext{
		App:          srv.app,
		Loc:          i18n.New(),
		SessionStore: session.NewSessionStore,
	}

	srv.app.Echo.Use(core.AppCtxMiddleware(&cc))

	body := fmt.Sprintf("name=%s&email=%s&password=secret", testuser.Name, testuser.Email)
	req := httptest.NewRequest("POST", "/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("could not signup new user, status want %d got %d", http.StatusOK, rec.Code)
	}
}

func TestSignInGet(t *testing.T) {
	a := srv.app.Echo.Group("/auth")
	a.GET("/signin", SignupGet)

	req := httptest.NewRequest("GET", "/auth/signin", nil)
	rec := httptest.NewRecorder()

	srv.app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("could not visit '/auth/signin' page, status want %d got %d", http.StatusOK, rec.Code)
	}
}

func TestSignInPost(t *testing.T) {
	a := srv.app.Echo.Group("/auth")
	a.POST("/signin", SignInPost)

	cc := core.AppContext{
		App:          srv.app,
		Loc:          i18n.New(),
		SessionStore: session.NewSessionStore,
	}

	srv.app.Echo.Use(core.AppCtxMiddleware(&cc))

	body := fmt.Sprintf("email=%s&password=secret", testuser.Email)
	req := httptest.NewRequest("POST", "/auth/signin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.app.Echo.ServeHTTP(rec, req)

	sess := cc.SessionStore(cc.Context)

	if !sess.GetBool("Auth") {
		t.Error("could not get login value from session store")
	}

	login := sess.GetBool("Auth")
	if !sess.GetBool("Auth") {
		t.Errorf("unexpected session auth value, want: %t, got: %t", true, login)
	}

	user := (sess.Values["User"]).(model.User)
	if user.Email != testuser.Email {
		t.Errorf("unexpected session user email, want: %s, got: %s", testuser.Email, user.Email)
	}

	if rec.Code != http.StatusSeeOther {
		t.Errorf("could not login user, status want %d got %d", http.StatusOK, rec.Code)
	}
}
